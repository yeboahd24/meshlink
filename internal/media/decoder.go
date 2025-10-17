package media

import (
	"encoding/json"
	"fmt"
)

type H264Decoder struct {
	isDecoding bool
}

func NewH264Decoder() *H264Decoder {
	return &H264Decoder{}
}

func (d *H264Decoder) Start() error {
	if d.isDecoding {
		return fmt.Errorf("decoder already started")
	}
	
	d.isDecoding = true
	return nil
}

func (d *H264Decoder) Stop() {
	d.isDecoding = false
}

func (d *H264Decoder) DecodeFrame(encodedData []byte) (*DecodedFrame, error) {
	if !d.isDecoding {
		return nil, fmt.Errorf("decoder not started")
	}
	
	if len(encodedData) < 4 {
		return nil, fmt.Errorf("invalid frame data: too short")
	}
	
	// Read metadata length
	metadataLen := int(encodedData[0])<<24 | int(encodedData[1])<<16 | int(encodedData[2])<<8 | int(encodedData[3])
	
	if len(encodedData) < 4+metadataLen {
		return nil, fmt.Errorf("invalid frame data: metadata length mismatch")
	}
	
	// Extract metadata
	metadataBytes := encodedData[4 : 4+metadataLen]
	var framePackage struct {
		Metadata FrameMetadata `json:"metadata"`
		DataSize int           `json:"data_size"`
	}
	
	if err := json.Unmarshal(metadataBytes, &framePackage); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}
	
	// Extract frame data
	frameData := encodedData[4+metadataLen:]
	
	if len(frameData) != framePackage.DataSize {
		return nil, fmt.Errorf("frame data size mismatch: expected %d, got %d", 
			framePackage.DataSize, len(frameData))
	}
	
	return &DecodedFrame{
		Metadata: framePackage.Metadata,
		Data:     frameData,
	}, nil
}

type DecodedFrame struct {
	Metadata FrameMetadata
	Data     []byte
}

func (f *DecodedFrame) GetFrameID() uint64 {
	return f.Metadata.FrameID
}

func (f *DecodedFrame) GetQuality() string {
	return f.Metadata.Quality
}

func (f *DecodedFrame) GetSize() int {
	return f.Metadata.Size
}

func (f *DecodedFrame) GetBitrate() int {
	return f.Metadata.Bitrate
}