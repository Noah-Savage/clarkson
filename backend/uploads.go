package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Attachment struct {
	ID        uint   `gorm:"primaryKey"`
	EntryID   uint   `gorm:"not null"`
	EntryType string `gorm:"not null"`
	Filename  string `gorm:"not null"`
	Path      string `gorm:"not null"`
}

func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}

type Application struct {
	db *gorm.DB
}

func (app *Application) compressImage(inputPath string, outputPath string, quality int) error {
	// Open the image
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image based on extension
	ext := strings.ToLower(filepath.Ext(inputPath))
	var img image.Image

	if ext == ".jpg" || ext == ".jpeg" {
		img, err = jpeg.Decode(file)
	} else if ext == ".png" {
		img, err = png.Decode(file)
	} else {
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return err
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode with compression
	if ext == ".jpg" || ext == ".jpeg" {
		opts := &jpeg.Options{Quality: quality}
		return jpeg.Encode(outFile, img, opts)
	} else if ext == ".png" {
		return png.Encode(outFile, img)
	}

	return nil
}

func (app *Application) handleFileUpload(c *gin.Context) {
	entryType := c.PostForm("entry_type") // fuel or expense
	entryID := c.PostForm("entry_id")

	if entryType == "" || entryID == "" {
		c.JSON(400, gin.H{"error": "Missing entry_type or entry_id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(400, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/webp":      true,
		"application/pdf": true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(400, gin.H{"error": "Invalid file type"})
		return
	}

	assetsPath := os.Getenv("ASSETS_PATH")
	if assetsPath == "" {
		assetsPath = "/assets"
	}

	os.MkdirAll(assetsPath, 0755)

	// Generate unique filename
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	filepath := filepath.Join(assetsPath, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "Upload failed"})
		return
	}

	// Compress image if it's a photo
	if strings.Contains(file.Header.Get("Content-Type"), "image") {
		compressedPath := filepath + ".compressed"
		if err := app.compressImage(filepath, compressedPath, 80); err == nil {
			// Replace original with compressed
			os.Remove(filepath)
			os.Rename(compressedPath, filepath)
		}
	}

	// Create attachment record
	attachment := Attachment{
		EntryID:   parseUint(entryID),
		EntryType: entryType,
		Filename:  file.Filename,
		Path:      filepath,
	}

	if err := app.db.Create(&attachment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(201, gin.H{
		"attachment": attachment,
		"url":        fmt.Sprintf("/api/download/%s", filename),
	})
}

func (app *Application) handleFileDownload(c *gin.Context) {
	filename := c.Param("id")

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(400, gin.H{"error": "Invalid filename"})
		return
	}

	assetsPath := os.Getenv("ASSETS_PATH")
	if assetsPath == "" {
		assetsPath = "/assets"
	}

	filepath := filepath.Join(assetsPath, filename)

	// Verify file exists and is in assets directory
	if _, err := os.Stat(filepath); err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	c.File(filepath)
}

func (app *Application) deleteAttachment(c *gin.Context) {
	attachmentID := c.Param("id")

	var attachment Attachment
	if err := app.db.First(&attachment, attachmentID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Attachment not found"})
		return
	}

	// Delete file from disk
	os.Remove(attachment.Path)

	// Delete database record
	if err := app.db.Delete(&attachment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Attachment deleted"})
}

func (app *Application) getAttachments(c *gin.Context) {
	entryType := c.Query("type")
	entryID := c.Query("entry_id")

	if entryType == "" || entryID == "" {
		c.JSON(400, gin.H{"error": "Missing parameters"})
		return
	}

	var attachments []Attachment
	if err := app.db.
		Where("entry_type = ? AND entry_id = ?", entryType, entryID).
		Find(&attachments).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, attachments)
}
