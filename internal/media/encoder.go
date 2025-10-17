package media

import (
	"encoding/json"
	"fmt"
	"time"
)

type H264Encoder struct {
	bitrate    int
	quality    string
	profile    string
	isEncoding bool
}

func NewH264Encoder(quality string) *H264Encoder {
	encoder := &H264Encoder{
		quality: quality,
		profile: "baseline",
	}
	
	// Set bitrate based on quality
	switch quality {
	case "1080p":
		encoder.bitrate = 4000000 // 4 Mbps
	case "720p":
		encoder.bitrate = 2000000 // 2 Mbps
	case "480p":
		encoder.bitrate = 1000000 // 1 Mbps
	default:
		encoder.bitrate = 2000000
	}
	
	return encoder
}

func (e *H264Encoder) Start() error {
	if e.isEncoding {
		return fmt.Errorf("encoder already started")
	}
	
	e.isEncoding = true
	return nil
}

func (e *H264Encoder) Stop() {
	e.isEncoding = false
}

func (e *H264Encoder) EncodeFrame(rawData []byte, frameID uint64) ([]byte, error) {
	if !e.isEncoding {
		return nil, fmt.Errorf("encoder not started")
	}
	
	// Create frame metadata
	frameInfo := FrameMetadata{
		FrameID:   frameID,
		Timestamp: time.Now(),
		Type:      "video",
		Codec:     "h264",
		Quality:   e.quality,
		Bitrate:   e.bitrate,
		Profile:   e.profile,
		Size:      len(rawData),
	}
	
	// Encode frame info + data
	return e.encodeWithMetadata(frameInfo, rawData)
}

type FrameMetadata struct {
	FrameID   uint64    `json:"frame_id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Codec     string    `json:"codec"`
	Quality   string    `json:"quality"`
	Bitrate   int       `json:"bitrate"`
	Profile   string    `json:"profile"`
	Size      int       `json:"size"`
}

func (e *H264Encoder) encodeWithMetadata(metadata FrameMetadata, data []byte) ([]byte, error) {
	// Create frame package with metadata + data
	framePackage := struct {
		Metadata FrameMetadata `json:"metadata"`
		DataSize int           `json:"data_size"`
	}{
		Metadata: metadata,
		DataSize: len(data),
	}
	
	// Serialize metadata
	metadataBytes, err := json.Marshal(framePackage)
	if err != nil {
		return nil, fmt.Errorf("failed to encode metadata: %w", err)
	}
	
	// Create final frame: [metadata_length][metadata][data]
	metadataLen := len(metadataBytes)
	totalSize := 4 + metadataLen + len(data) // 4 bytes for length + metadata + data
	
	result := make([]byte, totalSize)
	
	// Write metadata length (4 bytes)
	result[0] = byte(metadataLen >> 24)
	result[1] = byte(metadataLen >> 16)
	result[2] = byte(metadataLen >> 8)
	result[3] = byte(metadataLen)
	
	// Write metadata
	copy(result[4:4+metadataLen], metadataBytes)
	
	// Write frame data
	copy(result[4+metadataLen:], data)
	
	return result, nil
}