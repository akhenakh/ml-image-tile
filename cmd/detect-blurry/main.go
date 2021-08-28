package main

import (
	"fmt"
	"log"

	"github.com/namsral/flag"
	"gocv.io/x/gocv"
)

var source = flag.String("source", "", "Source image")

const thresold = 7_000

func main() {
	flag.Parse()

	img := gocv.IMRead(*source, gocv.IMReadAnyDepth)
	if img.Empty() {
		log.Fatalf("Error reading image from: %v\n", *source)
	}
	defer img.Close()

	laplacian := gocv.NewMat()
	defer laplacian.Close()
	result := gocv.NewMat()
	defer result.Close()
	mean := gocv.NewMat()
	defer mean.Close()

	gocv.Laplacian(img, &laplacian, img.Type(), 5, 1, 0, gocv.BorderDefault)
	gocv.MeanStdDev(laplacian, &mean, &result)
	deviation := result.Mean().Val1 * result.Mean().Val1
	if deviation < thresold {
		fmt.Println("Blurry", deviation, *source)
	}
}
