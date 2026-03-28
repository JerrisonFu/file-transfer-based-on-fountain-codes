package encoder

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"math/big"
	"os"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type EncodedPacket struct {
	Degree   uint32
	Indices  []uint32
	Symbols  []byte
	PacketID uint32
}

func (h *Handler) Encode(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	io.Copy(tmpFile, file)
	tmpFile.Seek(0, 0)

	fileInfo, _ := tmpFile.Stat()
	fileSize := fileInfo.Size()

	chunkSize := 1024
	numChunks := int(math.Ceil(float64(fileSize) / float64(chunkSize)))

	packets := generatePackets(tmpFile, numChunks, chunkSize)

	c.JSON(200, gin.H{
		"fileId":    header.Filename,
		"fileSize":  fileSize,
		"chunkSize": chunkSize,
		"numChunks": numChunks,
		"packets":   packets,
	})
}

func generatePackets(file *os.File, numChunks, chunkSize int) []EncodedPacket {
	data := make([]byte, numChunks*chunkSize)
	n, _ := file.Read(data)
	data = data[:n]

	packets := make([]EncodedPacket, 0, numChunks*2)

	for i := 0; i < numChunks*2; i++ {
		degree := int(selectDegree(numChunks))
		indices := selectIndices(degree, numChunks)

		symbols := make([]byte, chunkSize)
		for _, idx := range indices {
			start := idx * uint32(chunkSize)
			end := start + uint32(chunkSize)
			if end > uint32(len(data)) {
				end = uint32(len(data))
			}
			for j := start; j < end; j++ {
				symbols[j-start] ^= data[j]
			}
		}

		packets = append(packets, EncodedPacket{
			Degree:   uint32(degree),
			Indices:  indices,
			Symbols:  symbols[:min(chunkSize, len(symbols))],
			PacketID: uint32(i),
		})
	}

	return packets
}

func selectDegree(n int) uint32 {
	maxDegree := uint32(n)
	randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(maxDegree+1)))
	degree := uint32(randNum.Int64())
	if degree == 0 {
		degree = 1
	}
	return degree
}

func selectIndices(degree, n int) []uint32 {
	indices := make([]uint32, 0, degree)
	selected := make(map[uint32]bool)

	for len(indices) < degree {
		randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
		idx := uint32(randNum.Int64())
		if !selected[idx] {
			selected[idx] = true
			indices = append(indices, idx)
		}
	}
	return indices
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
