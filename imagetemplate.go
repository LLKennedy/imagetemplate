// Package imagetemplate generates images from JSON templates and application-specific variables.
package imagetemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"runtime/debug"

	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/LLKennedy/imagetemplate/v3/scaffold"
	"golang.org/x/tools/godoc/vfs"
)

// Loader creates image builders from several input options and writes the finished product to several output formats.
type Loader interface {
	Load() LoadOptions
	Write() WriteOptions
}

// LoadOptions chooses the input format for Loader.
type LoadOptions interface {
	FromBuilder(builder scaffold.Builder) (Loader, render.NamedProperties, error)
	FromBytes(bytes []byte) (Loader, render.NamedProperties, error)
	FromFile(path string) (Loader, render.NamedProperties, error)
	FromJSON(raw json.RawMessage) (Loader, render.NamedProperties, error)
	FromReader(reader io.Reader) (Loader, render.NamedProperties, error)
}

// WriteOptions chooses the output format for Loader.
type WriteOptions interface {
	ToBuilder(props render.NamedProperties) (scaffold.Builder, error)
	ToBMP(props render.NamedProperties) ([]byte, error)
	ToCanvas(props render.NamedProperties) (render.Canvas, error)
	ToImage(props render.NamedProperties) (image.Image, error)
	ToBMPReader(props render.NamedProperties) (io.Reader, error)
}

type loader struct {
	builder scaffold.Builder
	fs      vfs.FileSystem
}

// New returns a new loader with the default file system.
func New() Loader {
	return NewUsing(vfs.OS("."))
}

// NewUsing returns a new loader using a specified vfs.
func NewUsing(fs vfs.FileSystem) Loader {
	if fs == nil {
		fs = vfs.OS(".")
	}
	return loader{
		fs:      fs,
		builder: scaffold.NewBuilder(fs),
	}
}

// Load returns the load options for a loader.
func (l loader) Load() LoadOptions {
	return l
}

// Write returns the write options for a loader.
func (l loader) Write() WriteOptions {
	return l
}

// FromBuilder constructs a loader using a pre-existing builder.
func (l loader) FromBuilder(builder scaffold.Builder) (newLoader Loader, props render.NamedProperties, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in FromBuilder: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder = builder
	props = l.builder.GetNamedPropertiesList()
	return l, props, nil
}

// FromBytes constructs a loader from the bytes of a template file.
func (l loader) FromBytes(bytes []byte) (newLoader Loader, props render.NamedProperties, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in FromBytes: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = l.builder.LoadComponentsData(bytes)
	return l, l.builder.GetNamedPropertiesList(), err
}

// FromFile constructs a loader from the template file located at the provided path.
func (l loader) FromFile(path string) (newLoader Loader, props render.NamedProperties, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in FromFile: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = l.builder.LoadComponentsFile(path)
	return l, l.builder.GetNamedPropertiesList(), err
}

// FromJSON constructs a loader from the raw JSON template data provided.
func (l loader) FromJSON(raw json.RawMessage) (newLoader Loader, props render.NamedProperties, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in FromJSON: %v\n%s", r, debug.Stack())
		}
	}()
	rawData, _ := raw.MarshalJSON() //This function literally cannot error, so ignore the error output
	l.builder, err = l.builder.LoadComponentsData(rawData)
	return l, l.builder.GetNamedPropertiesList(), err
}

// FromReader constructs a loader from the streamed bytes of a template file.
func (l loader) FromReader(reader io.Reader) (newLoader Loader, props render.NamedProperties, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in FromReader: %v\n%s", r, debug.Stack())
		}
	}()
	var raw []byte
	raw, err = ioutil.ReadAll(reader)
	if err != nil {
		return l, nil, err
	}
	l.builder, err = l.builder.LoadComponentsData(raw)
	return l, l.builder.GetNamedPropertiesList(), err
}

// ToBuilder returns the finished render contained within its builder.
func (l loader) ToBuilder(props render.NamedProperties) (newBuilder scaffold.Builder, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in ToBuilder: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = applyProps(l.builder, props)
	return l.builder, err
}

// ToBMP returns the finished render as the bytes of a bitmap file.
func (l loader) ToBMP(props render.NamedProperties) (raw []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in ToBMP: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = applyProps(l.builder, props)
	if err != nil {
		return nil, err
	}
	return l.builder.WriteToBMP()
}

// ToCanvas returns the finished render as a canvas object.
func (l loader) ToCanvas(props render.NamedProperties) (c render.Canvas, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in ToCanvas: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = applyProps(l.builder, props)
	if err != nil {
		return nil, err
	}
	return l.builder.GetCanvas(), nil
}

// ToImage returns the finished render as an image.Image object.
func (l loader) ToImage(props render.NamedProperties) (img image.Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in ToImage: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = applyProps(l.builder, props)
	if err != nil {
		return nil, err
	}
	return l.builder.GetCanvas().GetUnderlyingImage(), nil
}

// ToBMPReader returns the finished render as streamed bytes of a bitmap file.
func (l loader) ToBMPReader(props render.NamedProperties) (rdr io.Reader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic in ToBMPReader: %v\n%s", r, debug.Stack())
		}
	}()
	l.builder, err = applyProps(l.builder, props)
	if err != nil {
		return nil, err
	}
	rawData, err := l.builder.WriteToBMP()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(rawData), nil
}

func applyProps(builder scaffold.Builder, props render.NamedProperties) (newBuilder scaffold.Builder, err error) {
	builder, err = builder.SetNamedProperties(props)
	if err != nil {
		return builder, err
	}
	builder, err = builder.ApplyComponents()
	return builder, err
}
