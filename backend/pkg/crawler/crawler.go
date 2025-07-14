package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/seo-crawler-app/internal/models"
)

// CrawlerService defines the interface for crawling operations
type CrawlerService interface {
	CrawlURL(userID int, url string, repo Repository) error
}

// Repository defines the interface for database operations needed by crawler
type Repository interface {
	GetCrawlResultIDByURL(userID int, url string) (int, error)
	UpdateCrawlResultStatus(userID int, url, status string) error
	UpdateCrawlData(userID int, url string, data *models.CrawlData) error
	CreateLink(crawlID int, link *models.LinkData) (int64, error)
	UpdateLinkStatus(linkID int64, statusCode int, isAccessible bool) error
	CreateHeading(crawlID int, heading *models.HeadingData) error
}

// Crawler implements the CrawlerService interface
type Crawler struct{}

// NewCrawler creates a new crawler instance
func NewCrawler() CrawlerService {
	return &Crawler{}
}

// CrawlURL crawls a given URL and stores the results
func (c *Crawler) CrawlURL(userID int, baseURL string, repo Repository) error {
	collector := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	// Limit the number of threads
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	data := &models.CrawlData{
		Headings:       make(map[string]int),
		Status:         "running",
		Links:          []models.LinkData{},
		HeadingDetails: []models.HeadingData{},
	}

	// Get the crawl result ID for storing detailed data
	var crawlResultID int
	maxTries := 5
	for i := 0; i < maxTries; i++ {
		var err error
		crawlResultID, err = repo.GetCrawlResultIDByURL(userID, baseURL)
		if err == nil {
			break
		}
		if i == maxTries-1 {
			log.Printf("Failed to get crawl result ID for %s after %d tries: %v", baseURL, maxTries, err)
			return err
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Update status to running
	if err := repo.UpdateCrawlResultStatus(userID, baseURL, "running"); err != nil {
		log.Printf("Failed to update status to running for %s: %v", baseURL, err)
		return err
	}

	// Set up response handler
	collector.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		if strings.Contains(body, "<!DOCTYPE html>") {
			data.HTMLVersion.String = "HTML5"
			data.HTMLVersion.Valid = true
		} else if strings.Contains(body, "<!DOCTYPE") {
			data.HTMLVersion.String = "HTML4"
			data.HTMLVersion.Valid = true
		} else {
			data.HTMLVersion.String = "Unknown"
			data.HTMLVersion.Valid = true
		}
	})

	// Set up title handler
	collector.OnHTML("title", func(e *colly.HTMLElement) {
		data.Title.String = strings.TrimSpace(e.Text)
		data.Title.Valid = true
	})

	// Set up heading handler
	headingOrder := 0
	collector.OnHTML("h1, h2, h3, h4, h5, h6", func(e *colly.HTMLElement) {
		data.Headings[e.Name]++
		headingOrder++
		
		headingText := strings.TrimSpace(e.Text)
		if headingText != "" {
			heading := &models.HeadingData{
				Level: e.Name,
				Text:  headingText,
				Order: headingOrder,
			}
			
			// Store heading details in database
			if err := repo.CreateHeading(crawlResultID, heading); err != nil {
				log.Printf("Failed to store heading for %s: %v", baseURL, err)
			}
			
			// Add to local data
			data.HeadingDetails = append(data.HeadingDetails, *heading)
		}
	})

	// Check for login forms
	collector.OnHTML("form", func(e *colly.HTMLElement) {
		formHTML, err := e.DOM.Html()
		if err != nil {
			return
		}
		formHTML = strings.ToLower(formHTML)
		if strings.Contains(formHTML, "login") || 
		   strings.Contains(formHTML, "password") || 
		   strings.Contains(formHTML, "username") ||
		   strings.Contains(formHTML, "email") {
			data.HasLoginForm = true
		}
	})

	// Set up link handler
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		u, err := url.Parse(link)
		if err != nil {
			return
		}
		
		baseURLParsed, err := url.Parse(baseURL)
		if err != nil {
			return
		}

		linkText := strings.TrimSpace(e.Text)
		linkType := "external"
		if u.Hostname() == baseURLParsed.Hostname() {
			data.InternalLinks++
			linkType = "internal"
		} else {
			data.ExternalLinks++
		}

		linkData := &models.LinkData{
			URL:          link,
			Text:         linkText,
			Type:         linkType,
			StatusCode:   0,
			IsAccessible: true,
		}

		// Store link in database
		linkID, err := repo.CreateLink(crawlResultID, linkData)
		if err != nil {
			log.Printf("Failed to store link for %s (link: %s): %v", baseURL, link, err)
			return
		}
		
		// Add to local data
		linkData.ID = int(linkID)
		data.Links = append(data.Links, *linkData)

		// Check link status asynchronously and update database
		go func(link string, linkID int64) {
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Head(link)
			statusCode := 0
			isAccessible := true
			
			if err != nil {
				statusCode = 0
				isAccessible = false
			} else {
				defer resp.Body.Close()
				statusCode = resp.StatusCode
				if resp.StatusCode >= 400 {
					isAccessible = false
					data.InaccessibleLinks++
				}
			}
			
			// Update link status in database
			if err := repo.UpdateLinkStatus(linkID, statusCode, isAccessible); err != nil {
				log.Printf("Failed to update link status for %s: %v", link, err)
			}
			
			// Update local data
			for i := range data.Links {
				if data.Links[i].ID == int(linkID) {
					data.Links[i].StatusCode = statusCode
					data.Links[i].IsAccessible = isAccessible
					break
				}
			}
		}(link, linkID)
	})

	// Set up completion handler
	collector.OnScraped(func(r *colly.Response) {
		if err := repo.UpdateCrawlData(userID, baseURL, data); err != nil {
			log.Printf("Failed to update crawl data for %s: %v", baseURL, err)
		} else {
			log.Printf("Finished crawling: %s", baseURL)
		}
	})

	// Set up error handler
	collector.OnError(func(r *colly.Response, err error) {
		if dbErr := repo.UpdateCrawlResultStatus(userID, baseURL, "error"); dbErr != nil {
			log.Printf("Failed to update error status for %s: %v", baseURL, dbErr)
		}
		log.Printf("Error crawling %s: %v", baseURL, err)
	})

	// Start crawling
	if err := collector.Visit(baseURL); err != nil {
		log.Printf("Failed to start crawling %s: %v", baseURL, err)
		if dbErr := repo.UpdateCrawlResultStatus(userID, baseURL, "error"); dbErr != nil {
			log.Printf("Failed to update error status for %s: %v", baseURL, dbErr)
		}
		return err
	}

	collector.Wait()
	return nil
} 