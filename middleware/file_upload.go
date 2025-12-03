package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// FileUploadConfig holds file upload configuration
type FileUploadConfig struct {
	MaxFileSize    int64    // Maximum file size in bytes
	AllowedTypes   []string // Allowed file extensions
	UploadPath     string   // Upload directory path
	CreateDirIfNot bool     // Create upload directory if not exists
}

// DefaultFileUploadConfig returns default configuration
func DefaultFileUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		MaxFileSize:    10485760, // 10MB
		AllowedTypes:   []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"},
		UploadPath:     "./uploads",
		CreateDirIfNot: true,
	}
}

// FileUploadMiddleware creates a middleware for handling file uploads
func FileUploadMiddleware(config ...FileUploadConfig) fiber.Handler {
	cfg := DefaultFileUploadConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Create upload directory if not exists
	if cfg.CreateDirIfNot {
		if err := os.MkdirAll(cfg.UploadPath, os.ModePerm); err != nil {
			utils.GlobalLogger.Error("Failed to create upload directory", err, map[string]interface{}{
				"path": cfg.UploadPath,
			})
		}
	}

	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// UploadFile handles single file upload
func UploadFile(c *fiber.Ctx, formFieldName string, config ...FileUploadConfig) (string, error) {
	cfg := DefaultFileUploadConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Get file from form
	file, err := c.FormFile(formFieldName)
	if err != nil {
		return "", fmt.Errorf("no file uploaded: %v", err)
	}

	// Check file size
	if file.Size > cfg.MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", cfg.MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedExt := range cfg.AllowedTypes {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("file type %s not allowed. Allowed types: %v", ext, cfg.AllowedTypes)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(cfg.UploadPath, filename)

	// Save file
	if err := c.SaveFile(file, filepath); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	utils.GlobalLogger.Info("File uploaded successfully", map[string]interface{}{
		"original_filename": file.Filename,
		"saved_filename":    filename,
		"size":              file.Size,
		"path":              filepath,
	})

	return filename, nil
}

// UploadMultipleFiles handles multiple file uploads
func UploadMultipleFiles(c *fiber.Ctx, formFieldName string, config ...FileUploadConfig) ([]string, error) {
	cfg := DefaultFileUploadConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %v", err)
	}

	files := form.File[formFieldName]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files uploaded")
	}

	var uploadedFiles []string

	for _, file := range files {
		// Check file size
		if file.Size > cfg.MaxFileSize {
			return uploadedFiles, fmt.Errorf("file %s exceeds maximum size", file.Filename)
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowed := false
		for _, allowedExt := range cfg.AllowedTypes {
			if ext == allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return uploadedFiles, fmt.Errorf("file type %s not allowed for %s", ext, file.Filename)
		}

		// Generate unique filename
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath := filepath.Join(cfg.UploadPath, filename)

		// Save file
		if err := c.SaveFile(file, savePath); err != nil {
			return uploadedFiles, fmt.Errorf("failed to save file %s: %v", file.Filename, err)
		}

		uploadedFiles = append(uploadedFiles, filename)

		utils.GlobalLogger.Info("File uploaded", map[string]interface{}{
			"original_filename": file.Filename,
			"saved_filename":    filename,
			"size":              file.Size,
		})
	}

	return uploadedFiles, nil
}

// DeleteFile deletes uploaded file
func DeleteFile(filename string, uploadPath ...string) error {
	path := "./uploads"
	if len(uploadPath) > 0 {
		path = uploadPath[0]
	}

	filepath := filepath.Join(path, filename)
	
	if err := os.Remove(filepath); err != nil {
		utils.GlobalLogger.Error("Failed to delete file", err, map[string]interface{}{
			"filename": filename,
			"path":     filepath,
		})
		return err
	}

	utils.GlobalLogger.Info("File deleted", map[string]interface{}{
		"filename": filename,
		"path":     filepath,
	})

	return nil
}
