package decoder

import (
	"github.com/gin-gonic/gin"
	"github.com/xssnick/tonutils-go/adnl/rldp/raptorq"
	"sync"
)

type Handler struct {
	mu       sync.Mutex
	rq       *raptorq.RaptorQ
	decoders map[string]*decoderState
}

type decoderState struct {
	decoder   *raptorq.Decoder
	received  int
	completed bool
}

func NewHandler() *Handler {
	return &Handler{
		rq:       raptorq.NewRaptorQ(1024),
		decoders: make(map[string]*decoderState),
	}
}

type DecodeRequest struct {
	FileID   string   `json:"fileId"`
	Packets  []Packet `json:"packets"`
	FileSize int      `json:"fileSize"`
}

type Packet struct {
	ID   uint32 `json:"id"`
	Data []byte `json:"data"`
}

func (h *Handler) Decode(c *gin.Context) {
	var req DecodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.mu.Lock()
	state, exists := h.decoders[req.FileID]
	if !exists {
		decoder, err := h.rq.CreateDecoder(uint32(req.FileSize))
		if err != nil {
			h.mu.Unlock()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		state = &decoderState{
			decoder:  decoder,
			received: 0,
		}
		h.decoders[req.FileID] = state
	}
	h.mu.Unlock()

	for _, pkt := range req.Packets {
		_, err := state.decoder.AddSymbol(pkt.ID, pkt.Data)
		if err != nil {
			continue
		}
		state.received++
	}

	complete, data, err := state.decoder.Decode()
	if err != nil {
		c.JSON(200, gin.H{
			"fileId":   req.FileID,
			"status":   "decoding",
			"received": state.received,
		})
		return
	}

	if complete {
		h.mu.Lock()
		state.completed = true
		h.mu.Unlock()

		c.JSON(200, gin.H{
			"fileId":   req.FileID,
			"status":   "complete",
			"data":     data,
			"received": state.received,
		})
	} else {
		c.JSON(200, gin.H{
			"fileId":   req.FileID,
			"status":   "decoding",
			"received": state.received,
		})
	}
}
