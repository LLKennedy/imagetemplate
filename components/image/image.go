// Package image is an embedded image component with support for jpg, png, bmp and tiff files.
package image

import (
	"fmt"
	"image"
	_ "image/jpeg" // jpeg imported for image decoding
	_ "image/png"  // png imported for image decoding
	"io/ioutil"
	"os"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/bmp"  // bmp imported for image decoding
	_ "golang.org/x/image/tiff" // tiff imported for image decoding
	"golang.org/x/tools/godoc/vfs"
)

// Component implements the Component interface for images.
type Component struct {
	/*
		NamedPropertiesMap maps user/application variables to properties of the component.
		This field is filled automatically by VerifyAndSetJSONData, then used in
		SetNamedProperties to determine whether a variable being passed in is relevant to this
		component.

		For example, map[string][]string{"photo": []string{"fileName"}} would indicate that
		the user specified variable "photo" will fill the Image property via an image file.
	*/
	NamedPropertiesMap map[string][]string
	// Image is the image to draw on the canvas.
	Image image.Image
	/*
		TopLeft is the coordinates of the top-left corner of the image relative to the
		top-left corner of the canvas.
	*/
	TopLeft image.Point
	// Width is the width to scale the image to.
	Width int
	// Height is the height to scale the image to.
	Height int
	// fs is the file system.
	fs vfs.FileSystem
}

type imageFormat struct {
	TopLeftX string `json:"topLeftX"`
	TopLeftY string `json:"topLeftY"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	FileName string `json:"fileName"`
	Data     string `json:"data"`
}

// Write draws an image on the canvas.
func (component Component) Write(canvas render.Canvas) (render.Canvas, error) {
	if len(component.NamedPropertiesMap) != 0 {
		return canvas, fmt.Errorf("cannot draw image, not all named properties are set: %v", component.NamedPropertiesMap)
	}
	c := canvas
	var err error
	scaledImage := imaging.Resize(component.Image, component.Width, component.Height, imaging.Lanczos)
	c, err = c.DrawImage(component.TopLeft, scaledImage)
	if err != nil {
		return canvas, err
	}
	return c, nil
}

// SetNamedProperties processes the named properties and sets them into the image properties.
func (component Component) SetNamedProperties(properties render.NamedProperties) (render.Component, error) {
	c := component
	var err error
	c.NamedPropertiesMap, err = render.StandardSetNamedProperties(properties, component.NamedPropertiesMap, (&c).delegatedSetProperties)
	if err != nil {
		return component, err
	}
	return c, nil
}

// GetJSONFormat returns the JSON structure of a image component.
func (component Component) GetJSONFormat() interface{} {
	return &imageFormat{}
}

// VerifyAndSetJSONData processes the data parsed from JSON and uses it to set image properties and fill the named properties map.
func (component Component) VerifyAndSetJSONData(data interface{}) (render.Component, render.NamedProperties, error) {
	c := component
	props := make(render.NamedProperties)
	stringStruct, ok := data.(*imageFormat)
	if !ok {
		return component, props, fmt.Errorf("failed to convert returned data to component properties")
	}
	return c.parseJSONFormat(stringStruct, props)
}

type osFileSystem struct {
}

func (ofs *osFileSystem) Open(path string) (vfs.ReadSeekCloser, error) {
	return os.Open(path)
}

func (ofs *osFileSystem) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

func (ofs *osFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (ofs *osFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func (ofs *osFileSystem) RootType(path string) vfs.RootType {
	return vfs.OS(".").RootType(path)
}

func (ofs *osFileSystem) String() string {
	return vfs.OS(".").String()
}

func (component Component) getFileSystem() vfs.FileSystem {
	return &osFileSystem{}
}

func init() {
	for _, name := range []string{"image", "img", "photo", "Image", "IMG", "Photo", "picture", "Picture", "IMAGE", "PHOTO", "PICTURE"} {
		render.RegisterComponent(name, func(fs vfs.FileSystem) render.Component { return Component{fs: fs} })
	}
}
