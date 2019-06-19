package imagetemplate

import (
	"encoding/json"
	"image"
	"io"
	"io/ioutil"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/LLKennedy/imagetemplate/v3/scaffold"
	"golang.org/x/tools/godoc/vfs"
)

// Loader creates image builders from several input options and writes the finished product to several output formats.
type Loader interface {
	LoadFrom() LoadOptions
	WriteTo() WriteOptions
}

// LoadOptions chooses the input format for Loader
type LoadOptions interface {
	Builder(builder scaffold.Builder) (Loader, render.NamedProperties, error)
	Bytes(bytes []byte) (Loader, render.NamedProperties, error)
	File(path string) (Loader, render.NamedProperties, error)
	JSON(raw json.RawMessage) (Loader, render.NamedProperties, error)
	Reader(reader io.Reader) (Loader, render.NamedProperties, error)
}

// WriteOptions chooses the output format for Loader
type WriteOptions interface {
	Builder(render.NamedProperties) (scaffold.Builder, error)
	BMP(render.NamedProperties) ([]byte, error)
	BMPFile(props render.NamedProperties, path string) error
	Bytes(render.NamedProperties) ([]byte, error)
	Canvas(render.NamedProperties) (render.Canvas, error)
	Image(render.NamedProperties) (image.Image, error)
	Reader(render.NamedProperties) (io.Reader, error)
}

type loader struct {
	builder scaffold.Builder
	fs      vfs.FileSystem
}

// NewLoader returns a new loader
func NewLoader(fs vfs.FileSystem) Loader {
	if fs == nil {
		fs = vfs.OS("")
	}
	return loader{
		builder: scaffold.NewBuilder(fs),
		fs:      fs,
	}
}

// LoadFrom returns the load options for a loader
func (l loader) LoadFrom() LoadOptions {
	return nil
}

// WriteTo returns the write options for a loader
func (l loader) WriteTo() WriteOptions {
	return nil
}

// LoadTemplate takes a file path and returns a Builder constructed from the template file
func LoadTemplate(path string) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
	builder := scaffold.NewBuilder(nil)
	builder, err := builder.LoadComponentsFile(path)
	if err != nil {
		return nil, nil, err
	}
	return loadBuilder(builder)
}

// LoadReader loads JSON data from a reader, returns a list of named properties, and accepts a callback with updated properties to create BMP data
func LoadReader(reader io.Reader) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
	bytes, err := ioutil.ReadAll(reader)
	builder := scaffold.NewBuilder(nil)
	builder, err = builder.LoadComponentsData(bytes)
	if err != nil {
		return nil, nil, err
	}
	return loadBuilder(builder)
}

func loadBuilder(builder scaffold.Builder) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
	props := builder.GetNamedPropertiesList()
	cont := func(inProps render.NamedProperties) ([]byte, error) {
		builder, err := builder.SetNamedProperties(props)
		if err != nil {
			return nil, err
		}
		builder, err = builder.ApplyComponents()
		if err != nil {
			return nil, err
		}
		var data []byte
		data, err = builder.WriteToBMP()
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return props, cont, nil
}
