package scaffold

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"

	"github.com/LLKennedy/imagetemplate/v3/cutils"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/disintegration/imaging"
)

func (builder ImageBuilder) setBackgroundImage(template Template) (b ImageBuilder, err error) {
	b = builder
	// Check the state of the optional and required properties
	dataSet := template.BaseImage.Data != ""
	fileSet := template.BaseImage.FileName != ""
	baseColourSet := template.BaseImage.BaseWidth != "" && template.BaseImage.BaseHeight != "" && (template.BaseImage.BaseColour.Red != "" || template.BaseImage.BaseColour.Green != "" || template.BaseImage.BaseColour.Blue != "" || template.BaseImage.BaseColour.Alpha != "")
	oneSet := dataSet || fileSet || baseColourSet
	if !oneSet {
		return builder.SetCanvas(builder.GetCanvas()).(ImageBuilder), nil
	}
	if cutils.ExclusiveNor(dataSet, fileSet, baseColourSet) {
		return builder, fmt.Errorf("cannot load base image from file and load from data string and generate from base colour, specify only data or fileName or base colour")
	}
	switch {
	case dataSet:
		b, err = b.setBaseData(template)
	case fileSet:
		b, err = b.setBaseFile(template)
	case baseColourSet:
		b, err = b.setBaseColour(template)
	}
	return
}

func (builder ImageBuilder) setBaseData(template Template) (ImageBuilder, error) {
	b := builder
	// Get image data from string
	sReader := strings.NewReader(template.BaseImage.Data)
	imageData := base64.NewDecoder(base64.StdEncoding, sReader)
	// Decode image data
	baseImage, _, err := image.Decode(imageData)
	if err != nil {
		return builder, err
	}
	return b.baseConvertAndResize(baseImage, template)
}

func (builder ImageBuilder) baseConvertAndResize(baseImage image.Image, template Template) (b ImageBuilder, err error) {
	b = builder
	if ycbcr, ok := baseImage.(*image.YCbCr); ok {
		var newImage draw.Image
		newImage = image.NewNRGBA(ycbcr.Rect)
		draw.Draw(newImage, ycbcr.Rect, ycbcr, ycbcr.Bounds().Min, draw.Over)
		baseImage = newImage
	}
	if b.Canvas == nil {
		// No current canvas, uses loaded image as canvas
		var drawImage draw.Image
		drawImage = image.NewNRGBA(baseImage.Bounds())
		draw.Draw(drawImage, baseImage.Bounds(), baseImage, baseImage.Bounds().Min, draw.Over)
		b = b.SetCanvas(render.ImageCanvas{Image: drawImage}).(ImageBuilder)
		ppi, err := strconv.ParseFloat(template.BaseImage.PPI, 64)
		if err != nil || ppi == 0 {
			ppi = float64(72)
		}
		canvas := b.GetCanvas().SetPPI(ppi)
		b = (b.SetCanvas(canvas)).(ImageBuilder)
		return b, nil
	}
	// Check if resizing is necessary
	currentHeight, currentWidth := baseImage.Bounds().Size().Y, baseImage.Bounds().Size().X
	targetHeight, targetWidth := b.GetCanvas().GetHeight(), b.GetCanvas().GetWidth()
	if targetHeight != currentHeight || targetWidth != currentWidth {
		// Compare aspect ratios
		targetAspect := float64(targetWidth) / float64(targetHeight)
		currentAspect := float64(currentWidth) / float64(currentHeight)
		var resizedWidth, resizedHeight int
		if targetAspect == currentAspect {
			// Identical apsect ratios
			resizedWidth = targetWidth
			resizedHeight = targetHeight
		} else if targetAspect < currentAspect {
			// Fit wide image into thin frame
			resizedHeight = targetHeight
		} else {
			// Fit thin image into wide frame
			resizedWidth = targetWidth
		}
		baseImage = imaging.Resize(baseImage, resizedWidth, resizedHeight, imaging.Lanczos)
	}
	ppi, err := strconv.ParseFloat(template.BaseImage.PPI, 64)
	if err != nil || ppi == 0 {
		ppi = float64(72)
	}
	canvas := b.GetCanvas().SetPPI(ppi)
	canvas, err = canvas.DrawImage(image.Point{X: 0, Y: 0}, baseImage)
	if err != nil {
		return builder, err
	}
	b = (b.SetCanvas(canvas)).(ImageBuilder)
	return b, nil
}

func (builder ImageBuilder) setBaseFile(template Template) (ImageBuilder, error) {
	b := builder
	// Get image data from file

	imgFile, err := b.fs.Open(template.BaseImage.FileName)
	if err != nil {
		return builder, err
	}
	defer imgFile.Close()
	imageData := imgFile
	// Decode image data
	baseImage, _, err := image.Decode(imageData)
	if err != nil {
		return builder, err
	}
	return b.baseConvertAndResize(baseImage, template)
}

func (builder ImageBuilder) setBaseColour(template Template) (b ImageBuilder, err error) {
	b = builder
	width64, err := strconv.ParseInt(template.BaseImage.BaseWidth, 10, 64) //Use ParseInt instead of Atoi for compatibility with go 1.7
	if err != nil {
		return builder, err
	}
	width := int(width64)
	height64, err := strconv.ParseInt(template.BaseImage.BaseHeight, 10, 64) //Use ParseInt instead of Atoi for compatibility with go 1.7
	if err != nil {
		return builder, err
	}
	height := int(height64)
	red64, err := strconv.ParseUint(template.BaseImage.BaseColour.Red, 0, 8)
	if err != nil {
		return builder, err
	}
	red := uint8(red64)
	green64, err := strconv.ParseUint(template.BaseImage.BaseColour.Green, 0, 8)
	if err != nil {
		return builder, err
	}
	green := uint8(green64)
	blue64, err := strconv.ParseUint(template.BaseImage.BaseColour.Blue, 0, 8)
	if err != nil {
		return builder, err
	}
	blue := uint8(blue64)
	alpha64, err := strconv.ParseUint(template.BaseImage.BaseColour.Alpha, 0, 8)
	if err != nil {
		return builder, err
	}
	alpha := uint8(alpha64)
	rectangle := image.Rect(0, 0, width, height)
	baseImage := image.NewNRGBA(rectangle)
	colourPlane := image.NewUniform(color.NRGBA{R: red, G: green, B: blue, A: alpha})
	draw.Draw(baseImage, rectangle, colourPlane, image.Point{X: 0, Y: 0}, draw.Over)
	return b.baseConvertAndResize(baseImage, template)
}
