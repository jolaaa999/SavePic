package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"savepic/backend/database"
	"savepic/backend/models"

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

type updateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTag 修改标签名称
func UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, 400, "标签 ID 无效")
		return
	}

	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, 400, "参数无效")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		fail(c, http.StatusBadRequest, 400, "标签名不能为空")
		return
	}

	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		fail(c, http.StatusNotFound, 404, "标签不存在")
		return
	}

	if tag.Name == name {
		success(c, tag)
		return
	}

	var existing models.Tag
	if err := database.DB.Where("name = ? AND id != ?", name, tag.ID).First(&existing).Error; err == nil {
		fail(c, http.StatusConflict, 409, "标签名已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}

	tag.Name = name
	if err := database.DB.Save(&tag).Error; err != nil {
		fail(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}
	success(c, tag)
}

// DeleteTag 删除标签并解除与所有表情包的关联
func DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, 400, "标签 ID 无效")
		return
	}

	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		fail(c, http.StatusNotFound, 404, "标签不存在")
		return
	}

	if err := database.DB.Model(&tag).Association("Memes").Clear(); err != nil {
		fail(c, http.StatusInternalServerError, 500, "解除关联失败")
		return
	}
	if err := database.DB.Delete(&tag).Error; err != nil {
		fail(c, http.StatusInternalServerError, 500, "删除标签失败")
		return
	}

	success(c, nil)
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
