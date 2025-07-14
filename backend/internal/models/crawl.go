package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// LinkData represents a single link found on the page
type LinkData struct {
	ID           int    `json:"id"`
	URL          string `json:"url"`
	Text         string `json:"text"`
	Type         string `json:"type"` // "internal" or "external"
	StatusCode   int    `json:"status_code"`
	IsAccessible bool   `json:"is_accessible"`
}

// HeadingData represents a single heading found on the page
type HeadingData struct {
	ID    int    `json:"id"`
	Level string `json:"level"` // "h1", "h2", etc.
	Text  string `json:"text"`
	Order int    `json:"order"`
}

// CrawlData represents the data collected for a crawled URL.
type CrawlData struct {
	HTMLVersion       sql.NullString `json:"-"`
	Title             sql.NullString `json:"-"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internal_links"`
	ExternalLinks     int            `json:"external_links"`
	InaccessibleLinks int            `json:"inaccessible_links"`
	HasLoginForm      bool           `json:"has_login_form"`
	Status            string         `json:"status"`
	CreatedAt         sql.NullTime   `json:"-"`
	UpdatedAt         sql.NullTime   `json:"-"`
	
	// JSON fields
	HTMLVersionStr    string         `json:"html_version"`
	TitleStr          string         `json:"title"`
	CreatedAtStr      string         `json:"created_at"`
	UpdatedAtStr      string         `json:"updated_at"`
	
	// Detailed data
	Links             []LinkData     `json:"links,omitempty"`
	HeadingDetails    []HeadingData  `json:"heading_details,omitempty"`
}

// MarshalJSON implements custom JSON marshaling
func (c CrawlData) MarshalJSON() ([]byte, error) {
	type Alias CrawlData
	
	// Convert nullable fields to strings
	if c.HTMLVersion.Valid {
		c.HTMLVersionStr = c.HTMLVersion.String
	}
	if c.Title.Valid {
		c.TitleStr = c.Title.String
	}
	if c.CreatedAt.Valid {
		c.CreatedAtStr = c.CreatedAt.Time.Format(time.RFC3339)
	}
	if c.UpdatedAt.Valid {
		c.UpdatedAtStr = c.UpdatedAt.Time.Format(time.RFC3339)
	}
	
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&c),
	})
}

// URLData stores the URL and the collected crawl data.
type URLData struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	CrawlData CrawlData `json:"crawl_data"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Data       []URLData `json:"data"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// CrawlRequest represents a crawl request
type CrawlRequest struct {
	URL string `json:"url" binding:"required"`
}

// BulkActionRequest represents bulk action requests
type BulkActionRequest struct {
	URLs []string `json:"urls" binding:"required"`
} 