package images

type Flyer struct {
	Id string `json:"id"`
	Design Design `json:"design"`
	Lang string `json:"lang"`
	Url string `json:"url"`
}

type Design struct {
	TemplateId string `json:"templateId"`
	Resolution Resolution `json:"resolution"`
	Type string `json:"type"`
	Tags []string `json:"tags"`
	FileFormat string `json:"fileFormat"`
	Orientation string `json:"orientation"`
}

type Resolution struct {
	Width int `json:"width"`
	Height int `json:"height"`
	Unit int `json:"unit"`
}

type FlyersMetadata struct {
	Version string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalFlyers int `json:"totalFlyers"`
	Url string `json:"url"`
	Schema Schema `json:"schema"`
}

type Schema struct {
	Format string `json:"format"`
	Encoding string `json:"encoding"`
	FileType string `json:"fileType"`
}