package main

import (
	img "github.com/LLKennedy/imagetemplate"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.Println("Starting imagetemplate demo")
	var builder img.Builder
	var canvas img.Canvas
	canvas, err := img.NewCanvas(100, 100)
	if err != nil {
		log.Fatalf("Failed to create canvas: %v", err)
	}
	builder, err = img.NewBuilder(canvas, color.Gray{Y: 128})
	if err != nil {
		log.Fatalf("Failed to create builder: %v", err)
	}
	err = canvas.Rectangle(image.Point{X: 10, Y: 40}, 15, 20, color.NRGBA{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		log.Fatalf("Failed to create rectangle: %v", err)
	}
	err = canvas.Circle(image.Point{X: 60, Y: 35}, 15, color.NRGBA{R: 0, G: 100, B: 255, A: 255})
	if err != nil {
		log.Fatalf("Failed to create circle: %v", err)
	}
	bytes, err := builder.WriteToBMP()
	if err != nil {
		log.Fatalf("Failed to write canvas to bitmap: %v", err)
	}
	err = ioutil.WriteFile("demo.bmp", bytes, os.ModeExclusive)
	if err != nil {
		log.Fatalf("Failed to write bitmap to file: %v", err)
	}
}
