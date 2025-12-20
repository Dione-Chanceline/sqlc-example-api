package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
)

type AttachmentHandler struct {
	queries *repo.Queries
}

func NewAttachmentHandler(q *repo.Queries) *AttachmentHandler {
	return &AttachmentHandler{queries: q}
}

func (h *AttachmentHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/attachment", h.handleUploadAttachment)
}

func (h *AttachmentHandler) handleUploadAttachment(c *gin.Context) {
	// Get message_id from form
	messageIDStr := c.PostForm("message_id")
	if messageIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message_id is required"})
		return
	}

	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message_id"})
		return
	}

	// Read the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Ensure uploads folder exists
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// Save file locally
	savePath := "./uploads/" + fileHeader.Filename

	err = c.SaveUploadedFile(fileHeader, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Insert into database with correct sqlc struct fields
	params := repo.InsertAttachmentParams{
		MessageID: messageID.String(),
		FileUrl:   savePath,
	}

	attachment, err := h.queries.InsertAttachment(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database insert failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"attachment": attachment,
	})
}
