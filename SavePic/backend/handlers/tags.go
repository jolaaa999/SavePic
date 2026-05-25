package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"savepic/database"
	"savepic/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tagWithCount struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// ListTags 获取所有标签（含使用次数）
func ListTags(c *gin.Context) {
	var tags []models.Tag
	if err := database.DB.Order("name ASC").Find(&tags).Error; err != nil {
		fail(c, http.StatusInternalServerError, 500, "获取标签失败")
		return
	}

	result := make([]tagWithCount, 0, len(tags))
	for _, tag := range tags {
		var count int64
		database.DB.Table("meme_tags").Where("tag_id = ?", tag.ID).Count(&count)
		result = append(result, tagWithCount{
			ID:    tag.ID,
			Name:  tag.Name,
			Count: count,
		})
	}
	success(c, result)
}

func normalizeTagNames(names []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(names))
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	return out
}

func parseTagNamesFromForm(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var arr []string
	if json.Unmarshal([]byte(raw), &arr) == nil {
		return normalizeTagNames(arr)
	}
	return normalizeTagNames(strings.Split(raw, ","))
}

// syncMemeTags 按名称创建标签并关联到表情包（替换原有关联）
func syncMemeTags(meme *models.Meme, names []string) error {
	names = normalizeTagNames(names)
	tags := make([]models.Tag, 0, len(names))

	for _, name := range names {
		var tag models.Tag
		err := database.DB.Where("name = ?", name).First(&tag).Error
		if err == gorm.ErrRecordNotFound {
			tag = models.Tag{Name: name}
			if err := database.DB.Create(&tag).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	return database.DB.Model(meme).Association("Tags").Replace(tags)
}
