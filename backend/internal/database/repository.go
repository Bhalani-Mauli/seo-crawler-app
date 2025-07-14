package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/seo-crawler-app/internal/models"
)

// Repository defines the interface for database operations
type Repository interface {
	CreateCrawlResult(userID int, url string) (int, error)
	GetCrawlResultIDByURL(userID int, url string) (int, error)
	UpdateCrawlResultStatus(userID int, url, status string) error
	GetCrawlResultByID(userID int, id int) (*models.URLData, error)
	GetCrawlResults(userID int, page, pageSize int, status, search, sortBy, sortOrder string) ([]models.URLData, int, error)
	UpdateCrawlData(userID int, url string, data *models.CrawlData) error
	GetLinksByCrawlID(crawlID int) ([]models.LinkData, error)
	GetHeadingsByCrawlID(crawlID int) ([]models.HeadingData, error)
	CreateLink(crawlID int, link *models.LinkData) (int64, error)
	UpdateLinkStatus(linkID int64, statusCode int, isAccessible bool) error
	CreateHeading(crawlID int, heading *models.HeadingData) error
	BulkUpdateStatus(userID int, urls []string, status string) error
	BulkDelete(userID int, urls []string) error
}

// CrawlRepository implements the Repository interface
type CrawlRepository struct {
	conn *Connection
}

// NewCrawlRepository creates a new crawl repository
func NewCrawlRepository(conn *Connection) Repository {
	return &CrawlRepository{conn: conn}
}

