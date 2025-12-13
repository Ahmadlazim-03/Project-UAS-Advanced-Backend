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

// UploadFile handles file upload with validation
func UploadFile(c *fiber.Ctx, fieldName string) (string, error) {
	// Get file from request
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %v", err)
	}

	// Check file size (max 10MB)
	maxSize := int64(10485760)
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds maximum limit of 10MB")
	}

	// Check file extension
	allowedExt := []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	isAllowed := false
	for _, allowed := range allowedExt {
		if ext == allowed {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return "", fmt.Errorf("file type not allowed. Allowed types: PDF, JPG, PNG, DOC, DOCX")
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Create upload directory if not exists
	uploadPath := "./uploads/achievements"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		utils.GlobalLogger.Error("Failed to create upload directory", err, nil)
		return "", fmt.Errorf("failed to create upload directory")
	}

	// Save file
	filePath := filepath.Join(uploadPath, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		utils.GlobalLogger.Error("Failed to save file", err, map[string]interface{}{
			"filename": filename,
		})
		return "", fmt.Errorf("failed to save file")
	}

	return filename, nil
}

// DeleteFile deletes a file from the uploads directory
func DeleteFile(filename string) error {
	filePath := filepath.Join("./uploads/achievements", filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		utils.GlobalLogger.Error("Failed to delete file", err, map[string]interface{}{
			"filename": filename,
		})
		return fmt.Errorf("failed to delete file")
	}

	utils.GlobalLogger.Info("File deleted successfully", map[string]interface{}{
		"filename": filename,
	})

	return nil
}
