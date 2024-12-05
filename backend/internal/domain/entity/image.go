package entity

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
    FileName    string     `json:"fileName"`
}

type Resolution struct {
    Width  int `json:"width" gorm:"column:width"`
    Height int `json:"height" gorm:"column:height"`
    Unit   int `json:"unit" gorm:"column:unit"`
} 