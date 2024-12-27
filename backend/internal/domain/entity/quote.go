package entity

type Quote struct {
    Id   int      `json:"id" gorm:"primaryKey;autoIncrement"`
    Text string   `json:"text"`
    Tags []string `json:"tags" gorm:"serializer:json"`
    Lang string   `json:"lang"`
} 