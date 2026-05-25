package models

import "time"

// Category 表情包分类（文件夹）
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Memes []Meme `gorm:"foreignKey:CategoryID" json:"memes,omitempty"`
}
