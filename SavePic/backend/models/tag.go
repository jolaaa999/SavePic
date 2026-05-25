package models

import "time"

// Tag 标签（交叉关键词）
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Memes []Meme `gorm:"many2many:meme_tags" json:"memes,omitempty"`
}
