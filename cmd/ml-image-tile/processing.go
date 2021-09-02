package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/h2non/bimg"
)

// processImage
//  on OSX CGO_CFLAGS_ALLOW="-Xpreprocessor" go get github.com/h2non/bimg
func processImageBimg(logger log.Logger, filePath, srcDir, dstDir string, smallerTile bool, resize, width, height int) error {
	buffer, err := bimg.Read(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	img := bimg.NewImage(buffer)
	if err != nil {
		return fmt.Errorf("can't open image %s %v", filePath, err)
	}
	size, err := img.Size()
	if err != nil {
		return fmt.Errorf("can't read sizego image %s %v", filePath, err)
	}

	if resize > 1 {
		buffer, err = img.Resize(size.Width/resize, size.Height/resize)
		if err != nil {
			return fmt.Errorf("can't resize image %s %v", filePath, err)
		}
		img = bimg.NewImage(buffer)

		size, err = img.Size()
		if err != nil {
			return fmt.Errorf("can't read resized image %s %v", filePath, err)
		}

		level.Debug(logger).Log(
			"msg", "resizing image",
			"sizex", size.Width,
			"sizey", size.Height,
			"src_path", filePath,
		)
	}

	// generate tiles starting from the center
	if size.Width < width || size.Height < height {
		return fmt.Errorf("too small to be tilled %s", filePath)
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
	var modx, mody int

	// do we allow repetition in tiles ?
	if smallerTile {
		modx = size.Width % width
		mody = size.Height % height
	}

	needx := size.Width / width
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

			ext := path.Ext(filePath)
			wpath := filePath[:len(filePath)-len(ext)]
			wpath = wpath[len(srcDir):]
			outFilename := fmt.Sprintf("%s-%d%s", wpath, count, ext)
			outFilePath := filepath.Join(dstDir, outFilename)

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
				"file_path", filePath,
				"src_dir", srcDir,
				"dst_dir", dstDir,
				"out_file_path", outFilePath,
			)

			crop := bimg.NewImage(buffer)
			cropb, err := crop.Extract(ypos, xpos, width, height)
			if err != nil {
				return fmt.Errorf("can't crop image %s %v", filePath, err)
			}

			err = bimg.Write(outFilePath, cropb)
			if err != nil {
				return fmt.Errorf("can't save image %s %v", outFilePath, err)
			}
			atomic.AddUint64(&tileCounter, 1)
			count++
		}
	}

	return nil
}

// processImage
//  on OSX CGO_CFLAGS_ALLOW="-Xpreprocessor" go get github.com/h2non/bimg
func randomTileImageBimg(logger log.Logger, filePath, srcDir, dstDir string, count, resize, width, height int) error {
	buffer, err := bimg.Read(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	img := bimg.NewImage(buffer)
	if err != nil {
		return fmt.Errorf("can't open image %s %v", filePath, err)
	}
	size, err := img.Size()
	if err != nil {
		return fmt.Errorf("can't read sizego image %s %v", filePath, err)
	}

	if resize > 1 {
		buffer, err = img.Resize(size.Width/resize, size.Height/resize)
		if err != nil {
			return fmt.Errorf("can't resize image %s %v", filePath, err)
		}
		img = bimg.NewImage(buffer)

		size, err = img.Size()
		if err != nil {
			return fmt.Errorf("can't read resized image %s %v", filePath, err)
		}

		level.Debug(logger).Log(
			"msg", "resizing image",
			"sizex", size.Width,
			"sizey", size.Height,
			"src_path", filePath,
		)
	}

	// generate randome tiles
	if size.Width < width || size.Height < height {
		return fmt.Errorf("too small to be tilled %s", filePath)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < count; i++ {
		xpos := rand.Intn(size.Width - width)
		ypos := rand.Intn(size.Height - height)

		ext := path.Ext(filePath)
		wpath := filePath[:len(filePath)-len(ext)]
		wpath = wpath[len(srcDir):]
		outFilename := fmt.Sprintf("%s-%d%s", wpath, i, ext)
		outFilePath := filepath.Join(dstDir, outFilename)

		level.Debug(logger).Log(
			"msg", "cropping random image",
			"count", count,
			"xpos", xpos,
			"ypos", ypos,
			"width", width,
			"height", height,
			"sizex", size.Width,
			"sizey", size.Height,
			"file_path", filePath,
			"src_dir", srcDir,
			"dst_dir", dstDir,
			"out_file_path", outFilePath,
		)

		crop := bimg.NewImage(buffer)
		cropb, err := crop.Extract(ypos, xpos, width, height)
		if err != nil {
			return fmt.Errorf("can't crop image %s %v", filePath, err)
		}

		err = bimg.Write(outFilePath, cropb)
		if err != nil {
			return fmt.Errorf("can't save image %s %v", outFilePath, err)
		}
		atomic.AddUint64(&tileCounter, 1)
	}
	return nil
}
