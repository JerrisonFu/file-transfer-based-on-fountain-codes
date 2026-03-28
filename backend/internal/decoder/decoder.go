package decoder

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"sync"
)

type Handler struct {
	mu      sync.Mutex
	decoded map[string][]byte
}

func NewHandler() *Handler {
	return &Handler{
		decoded: make(map[string][]byte),
	}
}

type DecodeRequest struct {
	FileID    string          `json:"fileId"`
	Packets   json.RawMessage `json:"packets"`
	FileSize  int64           `json:"fileSize"`
	ChunkSize int             `json:"chunkSize"`
}

func (h *Handler) Decode(c *gin.Context) {
	var req DecodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"fileId": req.FileID,
		"status": "decoding",
	})
}
