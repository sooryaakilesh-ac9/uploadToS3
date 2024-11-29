package quotes

import "github.com/lib/pq"

type GoogleSheetsLink struct {
	GoogleSheetsLink string `json:"googleSheetsLink"`
}

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

// Quote struct represents a quote and its attributes
type Quote struct {
	Id   int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Text string         `json:"text"`
	Tags pq.StringArray `gorm:"type:text[];column:tags" json:"tags"` // Explicitly use pq.StringArray
	Lang string         `json:"lang"`
}

type QuotesMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalQuotes int    `json:"totalQuotes"`
	Url         string `json:"url"`
	Schema      Schema `json:"schema"`
}

type Schema struct {
	Format   string `json:"format"`
	Encoding string `json:"encoding"`
	FileType string `json:"fileType"`
}
