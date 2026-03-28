package encoder

import (
	"github.com/gin-gonic/gin"
	"github.com/xssnick/tonutils-go/adnl/rldp/raptorq"
	"io"
)

type Handler struct {
	rq *raptorq.RaptorQ
}

func NewHandler() *Handler {
	return &Handler{
		rq: raptorq.NewRaptorQ(1024),
	}
}

type EncodedPacket struct {
	ID   uint32 `json:"id"`
	Data []byte `json:"data"`
}

func (h *Handler) Encode(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to read file"})
		return
	}

	encoder, err := h.rq.CreateEncoder(data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	symbolSize := h.rq.GetSymbolSize()
	baseSymbols := encoder.BaseSymbolsNum()
	totalSymbols := uint32(len(data))/symbolSize + 1

	packets := make([]EncodedPacket, 0, totalSymbols*2)
	for i := uint32(0); i < totalSymbols*2; i++ {
		symbol := encoder.GenSymbol(i)
		packets = append(packets, EncodedPacket{
			ID:   i,
			Data: symbol,
		})
	}

	c.JSON(200, gin.H{
		"fileId":       header.Filename,
		"fileSize":     len(data),
		"symbolSize":   symbolSize,
		"baseSymbols":  baseSymbols,
		"totalSymbols": totalSymbols,
		"packets":      packets,
	})
}
