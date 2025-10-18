package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type VideoWidget struct {
	widget.BaseWidget
	image    *canvas.Image
	width    int
	height   int
	frameNum uint64
}

func NewVideoWidget(width, height int) *VideoWidget {
	v := &VideoWidget{
		width:  width,
		height: height,
	}
	
	// Create initial placeholder image
	v.image = canvas.NewImageFromImage(v.createPlaceholder())
	v.image.FillMode = canvas.ImageFillContain
	
	v.ExtendBaseWidget(v)
	return v
}

func (v *VideoWidget) CreateRenderer() fyne.WidgetRenderer {
	return &videoRenderer{
		widget: v,
		image:  v.image,
	}
}

func (v *VideoWidget) UpdateFrame(data []byte) {
	v.frameNum++
	
	// Try to decode MJPEG frame
	img, err := v.decodeMJPEG(data)
	if err == nil {
		// Successfully decoded JPEG
		v.image.Image = img
	} else {
		// Fallback to animated pattern
		img := v.createAnimatedFrame(data)
		v.image.Image = img
	}
	
	v.image.Refresh()
}

func (v *VideoWidget) decodeMJPEG(data []byte) (image.Image, error) {
	// Decode JPEG image from bytes
	reader := bytes.NewReader(data)
	img, err := jpeg.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (v *VideoWidget) createPlaceholder() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, v.width, v.height))
	
	// Create gradient background
	for y := 0; y < v.height; y++ {
		for x := 0; x < v.width; x++ {
			gray := uint8((x + y) * 255 / (v.width + v.height))
			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}
	
	return img
}



func (v *VideoWidget) createAnimatedFrame(data []byte) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, v.width, v.height))
	
	// Create animated pattern based on frame number and data
	offset := int(v.frameNum % 100)
	dataHash := 0
	if len(data) > 0 {
		for _, b := range data[:min(len(data), 100)] {
			dataHash += int(b)
		}
	}
	
	// Animated background
	for y := 0; y < v.height; y++ {
		for x := 0; x < v.width; x++ {
			// Create moving pattern
			r := uint8((x + offset + dataHash) % 256)
			g := uint8((y + offset*2) % 256)
			b := uint8((x + y + offset) % 256)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	// Add frame indicator
	v.drawFrameInfo(img, data)
	
	return img
}



func (v *VideoWidget) drawFrameInfo(img *image.RGBA, data []byte) {
	// Draw frame number as colored bars
	barHeight := 20
	barWidth := v.width / 10
	
	for i := 0; i < 10; i++ {
		intensity := uint8((v.frameNum + uint64(i)) % 256)
		barColor := color.RGBA{intensity, 255 - intensity, 128, 255}
		
		for y := 0; y < barHeight; y++ {
			for x := i * barWidth; x < (i+1)*barWidth && x < v.width; x++ {
				img.Set(x, y, barColor)
			}
		}
	}
}

type videoRenderer struct {
	widget *VideoWidget
	image  *canvas.Image
}

func (r *videoRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
}

func (r *videoRenderer) MinSize() fyne.Size {
	return fyne.NewSize(320, 240)
}

func (r *videoRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.image}
}

func (r *videoRenderer) Refresh() {
	r.image.Refresh()
}

func (r *videoRenderer) Destroy() {}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}