package main

import (
	img "github.com/LLKennedy/imagetemplate"
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
	var canvas img.Canvas
	canvas, err := img.NewCanvas(1600, 1600)
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
	canvas, err = canvas.Barcode(img.BarcodeTypeQR, []byte("www.github.com/LLKennedy/imagetemplate"), img.BarcodeExtraData{QRLevel:qr.Q, QRMode: qr.Unicode}, image.ZP, 400, 400, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeAztec, []byte("www.github.com/LLKennedy/imagetemplate"), img.BarcodeExtraData{AztecMinECCPercent: 50, AztecUserSpecifiedLayers: 4}, image.Point{X:400, Y:0}, 400, 400, color.NRGBA{G:255, A:255}, color.NRGBA{B:255, A:255})
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypePDF, []byte("Luke"), img.BarcodeExtraData{PDFSecurityLevel: 4}, image.Point{X:800, Y:0}, 400, 400, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeDataMatrix, []byte("https://www.github.com/LLKennedy/imagetemplate/demo"), img.BarcodeExtraData{}, image.Point{X:1200, Y:0}, 400, 400, color.Transparent, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeCode93, []byte("Luke"), img.BarcodeExtraData{Code93IncludeChecksum: true, Code93FullAsciiMode: true}, image.Point{X:0, Y:400}, 400, 200, nil, color.Transparent)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeType2of5, []byte("12345678"), img.BarcodeExtraData{}, image.Point{X:400, Y:400}, 400, 200, color.NRGBA{R: 255, A:255}, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeType2of5Interleaved, []byte("12345678"), img.BarcodeExtraData{}, image.Point{X:800, Y:400}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeCodabar, []byte("B123456D"), img.BarcodeExtraData{}, image.Point{X:1200, Y:400}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeCode128, []byte("Luke"), img.BarcodeExtraData{}, image.Point{X:0, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeEAN13, []byte("5901234123457"), img.BarcodeExtraData{}, image.Point{X:400, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeEAN8, []byte("11223344"), img.BarcodeExtraData{}, image.Point{X:800, Y:600}, 400, 200, nil, nil)
	if err != nil {
		log.Println(err)
	}
	canvas, err = canvas.Barcode(img.BarcodeTypeCode39, []byte("Luke"), img.BarcodeExtraData{Code39IncludeChecksum: true, Code39FullAsciiMode: true}, image.Point{X:1200, Y:600}, 400, 200, nil, nil)
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