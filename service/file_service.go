package service

import (
	"student-achievement-system/middleware"
	"student-achievement-system/utils"

	"github.com/gofiber/fiber/v2"
)

type FileService interface {
	UploadAchievementFile(c *fiber.Ctx) error
	DeleteAchievementFile(c *fiber.Ctx) error
}

type fileService struct{}

func NewFileService() FileService {
	return &fileService{}
}

// UploadAchievementFile godoc
// @Summary      Upload achievement attachment
// @Description  Upload a file attachment for an achievement (PDF, JPG, PNG, DOC, DOCX)
// @Tags         Files
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "File to upload"
// @Success      200 {object} map[string]interface{} "File uploaded successfully"
// @Failure      400 {object} map[string]interface{} "Invalid file or file too large"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /files/upload [post]
func (s *fileService) UploadAchievementFile(c *fiber.Ctx) error {
	filename, err := middleware.UploadFile(c, "file")
	if err != nil {
		utils.GlobalLogger.Error("File upload failed", err, map[string]interface{}{
			"user_ip": c.IP(),
		})
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, "File uploaded successfully", fiber.Map{
		"filename": filename,
		"url":      "/uploads/" + filename,
	})
}

// DeleteAchievementFile godoc
// @Summary      Delete achievement attachment
// @Description  Delete an uploaded file attachment
// @Tags         Files
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        filename  path  string  true  "Filename to delete"
// @Success      200 {object} map[string]interface{} "File deleted successfully"
// @Failure      400 {object} map[string]interface{} "Failed to delete file"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /files/{filename} [delete]
func (s *fileService) DeleteAchievementFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	
	if err := middleware.DeleteFile(filename); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete file")
	}

	return utils.SuccessResponse(c, "File deleted successfully", nil)
}