func (r *CrawlRepository) CreateCrawlResult(userID int, url string) (int, error) {
	log.Printf("Creating crawl result for user_id: %d, url: %s", userID, url)
	
	result, err := r.conn.DB.Exec("INSERT INTO crawl_results (user_id, url, status) VALUES (?, ?, 'pending')", userID, url)
	if err != nil {
		log.Printf("Failed to create crawl result for user_id: %d, url: %s, error: %v", userID, url, err)
		return 0, fmt.Errorf("failed to create crawl result: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	log.Printf("Successfully created crawl result with id: %d for user_id: %d, url: %s", id, userID, url)
	
	// Verify the record was created with correct user_id
	var storedUserID int
	var storedURL string
	err = r.conn.DB.QueryRow("SELECT user_id, url FROM crawl_results WHERE id = ?", id).Scan(&storedUserID, &storedURL)
	if err != nil {
		log.Printf("Warning: Could not verify stored record for id %d: %v", id, err)
	} else {
		log.Printf("Verification: Stored record id=%d has user_id=%d, url=%s", id, storedUserID, storedURL)
	}
	
	return int(id), nil
}

func (r *CrawlRepository) GetCrawlResultIDByURL(userID int, url string) (int, error) {
	var id int
	err := r.conn.DB.QueryRow("SELECT id FROM crawl_results WHERE user_id = ? AND url = ?", userID, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to get crawl result ID for URL %s: %w", url, err)
	}
	return id, nil
}

func (r *CrawlRepository) UpdateCrawlResultStatus(userID int, url, status string) error {
	_, err := r.conn.DB.Exec("UPDATE crawl_results SET status = ?, updated_at = NOW() WHERE user_id = ? AND url = ?", status, userID, url)
	if err != nil {
		return fmt.Errorf("failed to update crawl result status: %w", err)
	}
	return nil
}

func (r *CrawlRepository) GetCrawlResultByID(userID int, id int) (*models.URLData, error) {
	var urlData models.URLData
	var headings []byte

	err := r.conn.DB.QueryRow(`
		SELECT id, url, html_version, title, headings, internal_links, external_links, 
			   inaccessible_links, has_login_form, status, created_at, updated_at 
		FROM crawl_results WHERE user_id = ? AND id = ?
	`, userID, id).Scan(
		&urlData.ID, &urlData.URL, &urlData.CrawlData.HTMLVersion, &urlData.CrawlData.Title,
		&headings, &urlData.CrawlData.InternalLinks, &urlData.CrawlData.ExternalLinks,
		&urlData.CrawlData.InaccessibleLinks, &urlData.CrawlData.HasLoginForm,
		&urlData.CrawlData.Status, &urlData.CrawlData.CreatedAt, &urlData.CrawlData.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get crawl result: %w", err)
	}

	if headings != nil {
		if err := json.Unmarshal(headings, &urlData.CrawlData.Headings); err != nil {
			log.Printf("Error unmarshaling headings: %v", err)
			urlData.CrawlData.Headings = make(map[string]int)
		}
	} else {
		urlData.CrawlData.Headings = make(map[string]int)
	}

	return &urlData, nil
}

func (r *CrawlRepository) GetCrawlResults(userID int, page, pageSize int, status, search, sortBy, sortOrder string) ([]models.URLData, int, error) {
	// Build query with filters
	whereClause := "WHERE user_id = ?"
	args := []interface{}{userID}

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if search != "" {
		whereClause += " AND (url LIKE ? OR title LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Validate sort column
	allowedSortColumns := map[string]bool{
		"url": true, "title": true, "status": true, "created_at": true, "updated_at": true,
	}
	if !allowedSortColumns[sortBy] {
		sortBy = "created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM crawl_results " + whereClause
	err := r.conn.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, url, html_version, title, headings, internal_links, external_links, 
			   inaccessible_links, has_login_form, status, created_at, updated_at 
		FROM crawl_results %s 
		ORDER BY %s %s 
		LIMIT ? OFFSET ?
	`, whereClause, sortBy, sortOrder)

	args = append(args, pageSize, offset)
	rows, err := r.conn.DB.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch results: %w", err)
	}
	defer rows.Close()

	var results []models.URLData
	for rows.Next() {
		var urlData models.URLData
		var headings []byte
		
		if err := rows.Scan(
			&urlData.ID, 
			&urlData.URL, 
			&urlData.CrawlData.HTMLVersion, 
			&urlData.CrawlData.Title,
			&headings, 
			&urlData.CrawlData.InternalLinks, 
			&urlData.CrawlData.ExternalLinks,
			&urlData.CrawlData.InaccessibleLinks, 
			&urlData.CrawlData.HasLoginForm,
			&urlData.CrawlData.Status, 
			&urlData.CrawlData.CreatedAt, 
			&urlData.CrawlData.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if headings != nil {
			if err := json.Unmarshal(headings, &urlData.CrawlData.Headings); err != nil {
				log.Printf("Error unmarshaling headings: %v", err)
				urlData.CrawlData.Headings = make(map[string]int)
			}
		} else {
			urlData.CrawlData.Headings = make(map[string]int)
		}

		results = append(results, urlData)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, total, nil
}

func (r *CrawlRepository) UpdateCrawlData(userID int, url string, data *models.CrawlData) error {
	headingsJSON, _ := json.Marshal(data.Headings)
	
	// Handle nullable values
	var htmlVersion, title interface{}
	if data.HTMLVersion.Valid {
		htmlVersion = data.HTMLVersion.String
	} else {
		htmlVersion = nil
	}
	if data.Title.Valid {
		title = data.Title.String
	} else {
		title = nil
	}
	
	_, err := r.conn.DB.Exec(`
		UPDATE crawl_results 
		SET html_version = ?, title = ?, headings = ?, internal_links = ?, 
			external_links = ?, inaccessible_links = ?, has_login_form = ?, 
			status = 'done', updated_at = NOW() 
		WHERE user_id = ? AND url = ?
	`, htmlVersion, title, headingsJSON, data.InternalLinks,
		data.ExternalLinks, data.InaccessibleLinks, data.HasLoginForm, userID, url)
	
	if err != nil {
		return fmt.Errorf("failed to update crawl data: %w", err)
	}
	return nil
}

func (r *CrawlRepository) GetLinksByCrawlID(crawlID int) ([]models.LinkData, error) {
	rows, err := r.conn.DB.Query(`
		SELECT id, link_url, link_text, link_type, status_code, is_accessible
		FROM crawl_links WHERE crawl_result_id = ? ORDER BY id
	`, crawlID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch links: %w", err)
	}
	defer rows.Close()

	var links []models.LinkData
	for rows.Next() {
		var link models.LinkData
		err := rows.Scan(&link.ID, &link.URL, &link.Text, &link.Type, &link.StatusCode, &link.IsAccessible)
		if err == nil {
			links = append(links, link)
		}
	}

	return links, nil
}

func (r *CrawlRepository) GetHeadingsByCrawlID(crawlID int) ([]models.HeadingData, error) {
	rows, err := r.conn.DB.Query(`
		SELECT id, heading_level, heading_text, heading_order
		FROM crawl_headings WHERE crawl_result_id = ? ORDER BY heading_order
	`, crawlID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch headings: %w", err)
	}
	defer rows.Close()

	var headings []models.HeadingData
	for rows.Next() {
		var heading models.HeadingData
		err := rows.Scan(&heading.ID, &heading.Level, &heading.Text, &heading.Order)
		if err == nil {
			headings = append(headings, heading)
		}
	}

	return headings, nil
}

func (r *CrawlRepository) CreateLink(crawlID int, link *models.LinkData) (int64, error) {
	result, err := r.conn.DB.Exec(`
		INSERT INTO crawl_links (crawl_result_id, link_url, link_text, link_type, status_code, is_accessible)
		VALUES (?, ?, ?, ?, ?, ?)
	`, crawlID, link.URL, link.Text, link.Type, link.StatusCode, link.IsAccessible)
	if err != nil {
		return 0, fmt.Errorf("failed to create link: %w", err)
	}
	
	return result.LastInsertId()
}

func (r *CrawlRepository) UpdateLinkStatus(linkID int64, statusCode int, isAccessible bool) error {
	_, err := r.conn.DB.Exec(`
		UPDATE crawl_links 
		SET status_code = ?, is_accessible = ? 
		WHERE id = ?
	`, statusCode, isAccessible, linkID)
	if err != nil {
		return fmt.Errorf("failed to update link status: %w", err)
	}
	return nil
}

func (r *CrawlRepository) CreateHeading(crawlID int, heading *models.HeadingData) error {
	_, err := r.conn.DB.Exec(`
		INSERT INTO crawl_headings (crawl_result_id, heading_level, heading_text, heading_order)
		VALUES (?, ?, ?, ?)
	`, crawlID, heading.Level, heading.Text, heading.Order)
	if err != nil {
		return fmt.Errorf("failed to create heading: %w", err)
	}
	return nil
}

func (r *CrawlRepository) BulkUpdateStatus(userID int, urls []string, status string) error {
	placeholders := strings.Repeat("?,", len(urls)-1) + "?"
	query := fmt.Sprintf("UPDATE crawl_results SET status = ?, updated_at = NOW() WHERE user_id = ? AND url IN (%s)", placeholders)
	
	args := make([]interface{}, len(urls)+2)
	args[0] = status
	args[1] = userID
	for i, url := range urls {
		args[i+2] = url
	}

	_, err := r.conn.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk update status: %w", err)
	}
	return nil
}

func (r *CrawlRepository) BulkDelete(userID int, urls []string) error {
	placeholders := strings.Repeat("?,", len(urls)-1) + "?"
	query := fmt.Sprintf("DELETE FROM crawl_results WHERE user_id = ? AND url IN (%s)", placeholders)
	
	args := make([]interface{}, len(urls)+1)
	args[0] = userID
	for i, url := range urls {
		args[i+1] = url
	}

	_, err := r.conn.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk delete: %w", err)
	}
	return nil
} 