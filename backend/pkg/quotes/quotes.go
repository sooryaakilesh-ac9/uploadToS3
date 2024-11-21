package quotes

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

type Quote struct {
	Id   int      `json:"id"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
	Lang string   `json:"lang"`
}

type QuotesMetadata struct {
	Version string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalQuotes int `json:"totalQuotes"`
	Url string `json:"url"`
	Schema Schema `json:"schema"`
}

type Schema struct {
	Format string `json:"format"`
	Encoding string `json:"encoding"`
	FileType string `json:"fileType"`
}
