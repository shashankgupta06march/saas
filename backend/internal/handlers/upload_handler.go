package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles image uploads (e.g. the chatbot launcher icon) and
// stores them on disk so they can be served back as static files.
type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	// Best-effort: make sure the directory exists at startup.
	_ = os.MkdirAll(uploadDir, 0o755)
	return &UploadHandler{uploadDir: uploadDir}
}

// allowedImageExts maps accepted upload extensions to their canonical form.
var allowedImageExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".webp": true,
}

const maxImageUploadBytes = 3 << 20 // 3 MB

// UploadImage accepts a multipart "file" field, validates it's an allowed image
// type within the size limit, stores it under the uploads directory and returns
// an absolute URL the widget can load from any domain.
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if header.Size > maxImageUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is too large (max 3 MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image type. Allowed: PNG, JPEG, GIF, WebP"})
		return
	}

	// Generate a random, collision-free filename; never trust the client's name.
	nameBytes := make([]byte, 16)
	if _, err := rand.Read(nameBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate filename"})
		return
	}
	filename := hex.EncodeToString(nameBytes) + ext
	destPath := filepath.Join(h.uploadDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, io.LimitReader(file, maxImageUploadBytes+1)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	publicURL := fmt.Sprintf("%s/uploads/%s", baseURL(c), filename)
	c.JSON(http.StatusCreated, gin.H{
		"url":      publicURL,
		"filename": filename,
	})
}

// baseURL reconstructs the externally-visible base URL (scheme + host) from the
// request, honouring reverse-proxy headers so uploaded URLs work in production.
func baseURL(c *gin.Context) string {
	scheme := "http"
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if c.Request.TLS != nil {
		scheme = "https"
	}

	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	return scheme + "://" + host
}
