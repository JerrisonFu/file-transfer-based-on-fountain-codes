package transfer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"os"
)

func HandleUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	fileID := uuid.New().String()
	tmpPath := fmt.Sprintf("data/%s_%s", fileID, header.Filename)

	os.MkdirAll("data", 0755)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer tmpFile.Close()

	io.Copy(tmpFile, file)

	c.JSON(200, gin.H{
		"fileId":   fileID,
		"filename": header.Filename,
		"size":     header.Size,
	})
}

func HandleDownload(c *gin.Context) {
	fileID := c.Param("fileId")
	filePath := fmt.Sprintf("data/%s", fileID)

	c.File(filePath)
}

func HandleStatus(c *gin.Context) {
	fileID := c.Param("fileId")

	_, err := os.Stat(fmt.Sprintf("data/%s", fileID))
	if err != nil {
		c.JSON(404, gin.H{"error": "file not found"})
		return
	}

	c.JSON(200, gin.H{
		"fileId": fileID,
		"status": "ready",
	})
}
