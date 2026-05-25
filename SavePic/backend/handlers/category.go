package handlers

import (
	"net/http"

	"savepic/database"
	"savepic/models"

	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order"`
	Sort      int    `json:"sort"` // 兼容旧字段
}

type categoryWithCount struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	Count     int64  `json:"count"`
}

// ListCategories 获取所有分类（含表情包数量）
func ListCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		fail(c, http.StatusInternalServerError, 500, "获取分类失败")
		return
	}

	result := make([]categoryWithCount, 0, len(categories))
	for _, cat := range categories {
		var count int64
		database.DB.Model(&models.Meme{}).Where("category_id = ?", cat.ID).Count(&count)
		result = append(result, categoryWithCount{
			ID:        cat.ID,
			Name:      cat.Name,
			SortOrder: cat.SortOrder,
			Count:     count,
		})
	}
	success(c, result)
}

// CreateCategory 创建新分类
func CreateCategory(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, 400, "参数无效：名称不能为空")
		return
	}

	sortOrder := req.SortOrder
	if sortOrder == 0 && req.Sort != 0 {
		sortOrder = req.Sort
	}

	category := models.Category{
		Name:      req.Name,
		SortOrder: sortOrder,
	}
	if err := database.DB.Create(&category).Error; err != nil {
		fail(c, http.StatusInternalServerError, 500, "创建分类失败")
		return
	}

	success(c, category)
}
