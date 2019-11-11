package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"

	"github.com/nfnt/resize"
)

func resizeImage(bigImage []byte) []byte {
	// Decode []byte into image.Image
	image, _, err := image.Decode(bytes.NewReader(bigImage))
	if err != nil {
		log.Fatal(err)
	}

	// Resize to width 250 using Lanczos resampling
	// and preserve aspect ratio
	smallImage := resize.Resize(250, 0, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	err2 := jpeg.Encode(buf, smallImage, nil)
	if err2 != nil {
		log.Fatal(err2)
	}
	return buf.Bytes()
}
