
package main

import (
	"fmt"
	vidio "github.com/AlexEidt/Vidio"
	"image"
	"image/color"
	//"image/jpeg"
	"strings"
)

const (
	LENASCII            = uint16(len(ASCIILIST) - 1)
	MAXGRAYCOLOR uint16 = 65535
	//QUALITY = 100
)

var ASCIILIST = [...]rune{':',':','0','0'} 


func nearestNeighborScaling(img image.RGBA, w, h int) image.Image {
	widthimg := img.Bounds().Max.X
	heightimg := img.Bounds().Max.Y
	MinPoint := image.Point{0, 0}
	MaxPoint := image.Point{w, h}
	scalonateImg := image.NewRGBA(image.Rectangle{MinPoint, MaxPoint})
	for i := 0; i < widthimg; i++ {
		for j := 0; j < heightimg; j++ {
			srcX := float32(i) / float32(w) * float32(widthimg)
			srcY := float32(j) / float32(h) * float32(heightimg)
			srcX = min(srcX, float32(widthimg))
			srcY = min(srcY, float32(heightimg))

			r, g, b, a := img.At(int(srcX), int(srcY)).RGBA()
			newPix := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			scalonateImg.SetRGBA(i, j, newPix)
		}
	}
	return scalonateImg
}

func Gray2Accis(img image.Gray16) string {
	var asciiImg strings.Builder
	for i := 0; i < img.Bounds().Max.X; i++ {
		for j := 0; j < img.Bounds().Max.Y; j++ {
			convertValue := float32(LENASCII) * (float32(img.Gray16At(j, i).Y) / float32(MAXGRAYCOLOR))
			asciiImg.WriteRune(ASCIILIST[int(convertValue)])
			asciiImg.WriteRune(' ')
		}
		asciiImg.WriteRune('\n')
	}
	return asciiImg.String()
}

func RGB2GrayColor(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r + g + b) / 3)
}

func graySacaleImage(img image.Image) image.Gray16 {
	newImage := image.NewGray16(image.Rectangle{img.Bounds().Min, img.Bounds().Max})
	for i := 0; i < img.Bounds().Max.X; i++ {
		for j := 0; j < img.Bounds().Max.Y; j++ {
			newImage.SetGray16(i, j, color.Gray16{RGB2GrayColor(img.At(i, j))})
		}
	}
	return *newImage
}

func main(){
	video, erro := vidio.NewVideo("badApple.mp4")

	if erro != nil{
		fmt.Println(erro)
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	video.SetFrameBuffer(img.Pix)

	options := vidio.Options{
		FPS: video.FPS(),
		Bitrate: video.Bitrate(),
	}

	if video.HasStreams() {
		options.StreamFile = video.FileName()
	}

	for video.Read(){
		scalonateImage := nearestNeighborScaling(*img, 90, 90)
		grayImage := graySacaleImage(scalonateImage)
		accisImage := Gray2Accis(grayImage)
		fmt.Println(accisImage)
		fmt.Println(accisImage)
	}
}