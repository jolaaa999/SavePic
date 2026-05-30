package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	blobAPIVersion     = "12"
	defaultBlobAPIURL  = "https://vercel.com/api/blob"
	uploadsPathPrefix  = "uploads/"
)

type blobAuth struct {
	token   string
	storeID string
}

type blobPutResponse struct {
	URL string `json:"url"`
}

type blobDeleteRequest struct {
	URLs []string `json:"urls"`
}

func isVercel() bool {
	return strings.TrimSpace(os.Getenv("VERCEL")) == "1"
}

// IsLocal 未配置 Vercel Blob 时使用本地 uploads 目录（Vercel 上始终为 false）
func IsLocal() bool {
	if isVercel() {
		return false
	}
	return blobCredentials() == nil
}

// Save 保存图片并返回可访问 URL（本地为 /uploads/...，云端为 Blob 公网 URL）
func Save(data []byte, ext string) (string, error) {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	filename := uploadsPathPrefix + uuid.New().String() + ext

	if creds := blobCredentials(); creds != nil {
		return saveBlob(filename, data, creds)
	}
	if isVercel() {
		return "", fmt.Errorf("未配置 Blob 存储：请在 Vercel 项目中创建 Blob 并关联本项目（需要 BLOB_STORE_ID + VERCEL_OIDC_TOKEN，或 BLOB_READ_WRITE_TOKEN）")
	}
	return saveLocal(strings.TrimPrefix(filename, uploadsPathPrefix), data)
}

func saveLocal(filename string, data []byte) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}
	path := filepath.Join("uploads", filename)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return "/uploads/" + filename, nil
}

/**
 * blobCredentials 解析 Blob 认证，优先级与 @vercel/blob 一致：
 * 1. VERCEL_OIDC_TOKEN + BLOB_STORE_ID
 * 2. BLOB_READ_WRITE_TOKEN
 */
func blobCredentials() *blobAuth {
	oidc := strings.TrimSpace(os.Getenv("VERCEL_OIDC_TOKEN"))
	if oidc != "" {
		storeID := normalizeStoreID(strings.TrimSpace(os.Getenv("BLOB_STORE_ID")))
		if storeID != "" {
			return &blobAuth{token: oidc, storeID: storeID}
		}
	}

	rw := strings.TrimSpace(os.Getenv("BLOB_READ_WRITE_TOKEN"))
	if rw != "" {
		storeID := parseStoreIDFromReadWriteToken(rw)
		if storeID != "" {
			return &blobAuth{token: rw, storeID: storeID}
		}
	}
	return nil
}

func normalizeStoreID(storeID string) string {
	return strings.TrimPrefix(storeID, "store_")
}

func parseStoreIDFromReadWriteToken(token string) string {
	parts := strings.Split(token, "_")
	if len(parts) < 4 {
		return ""
	}
	return normalizeStoreID(parts[3])
}

func blobAPIBase() string {
	if u := strings.TrimSpace(os.Getenv("VERCEL_BLOB_API_URL")); u != "" {
		return strings.TrimRight(u, "/")
	}
	return defaultBlobAPIURL
}

func saveBlob(pathname string, data []byte, auth *blobAuth) (string, error) {
	apiURL := blobAPIBase() + "/?" + url.Values{"pathname": {pathname}}.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, apiURL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+auth.token)
	req.Header.Set("x-vercel-blob-store-id", auth.storeID)
	req.Header.Set("x-api-version", blobAPIVersion)
	req.Header.Set("x-content-length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("x-vercel-blob-access", "public")
	req.Header.Set("x-content-type", contentTypeForExt(filepath.Ext(pathname)))
	req.Header.Set("x-add-random-suffix", "0")
	req.Header.Set("x-allow-overwrite", "1")
	req.Header.Set("x-api-blob-request-id", fmt.Sprintf("%s:%d", auth.storeID, time.Now().UnixNano()))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("blob upload failed (%d): %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var out blobPutResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.URL == "" {
		return "", fmt.Errorf("blob upload returned empty url")
	}
	return out.URL, nil
}

// Delete 删除存储中的文件（按 URL 判断本地或 Blob）
func Delete(fileURL string) error {
	if fileURL == "" {
		return nil
	}
	if strings.HasPrefix(fileURL, "/uploads/") {
		return os.Remove(filepath.Join(".", fileURL))
	}

	auth := blobCredentials()
	if auth == nil || !strings.Contains(fileURL, "blob.vercel-storage.com") {
		return nil
	}

	apiURL := blobAPIBase() + "/delete"
	payload, _ := json.Marshal(blobDeleteRequest{URLs: []string{fileURL}})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, apiURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.token)
	req.Header.Set("x-vercel-blob-store-id", auth.storeID)
	req.Header.Set("x-api-version", blobAPIVersion)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(res.Body)
	return fmt.Errorf("blob delete failed (%d): %s", res.StatusCode, strings.TrimSpace(string(body)))
}

func contentTypeForExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
