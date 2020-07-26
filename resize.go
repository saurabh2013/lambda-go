package main

import (
	"errors"
	"image"
	"image/color"
)

// Resize image based on input scaling parameter.
func ResizeImg(img image.Image, w, h uint) (image.Image, error) {
	if img == nil {
		return nil, errors.New("Input image is found nil.")
	}
	return resize(img, img.Bounds(), 100, 100), nil
}

func resize(m image.Image, r image.Rectangle, w, h int) image.Image {

	if w == 0 || h == 0 || r.Dx() <= 0 || r.Dy() <= 0 {
		return image.NewRGBA64(image.Rect(0, 0, w, h))
	}
	curw, curh := r.Dx(), r.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			subx := x * curw / w
			suby := y * curh / h
			r32, g32, b32, a32 := m.At(subx, suby).RGBA()
			r := uint8(r32 >> 8)
			g := uint8(g32 >> 8)
			b := uint8(b32 >> 8)
			a := uint8(a32 >> 8)
			img.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}
	return img
}

// Resize images to given parameters.
// func Resize(width uint, height uint, quality uint, img []byte) (bytes []byte, err error) {
// 	mw := imagick.NewMagickWand()
// 	defer mw.Destroy()

// 	if err = mw.ReadImageBlob(img); err != nil {
// 		return nil, err
// 	}

// 	if width == 0 || height == 0 {
// 		imgWidth := mw.GetImageWidth()
// 		imgHeight := mw.GetImageHeight()

// 		if width == 0 && height == 0 {
// 			width = imgWidth
// 			height = imgHeight
// 		} else {
// 			aspectRatio := float64(imgWidth) / float64(imgHeight)
// 			if height == 0 {
// 				height = uint(float64(width) / aspectRatio)
// 			} else {
// 				width = uint(float64(height) * aspectRatio)
// 			}
// 		}
// 	}

// 	if err = mw.ThumbnailImage(width, height); err != nil {
// 		return nil, err
// 	}

// 	if err = mw.SetImageCompressionQuality(quality); err != nil {
// 		return nil, err
// 	}

// 	if bytes := mw.GetImageBlob(); len(bytes) == 0 {
// 		err = mw.GetLastError()
// 	}

// 	return bytes, err
// }
