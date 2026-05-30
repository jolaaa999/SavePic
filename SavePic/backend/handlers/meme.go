package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"savepic/backend/database"
	"savepic/backend/internal/imgutil"
	"savepic/backend/models"
	"savepic/backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var allowedImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

type updateMemeRequest struct {
	CategoryID *uint    `json:"category_id"`
	Tags       []string `json:"tags"`
}

// UploadMeme 上传表情包（MD5 去重 + 尺寸 + 标签）
func UploadMeme(c *gin.Context) {
	categoryIDStr := c.PostForm("category_id")
	if categoryIDStr == "" {
		fail(c, http.StatusBadRequest, 400, "缺少 category_id")
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, 400, "category_id 格式无效")
		return
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		fail(c, http.StatusNotFound, 404, "分类不存在")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		fail(c, http.StatusBadRequest, 400, "请上传图片文件")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		fail(c, http.StatusBadRequest, 400, "仅支持 jpg、jpeg、png、gif、webp 格式")
		return
	}

	tmpPath := filepath.Join(os.TempDir(), uuid.New().String()+ext)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		fail(c, http.StatusInternalServerError, 500, "保存临时文件失败")
		return
	}
	defer os.Remove(tmpPath)

	fileHash, err := imgutil.FileMD5(tmpPath)
	if err != nil {
		fail(c, http.StatusInternalServerError, 500, "计算文件哈希失败")
		return
	}

	tagNames := parseTagNamesFromForm(c.PostForm("tags"))

	var existing models.Meme
	if err := database.DB.Where("file_hash = ?", fileHash).Preload("Tags").First(&existing).Error; err == nil {
		if len(tagNames) > 0 {
			_ = syncMemeTags(&existing, tagNames)
			database.DB.Preload("Tags").First(&existing, existing.ID)
		}
		success(c, gin.H{"meme": existing, "duplicate": true})
		return
	} else if err != gorm.ErrRecordNotFound {
		fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}

	width, height, _ := imgutil.Dimensions(tmpPath)
	info, _ := os.Stat(tmpPath)
	fileSize := int64(0)
	if info != nil {
		fileSize = info.Size()
	}

	data, readErr := os.ReadFile(tmpPath)
	if readErr != nil {
		fail(c, http.StatusInternalServerError, 500, "读取文件失败")
		return
	}
	fileURL, saveErr := storage.Save(c.Request.Context(), data, ext)
	if saveErr != nil {
		msg := "保存文件失败"
		if isVercel := os.Getenv("VERCEL") == "1"; isVercel {
			msg = saveErr.Error()
		}
		fail(c, http.StatusInternalServerError, 500, msg)
		return
	}

	meme := models.Meme{
		CategoryID: uint(categoryID),
		FileURL:    fileURL,
		FileHash:   fileHash,
		Width:      width,
		Height:     height,
		Size:       fileSize,
	}
	if err := database.DB.Create(&meme).Error; err != nil {
		_ = storage.Delete(c.Request.Context(), fileURL)
		fail(c, http.StatusInternalServerError, 500, "保存记录失败")
		return
	}

	if len(tagNames) > 0 {
		if err := syncMemeTags(&meme, tagNames); err != nil {
			fail(c, http.StatusInternalServerError, 500, "保存标签失败")
			return
		}
	}

	database.DB.Preload("Tags").First(&meme, meme.ID)
	success(c, gin.H{"meme": meme, "duplicate": false})
}

// ListMemes 复杂查询：category_id + tag_ids（交叉 AND）+ sort
func ListMemes(c *gin.Context) {
	query := database.DB.Model(&models.Meme{}).Preload("Tags")

	if cid := c.Query("category_id"); cid != "" {
		query = query.Where("category_id = ?", cid)
	}

	tagIDs := parseUintList(c.Query("tag_ids"))
	if len(tagIDs) > 0 {
		query = query.Where(
			"id IN (?)",
			database.DB.Table("meme_tags").
				Select("meme_id").
				Where("tag_id IN ?", tagIDs).
				Group("meme_id").
				Having("COUNT(DISTINCT tag_id) = ?", len(tagIDs)),
		)
	}

	sort := c.DefaultQuery("sort", "desc")
	if sort == "asc" {
		query = query.Order("created_at ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	var memes []models.Meme
	if err := query.Find(&memes).Error; err != nil {
		fail(c, http.StatusInternalServerError, 500, "获取表情包失败")
		return
	}

	success(c, memes)
}

// ListMemesByCategory 兼容旧路由
func ListMemesByCategory(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Set("category_id", c.Param("id"))
	if q.Get("sort") == "" {
		q.Set("sort", "desc")
	}
	c.Request.URL.RawQuery = q.Encode()
	ListMemes(c)
}

// UpdateMeme 更新分类或标签
func UpdateMeme(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, 400, "表情包 ID 无效")
		return
	}

	var meme models.Meme
	if err := database.DB.First(&meme, id).Error; err != nil {
		fail(c, http.StatusNotFound, 404, "表情包不存在")
		return
	}

	var req updateMemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, 400, "参数无效")
		return
	}

	if req.CategoryID != nil {
		var cat models.Category
		if err := database.DB.First(&cat, *req.CategoryID).Error; err != nil {
			fail(c, http.StatusNotFound, 404, "分类不存在")
			return
		}
		meme.CategoryID = *req.CategoryID
		if err := database.DB.Save(&meme).Error; err != nil {
			fail(c, http.StatusInternalServerError, 500, "更新失败")
			return
		}
	}

	if req.Tags != nil {
		if err := syncMemeTags(&meme, req.Tags); err != nil {
			fail(c, http.StatusInternalServerError, 500, "更新标签失败")
			return
		}
	}

	database.DB.Preload("Tags").First(&meme, meme.ID)
	success(c, meme)
}

// deleteMemeRecord 删除表情包记录、关联标签及存储文件
func deleteMemeRecord(ctx context.Context, meme *models.Meme) error {
	_ = storage.Delete(ctx, meme.FileURL)
	_ = database.DB.Model(meme).Association("Tags").Clear()
	return database.DB.Delete(meme).Error
}

// DeleteMeme 删除表情包
func DeleteMeme(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, 400, "表情包 ID 无效")
		return
	}

	var meme models.Meme
	if err := database.DB.First(&meme, id).Error; err != nil {
		fail(c, http.StatusNotFound, 404, "表情包不存在")
		return
	}

	if err := deleteMemeRecord(c.Request.Context(), &meme); err != nil {
		fail(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}

	success(c, nil)
}

func parseUintList(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseUint(p, 10, 64)
		if err == nil {
			ids = append(ids, uint(n))
		}
	}
	return ids
}
