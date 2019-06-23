package datetime

import (
	"runtime/debug"
	"testing"

	"github.com/LLKennedy/gosysfonts"
	"github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs"
)

func TestDateTimeSetNamedPropertiesOS(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input render.NamedProperties
		res   Component
		err   string
	}
	tests := []testSet{
		{
			name: "load font file from real fs",
			start: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
			},
			input: render.NamedProperties{
				"aProp": "gibberish file that doesn't exist",
			},
			res: Component{
				NamedPropertiesMap: map[string][]string{
					"aProp": {"fontFile"},
				},
				fs: vfs.OS("."),
			},
			err: "open gibberish file that doesn't exist: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					assert.Failf(t, "caught panic", "%v\n%s", r, debug.Stack())
				}
			}()
			res, err := test.start.SetNamedProperties(test.input)
			assert.Equal(t, test.res, res)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err)
			}
			if mockFs, isMock := test.start.fs.(*filesystem.MockFileSystem); isMock {
				mockFs.AssertExpectations(t)
			}
		})
	}
}

func TestDateTimeVerifyAndTestTextJSONDataOS(t *testing.T) {
	type testSet struct {
		name  string
		start Component
		input interface{}
		res   Component
		props render.NamedProperties
		err   string
	}
	tests := []testSet{
		{
			name: "gibberish font file",
			start: Component{
				fontPool: fakeSysFonts{},
			},
			input: &datetimeFormat{
				Font: struct {
					FontName string `json:"fontName"`
					FontFile string `json:"fontFile"`
					FontURL  string `json:"fontURL"`
				}{
					FontFile: "gibberish file that doesn't exist",
				},
			},
			res: Component{
				fontPool: fakeSysFonts{},
			},
			props: render.NamedProperties{},
			err: "open gibberish file that doesn't exist: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, props, err := test.start.VerifyAndSetJSONData(test.input)
			assert.Equal(t, test.res, res)
			assert.Equal(t, test.props, props)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.err)
			}
		})
	}
}

func TestInit(t *testing.T) {
	c, err := render.Decode("datetime")
	assert.NoError(t, err)
	assert.Equal(t, Component{fs: vfs.OS("."), fontPool: gosysfonts.OSXPool{}}, c)
}
