package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/h2non/bimg"
	"golang.org/x/image/draw"
)

type Direction uint8

const (
	Center Direction = 1 << iota
	North
	East
	South
	West
)

func (d Direction) HasDirection(dir Direction) bool { return d&dir != 0 }
func (d *Direction) AddDirection(dir Direction)     { *d |= dir }

// processImage
//  on OSX CGO_CFLAGS_ALLOW="-Xpreprocessor" go get github.com/h2non/bimg
func processImageBimg(logger log.Logger, srcPath string, outDir string, resize, width, height int) error {
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

	if resize > 1 {
		buffer, err = img.Resize(size.Width/resize, size.Height/resize)
		if err != nil {
			return fmt.Errorf("can't resize image %s %v", srcPath, err)
		}
		img = bimg.NewImage(buffer)

		size, err = img.Size()
		if err != nil {
			return fmt.Errorf("can't read resized image %s %v", srcPath, err)
		}

		level.Debug(logger).Log(
			"msg", "resizing image",
			"sizex", size.Width,
			"sizey", size.Height,
			"src_path", srcPath,
		)
	}

	// generate tiles starting from the center
	if size.Width < width || size.Height < height {
		return fmt.Errorf("too small to be tilled %s", srcPath)
	}

	count := 0

	// start at the top left
	var xpos, ypos int

	// find how many tiles we need
	// we want the tile a c if we have enough material (at least x/2)
	// we want the tile d e if we have enough material (at least y/2)
	// a | b | c
	// d | e | f

	// a line is the number of slice + the extra half overlap
	modx := size.Width % width
	needx := size.Width / width

	mody := size.Height % height
	needy := size.Height / height

	if modx > width/2 {
		needx += 2
	}

	if mody > height/2 {
		needy += 2
	}

	// descending loop
	for cuty := 0; cuty < needy; cuty++ {
		if mody < height/2 {
			ypos = cuty*height + mody/2
		} else {
			ypos = (cuty-1)*height + mody/2
			if cuty == 0 {
				ypos = 0
			}
			if cuty == needy-1 {
				ypos = size.Height - height
			}
		}

		// save an horizontal slice
		for cutx := 0; cutx < needx; cutx++ {
			if modx < width/2 {
				xpos = cutx*width + modx/2
			} else {
				xpos = (cutx-1)*width + modx/2
				if cutx == 0 {
					xpos = 0
				}
				if cutx == needx-1 {
					xpos = size.Width - width
				}
			}

			level.Debug(logger).Log(
				"msg", "cropping image",
				"count", count,
				"modx", modx,
				"mody", mody,
				"needx", needx,
				"needy", needy,
				"xpos", xpos,
				"ypos", ypos,
				"cutx", cutx,
				"cuty", cuty,
				"width", width,
				"height", height,
				"sizex", size.Width,
				"sizey", size.Height,
				"src_path", srcPath,
			)
			err := saveTile(buffer, srcPath, outDir, ypos, xpos, width, height, count)
			if err != nil {
				return err
			}

			count++
		}
	}

	return nil
}

func saveTile(
	buffer []byte,
	srcPath string,
	outDir string,
	ypos, xpos, width, height, count int,
) error {
	crop := bimg.NewImage(buffer)
	cropb, err := crop.Extract(ypos, xpos, width, height)
	if err != nil {
		return fmt.Errorf("can't crop image %s %v", srcPath, err)
	}

	ext := path.Ext(srcPath)
	wpath := srcPath[:len(srcPath)-len(ext)]
	outPath := fmt.Sprintf("%s/%s-%d%s", outDir, wpath, count, ext)
	err = bimg.Write(outPath, cropb)
	if err != nil {
		return fmt.Errorf("can't save image %s %v", outPath, err)
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
