package services

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/seo-crawler-app/internal/database"
	"github.com/seo-crawler-app/internal/models"
	"github.com/seo-crawler-app/pkg/crawler"
)

// CrawlService defines the interface for crawl business logic
type CrawlService interface {
	SubmitCrawl(userID int, crawlReq *models.CrawlRequest) error
	GetCrawlResults(userID int, page, pageSize int, status, search, sortBy, sortOrder string) (*models.PaginationResponse, error)
	GetCrawlResultByID(userID int, id string) (*models.URLData, error)
	GetLinksByCrawlID(id string) ([]models.LinkData, error)
	GetHeadingsByCrawlID(id string) ([]models.HeadingData, error)
	BulkRerun(userID int, urls []string) error
	BulkDelete(userID int, urls []string) error
	StopCrawl(userID int, id string) error
}

// crawlService implements the CrawlService interface
type crawlService struct {
	repo    database.Repository
	crawler crawler.CrawlerService
}

// NewCrawlService creates a new crawl service
func NewCrawlService(repo database.Repository, crawler crawler.CrawlerService) CrawlService {
	return &crawlService{
		repo:    repo,
		crawler: crawler,
	}
}

// SubmitCrawl submits a URL for crawling
func (s *crawlService) SubmitCrawl(userID int, crawlReq *models.CrawlRequest) error {
	log.Printf("SubmitCrawl called with userID: %d, URL: %s", userID, crawlReq.URL)
	
	// Validate URL
	if _, err := url.ParseRequestURI(crawlReq.URL); err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if URL already exists and update status to pending for re-crawl
	// For now, we'll create a new entry each time
	_, err := s.repo.CreateCrawlResult(userID, crawlReq.URL)
	if err != nil {
		return fmt.Errorf("failed to create crawl result: %w", err)
	}

	// Start crawling in background
	go func() {
		if err := s.crawler.CrawlURL(userID, crawlReq.URL, s.repo); err != nil {
			log.Printf("Failed to crawl URL %s: %v", crawlReq.URL, err)
		}
	}()

	return nil
}

// GetCrawlResults retrieves paginated crawl results
func (s *crawlService) GetCrawlResults(userID int, page, pageSize int, status, search, sortBy, sortOrder string) (*models.PaginationResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	results, total, err := s.repo.GetCrawlResults(userID, page, pageSize, status, search, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to get crawl results: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &models.PaginationResponse{
		Data:       results,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetCrawlResultByID retrieves a specific crawl result by ID
func (s *crawlService) GetCrawlResultByID(userID int, id string) (*models.URLData, error) {
	crawlID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	urlData, err := s.repo.GetCrawlResultByID(userID, crawlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get crawl result: %w", err)
	}

	// Get detailed links
	links, err := s.repo.GetLinksByCrawlID(crawlID)
	if err != nil {
		log.Printf("Failed to get links for crawl ID %d: %v", crawlID, err)
	} else {
		urlData.CrawlData.Links = links
	}

	// Get detailed headings
	headings, err := s.repo.GetHeadingsByCrawlID(crawlID)
	if err != nil {
		log.Printf("Failed to get headings for crawl ID %d: %v", crawlID, err)
	} else {
		urlData.CrawlData.HeadingDetails = headings
	}

	return urlData, nil
}

// GetLinksByCrawlID retrieves links for a specific crawl result
func (s *crawlService) GetLinksByCrawlID(id string) ([]models.LinkData, error) {
	crawlID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	links, err := s.repo.GetLinksByCrawlID(crawlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get links: %w", err)
	}

	return links, nil
}

// GetHeadingsByCrawlID retrieves headings for a specific crawl result
func (s *crawlService) GetHeadingsByCrawlID(id string) ([]models.HeadingData, error) {
	crawlID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	headings, err := s.repo.GetHeadingsByCrawlID(crawlID)
	if err != nil {
		return nil, fmt.Errorf("failed to get headings: %w", err)
	}

	return headings, nil
}

// BulkRerun re-runs crawling for multiple URLs
func (s *crawlService) BulkRerun(userID int, urls []string) error {
	if err := s.repo.BulkUpdateStatus(userID, urls, "pending"); err != nil {
		return fmt.Errorf("failed to update URLs status: %w", err)
	}

	// Start crawling for each URL
	for _, url := range urls {
		go func(u string) {
			if err := s.crawler.CrawlURL(userID, u, s.repo); err != nil {
				log.Printf("Failed to crawl URL %s: %v", u, err)
			}
		}(url)
	}

	return nil
}

// BulkDelete deletes multiple crawl results
func (s *crawlService) BulkDelete(userID int, urls []string) error {
	if err := s.repo.BulkDelete(userID, urls); err != nil {
		return fmt.Errorf("failed to delete URLs: %w", err)
	}

	return nil
}

// StopCrawl stops crawling for a specific URL
func (s *crawlService) StopCrawl(userID int, id string) error {
	// This would need to be implemented in the repository
	// For now, we'll update the status to stopped
	crawlID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Get the URL first, then update its status
	urlData, err := s.repo.GetCrawlResultByID(userID, crawlID)
	if err != nil {
		return fmt.Errorf("failed to get crawl result: %w", err)
	}

	if err := s.repo.UpdateCrawlResultStatus(userID, urlData.URL, "stopped"); err != nil {
		return fmt.Errorf("failed to stop crawling: %w", err)
	}

	return nil
} 