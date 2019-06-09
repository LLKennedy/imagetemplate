package imagetemplate

import (
	"io/ioutil"
	"os"

	"github.com/LLKennedy/imagetemplate/v2/render"
)

// LoadTemplate takes a file path and returns a Builder constructed from the template file
func LoadTemplate(path string) (render.NamedProperties, func(render.NamedProperties) ([]byte, error), error) {
	builder := NewBuilder()
	builder, err := builder.LoadComponentsFile(path)
	if err != nil {
		return nil, nil, err
	}
	props := builder.GetNamedPropertiesList()
	cont := func(inProps render.NamedProperties) ([]byte, error) {
		builder, err = builder.SetNamedProperties(props)
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
		err = ioutil.WriteFile("simple-static.bmp", data, os.ModeExclusive)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return props, cont, nil
}
