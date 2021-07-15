package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"

	"github.com/h2non/bimg"
	"golang.org/x/image/draw"
)

// processImage
//  on OSX CGO_CFLAGS_ALLOW="-Xpreprocessor" go get github.com/h2non/bimg
func processImageBimg(srcPath string, outDir string, resize, x, y int) error {
	buffer, err := bimg.Read(srcPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	img := bimg.NewImage(buffer)
	if err != nil {
		return fmt.Errorf("can't open image %s %v", srcPath, err)
	}
	size, err := img.Size()
	if err != nil {
		return fmt.Errorf("can't read size image %s %v", srcPath, err)
	}

	b, err := img.Resize(size.Width/2, size.Height/2)
	if err != nil {
		return fmt.Errorf("can't resize image %s %v", srcPath, err)
	}
	img = bimg.NewImage(b)

	// generate tiles starting from the center
	if size.Width < x || size.Height < y {
		return fmt.Errorf("too small to be tilled %s", srcPath)
	}

	// find how many tiles we need
	// we want the tile a c if we have enough material (at least x/2)
	// we want the tile d e if we have enough material (at least y/2)
	// a | b | c
	// d | e | f

	count := 0

	options := bimg.Options{
		Width:  x,
		Height: y,
		Crop:   true,
		Top:    0,
		Left:   (size.Width / 2) - (x / 2),
	}

	crop, err := img.Process(options)
	if err != nil {
		return fmt.Errorf("can't crop image %s %v", srcPath, err)
	}

	ext := path.Ext(srcPath)
	wpath := srcPath[:len(srcPath)-len(ext)]
	outPath := fmt.Sprintf("%s-%d%s", wpath, count, ext)
	err = bimg.Write(, crop)
	if err != nil {
		return fmt.Errorf("can't save image %s %v", outpath, err)
	}
	return nil
}

// processImage using image/draw
func processImage(srcPath string, outPath string, resize, x, y int) error {
	input, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer input.Close()

	src, _, err := image.Decode(input)
	if err != nil {
		return err
	}

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return nil
}
