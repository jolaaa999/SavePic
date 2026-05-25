package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const blobAPIVersion = "12"

// IsLocal 未配置 Vercel Blob 时使用本地 uploads 目录
func IsLocal() bool {
	return strings.TrimSpace(os.Getenv("BLOB_READ_WRITE_TOKEN")) == ""
}

// Save 保存图片并返回可访问 URL（本地为 /uploads/...，云端为 Blob 公网 URL）
func Save(data []byte, ext string) (string, error) {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	filename := uuid.New().String() + ext

	token := strings.TrimSpace(os.Getenv("BLOB_READ_WRITE_TOKEN"))
	if token == "" {
		return saveLocal(filename, data)
	}
	return saveBlob(filename, data, token)
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

func saveBlob(filename string, data []byte, token string) (string, error) {
	storeID := parseStoreIDFromToken(token)
	if storeID == "" {
		return "", fmt.Errorf("invalid BLOB_READ_WRITE_TOKEN: cannot parse store id")
	}

	apiURL := os.Getenv("VERCEL_BLOB_API_URL")
	if strings.TrimSpace(apiURL) == "" {
		apiURL = "https://vercel.com/api/blob"
	}

	reqURL := apiURL + "/?" + url.Values{"pathname": {filename}}.Encode()
	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-vercel-blob-store-id", storeID)
	req.Header.Set("x-api-version", blobAPIVersion)
	req.Header.Set("x-content-length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("x-vercel-blob-access", "public")
	req.Header.Set("x-content-type", contentTypeForExt(filepath.Ext(filename)))
	req.Header.Set("x-add-random-suffix", "false")
	req.Header.Set("x-allow-overwrite", "true")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("blob upload failed (%d): %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var out struct {
		URL string `json:"url"`
	}
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
	token := strings.TrimSpace(os.Getenv("BLOB_READ_WRITE_TOKEN"))
	if token == "" || !strings.Contains(fileURL, "blob.vercel-storage.com") {
		return nil
	}

	storeID := parseStoreIDFromToken(token)
	apiURL := os.Getenv("VERCEL_BLOB_API_URL")
	if strings.TrimSpace(apiURL) == "" {
		apiURL = "https://vercel.com/api/blob"
	}

	u, err := url.Parse(fileURL)
	if err != nil {
		return err
	}
	pathname := strings.TrimPrefix(u.Path, "/")
	reqURL := apiURL + "/?" + url.Values{"pathname": {pathname}}.Encode()

	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-vercel-blob-store-id", storeID)
	req.Header.Set("x-api-version", blobAPIVersion)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("blob delete failed: %d", res.StatusCode)
}

func parseStoreIDFromToken(token string) string {
	if id := strings.TrimSpace(os.Getenv("BLOB_STORE_ID")); id != "" {
		return strings.TrimPrefix(id, "store_")
	}
	parts := strings.Split(token, "_")
	if len(parts) < 4 {
		return ""
	}
	storeID := parts[3]
	if strings.HasPrefix(storeID, "store_") {
		return strings.TrimPrefix(storeID, "store_")
	}
	return storeID
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
