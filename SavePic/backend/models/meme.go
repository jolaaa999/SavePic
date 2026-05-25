package models

import "time"

// Meme 表情包
type Meme struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"index;not null" json:"category_id"`
	FileURL    string    `gorm:"size:255;not null" json:"file_url"`
	FileHash   string    `gorm:"size:64;uniqueIndex;not null" json:"file_hash"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag    `gorm:"many2many:meme_tags" json:"tags,omitempty"`
}
