package imagetemplate

import (
	"io"
	"io/ioutil"

	"github.com/LLKennedy/imagetemplate/v3/render"
)

// LoadTemplate takes a file path and returns a Builder constructed from the template file
func LoadTemplate(path string) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
	builder := NewBuilder()
	builder, err := builder.LoadComponentsFile(path)
	if err != nil {
		return nil, nil, err
	}
	return loadBuilder(builder)
}

// LoadReader loads JSON data from a reader, returns a list of named properties, and accepts a callback with updated properties to create BMP data
func LoadReader(reader io.Reader) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
	bytes, err := ioutil.ReadAll(reader)
	builder := NewBuilder()
	builder, err = builder.LoadComponentsData(bytes)
	if err != nil {
		return nil, nil, err
	}
	return loadBuilder(builder)
}

func loadBuilder(builder Builder) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
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
