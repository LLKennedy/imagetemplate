package main

import (
	img "github.com/LLKennedy/imagetemplate"
	"github.com/LLKennedy/imagetemplate/render"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"github.com/boombuler/barcode/qr"
	"os"

)

func main() {
	log.Println("Starting imagetemplate demo")
	var builder img.Builder
	var canvas render.Canvas
	canvas, err := render.NewCanvas(1600, 1600)
	if err != nil {
		log.Fatalf("Failed to create canvas: %v", err)
	}
	builder, err = img.NewBuilder(canvas, color.Gray{Y: 180})
	if err != nil {
		log.Fatalf("Failed to create builder: %v", err)
	}
	// canvas, err = canvas.Rectangle(image.Point{X: 110, Y: 40}, 60, 87, color.NRGBA{R: 255, G: 100, B: 0, A: 255})
	// if err != nil {
	// 	log.Fatalf("Failed to create rectangle: %v", err)
	// }
	// canvas, err = canvas.Circle(image.Point{X: 301, Y: 253}, 57, color.NRGBA{R: 0, G: 100, B: 255, A: 255})
	// if err != nil {
	// 	log.Fatalf("Failed to create circle: %v", err)
	// }
	canvas, err = canvas.Barcode(render.BarcodeTypeQR, []byte("www.github.com/LLKennedy/imagetemplate"), render.BarcodeExtraData{QRLevel:qr.Q, QRMode: qr.Unicode}, image.ZP, 400, 400, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeAztec, []byte("www.github.com/LLKennedy/imagetemplate"), render.BarcodeExtraData{AztecMinECCPercent: 50, AztecUserSpecifiedLayers: 4}, image.Point{X:400, Y:0}, 400, 400, color.NRGBA{G:255, A:255}, color.NRGBA{B:255, A:255})
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypePDF, []byte("Luke"), render.BarcodeExtraData{PDFSecurityLevel: 4}, image.Point{X:800, Y:0}, 400, 400, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeDataMatrix, []byte("https://www.github.com/LLKennedy/imagetemplate/demo"), render.BarcodeExtraData{}, image.Point{X:1200, Y:0}, 400, 400, color.Transparent, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeCode93, []byte("Luke"), render.BarcodeExtraData{Code93IncludeChecksum: true, Code93FullASCIIMode: true}, image.Point{X:0, Y:400}, 400, 200, nil, color.Transparent)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeType2of5, []byte("12345678"), render.BarcodeExtraData{}, image.Point{X:400, Y:400}, 400, 200, color.NRGBA{R: 255, A:255}, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeType2of5Interleaved, []byte("12345678"), render.BarcodeExtraData{}, image.Point{X:800, Y:400}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeCodabar, []byte("B123456D"), render.BarcodeExtraData{}, image.Point{X:1200, Y:400}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeCode128, []byte("Luke"), render.BarcodeExtraData{}, image.Point{X:0, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeEAN13, []byte("5901234123457"), render.BarcodeExtraData{}, image.Point{X:400, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeEAN8, []byte("11223344"), render.BarcodeExtraData{}, image.Point{X:800, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(render.BarcodeTypeCode39, []byte("Luke"), render.BarcodeExtraData{Code39IncludeChecksum: true, Code39FullASCIIMode: true}, image.Point{X:1200, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
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