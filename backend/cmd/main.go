package main

import (
	"file-transfer/backend/internal/decoder"
	"file-transfer/backend/internal/encoder"
	"file-transfer/backend/internal/transfer"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/upload", transfer.HandleUpload)
		api.GET("/download/:fileId", transfer.HandleDownload)
		api.GET("/status/:fileId", transfer.HandleStatus)
	}

	encoderHandler := encoder.NewHandler()
	decoderHandler := decoder.NewHandler()

	api.POST("/encode", encoderHandler.Encode)
	api.POST("/decode", decoderHandler.Decode)

	r.Run(":8080")
}
