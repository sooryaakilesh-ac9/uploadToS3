package handler

import (
	"backend/pkg/images"
	"backend/pkg/quotes"
)

type MetadataFlyers struct {
	Flyers         []Flyer        `json:"flyers"`
	FlyersMetadata FlyersMetadata `json:"metadata"`
}

type Flyer struct {
	Id     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Design Design `json:"design" gorm:"embedded"`
	Lang   string `json:"lang"`
	Url    string `json:"url"`
}

type Design struct {
	TemplateId  string     `json:"templateId" gorm:"column:template_id"`
	Resolution  Resolution `json:"resolution" gorm:"embedded"`
	Type        string     `json:"type"`
	Tags        []string   `json:"tags" gorm:"serializer:json"`
	FileFormat  string     `json:"fileFormat" gorm:"column:file_format"`
	Orientation string     `json:"orientation"`
}

type Resolution struct {
	Width  int `json:"width" gorm:"column:width"`
	Height int `json:"height" gorm:"column:height"`
	Unit   int `json:"unit" gorm:"column:unit"`
}

type FlyersMetadata struct {
	Version     string `json:"version"`
	LastUpdated string `json:"lastUpdated"`
	TotalFlyers int    `json:"total"`
	Url         string `json:"url"`
	Schema      Schema `json:"schema,omitempty"`
}

type MetaDataImages struct {
	Images         []images.Flyer `json:"images"`
	ImagesMetadata FlyersMetadata `json:"metadata"`
}

type MetaDataQuotes struct {
	Quotes         []quotes.Quote `json:"quotes"`
	QuotesMetadata QuotesMetadata `json:"quotesMetadata"`
}
