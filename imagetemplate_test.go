package imagetemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"testing"

	fs "github.com/LLKennedy/imagetemplate/v3/internal/filesystem"
	"github.com/LLKennedy/imagetemplate/v3/render"
	"github.com/LLKennedy/imagetemplate/v3/scaffold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBuilder struct {
	mock.Mock
}

func (b *mockBuilder) GetCanvas() render.Canvas {
	args := b.Called()
	return args.Get(0).(render.Canvas)
}
func (b *mockBuilder) SetCanvas(newCanvas render.Canvas) scaffold.Builder {
	args := b.Called(newCanvas)
	return args.Get(0).(scaffold.Builder)
}
func (b *mockBuilder) GetComponents() []render.Component {
	args := b.Called()
	return args.Get(0).([]render.Component)
}
func (b *mockBuilder) SetComponents(components []scaffold.ToggleableComponent) scaffold.Builder {
	args := b.Called(components)
	return args.Get(0).(scaffold.Builder)
}
func (b *mockBuilder) GetNamedPropertiesList() render.NamedProperties {
	args := b.Called()
	return args.Get(0).(render.NamedProperties)
}
func (b *mockBuilder) SetNamedProperties(properties render.NamedProperties) (scaffold.Builder, error) {
	args := b.Called(properties)
	return args.Get(0).(scaffold.Builder), args.Error(1)
}
func (b *mockBuilder) ApplyComponents() (scaffold.Builder, error) {
	args := b.Called()
	return args.Get(0).(scaffold.Builder), args.Error(1)
}
func (b *mockBuilder) LoadComponentsFile(fileName string) (scaffold.Builder, error) {
	args := b.Called(fileName)
	return args.Get(0).(scaffold.Builder), args.Error(1)
}
func (b *mockBuilder) LoadComponentsData(fileData []byte) (scaffold.Builder, error) {
	args := b.Called(fileData)
	return args.Get(0).(scaffold.Builder), args.Error(1)
}
func (b *mockBuilder) WriteToBMP() ([]byte, error) {
	args := b.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func TestLoadWrite(t *testing.T) {
	l := New()
	assert.Equal(t, l.Load(), l)
	assert.Equal(t, l.Write(), l)
}

func TestLoadMethods(t *testing.T) {
	b := new(mockBuilder)
	nilProps := render.NamedProperties(nil)
	b.On("GetNamedPropertiesList").Return(nilProps)
	mfs := fs.NewMockFileSystem()
	l := loader{
		builder: b,
		fs:      mfs,
	}
	t.Run("FromBuilder", func(t *testing.T) {
		l2, props, err := l.FromBuilder(b)
		assert.Equal(t, l, l2)
		assert.Equal(t, nilProps, props)
		assert.NoError(t, err)
	})
	t.Run("FromBytes", func(t *testing.T) {
		b.On("LoadComponentsData", []byte("hello")).Return(b, fmt.Errorf("some error"))
		l2, props, err := l.FromBytes([]byte("hello"))
		assert.Equal(t, l, l2)
		assert.Equal(t, nilProps, props)
		assert.EqualError(t, err, "some error")
	})
	t.Run("FromFile", func(t *testing.T) {
		b.On("LoadComponentsFile", "testfile").Return(b, fmt.Errorf("file load error"))
		l2, props, err := l.FromFile("testfile")
		assert.Equal(t, l, l2)
		assert.Equal(t, nilProps, props)
		assert.EqualError(t, err, "file load error")
	})
	t.Run("FromJSON", func(t *testing.T) {
		jsonBytes := []byte(`
		{
			"testKey": "testVal"
		}
		`)
		type rawStuff struct {
			TestKey json.RawMessage `json:"testKey"`
		}
		newRaw := &rawStuff{}
		err := json.Unmarshal(jsonBytes, newRaw)
		assert.NoError(t, err)
		b.On("LoadComponentsData", []byte(`"testVal"`)).Return(b, fmt.Errorf("json error"))
		l2, props, err := l.FromJSON(newRaw.TestKey)
		assert.Equal(t, l, l2)
		assert.Equal(t, nilProps, props)
		assert.EqualError(t, err, "json error")
	})
	t.Run("FromReader", func(t *testing.T) {
		t.Run("invalid reader", func(t *testing.T) {
			reader := badReader{}
			l2, props, err := l.FromReader(reader)
			assert.Equal(t, l, l2)
			assert.Nil(t, props)
			assert.EqualError(t, err, "not a real reader")
		})
		t.Run("valid reader", func(t *testing.T) {
			reader := bytes.NewReader([]byte("some data"))
			b.On("LoadComponentsData", []byte("some data")).Return(b, fmt.Errorf("reader error"))
			l2, props, err := l.FromReader(reader)
			assert.Equal(t, l, l2)
			assert.Equal(t, nilProps, props)
			assert.EqualError(t, err, "reader error")
		})
	})
	b.AssertExpectations(t)
}

type badReader struct {
}

func (r badReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("not a real reader")
}

func TestWriteMethods(t *testing.T) {
	b := new(mockBuilder)
	badProps := render.NamedProperties{"a": 1}
	nilProps := render.NamedProperties(nil)
	b.On("SetNamedProperties", nilProps).Return(b, nil)
	b.On("ApplyComponents").Return(b, nil)
	b.On("SetNamedProperties", badProps).Return(b, fmt.Errorf("bad props"))
	mockCanvas := new(render.MockCanvas)
	baseImage := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	mockCanvas.On("GetUnderlyingImage").Return(baseImage)
	b.On("GetCanvas").Return(mockCanvas)
	b.On("WriteToBMP").Return([]byte("a bmp"), fmt.Errorf("bmp error"))
	mfs := fs.NewMockFileSystem()
	l := loader{
		builder: b,
		fs:      mfs,
	}
	t.Run("ToBuilder", func(t *testing.T) {
		t.Run("invalid props", func(t *testing.T) {
			res, err := l.Write().ToBuilder(badProps)
			assert.Equal(t, b, res)
			assert.EqualError(t, err, "bad props")
		})
		t.Run("valid props", func(t *testing.T) {
			res, err := l.Write().ToBuilder(nilProps)
			assert.Equal(t, b, res)
			assert.NoError(t, err)
		})
	})
	t.Run("ToBMP", func(t *testing.T) {
		t.Run("invalid props", func(t *testing.T) {
			res, err := l.Write().ToBMP(badProps)
			assert.Nil(t, res)
			assert.EqualError(t, err, "bad props")
		})
		t.Run("valid props", func(t *testing.T) {
			res, err := l.Write().ToBMP(nilProps)
			assert.Equal(t, []byte("a bmp"), res)
			assert.EqualError(t, err, "bmp error")
		})
	})
	t.Run("ToCanvas", func(t *testing.T) {
		t.Run("invalid props", func(t *testing.T) {
			res, err := l.Write().ToCanvas(badProps)
			assert.Nil(t, res)
			assert.EqualError(t, err, "bad props")
		})
		t.Run("valid props", func(t *testing.T) {
			res, err := l.Write().ToCanvas(nilProps)
			assert.Equal(t, mockCanvas, res)
			assert.NoError(t, err)
		})
	})
	t.Run("ToImage", func(t *testing.T) {
		t.Run("invalid props", func(t *testing.T) {
			res, err := l.Write().ToImage(badProps)
			assert.Nil(t, res)
			assert.EqualError(t, err, "bad props")
		})
		t.Run("valid props", func(t *testing.T) {
			res, err := l.Write().ToImage(nilProps)
			assert.Equal(t, baseImage, res)
			assert.NoError(t, err)
		})
	})
	t.Run("ToBMPReader", func(t *testing.T) {
		t.Run("invalid props", func(t *testing.T) {
			res, err := l.Write().ToBMPReader(badProps)
			assert.Nil(t, res)
			assert.EqualError(t, err, "bad props")
		})
		t.Run("valid props", func(t *testing.T) {
			t.Run("error writing to BMP", func(t *testing.T) {
				res, err := l.Write().ToBMPReader(nilProps)
				assert.Nil(t, res)
				assert.EqualError(t, err, "bmp error")
			})
			t.Run("success", func(t *testing.T) {
				b := new(mockBuilder)
				nilProps := render.NamedProperties(nil)
				b.On("SetNamedProperties", nilProps).Return(b, nil)
				b.On("ApplyComponents").Return(b, nil)
				b.On("WriteToBMP").Return([]byte("a new bmp"), nil)
				l := loader{
					builder: b,
				}
				res, err := l.Write().ToBMPReader(nilProps)
				assert.Equal(t, bytes.NewReader([]byte("a new bmp")), res)
				assert.NoError(t, err)
			})
		})
	})
	b.AssertExpectations(t)
}

func TestPanics(t *testing.T) {
	l := loader{}
	nilProps := render.NamedProperties(nil)
	t.Run("FromBuilder", func(t *testing.T) {
		l2, props, err := l.FromBuilder(nil)
		assert.Nil(t, l2)
		assert.Equal(t, nilProps, props)
		assert.Error(t, err)
	})
	t.Run("FromBytes", func(t *testing.T) {
		l2, props, err := l.FromBytes(nil)
		assert.Nil(t, l2)
		assert.Equal(t, nilProps, props)
		assert.Error(t, err)
	})
	t.Run("FromFile", func(t *testing.T) {
		l2, props, err := l.FromFile("")
		assert.Nil(t, l2)
		assert.Equal(t, nilProps, props)
		assert.Error(t, err)
	})
	t.Run("FromJSON", func(t *testing.T) {
		l2, props, err := l.FromJSON(nil)
		assert.Nil(t, l2)
		assert.Equal(t, nilProps, props)
		assert.Error(t, err)
	})
	t.Run("FromReader", func(t *testing.T) {
		l2, props, err := l.FromReader(nil)
		assert.Nil(t, l2)
		assert.Equal(t, nilProps, props)
		assert.Error(t, err)
	})
	t.Run("ToBuilder", func(t *testing.T) {
		b, err := l.ToBuilder(nil)
		assert.Nil(t, b)
		assert.Error(t, err)
	})
	t.Run("ToBMP", func(t *testing.T) {
		raw, err := l.ToBMP(nil)
		assert.Nil(t, raw)
		assert.Error(t, err)
	})
	t.Run("ToCanvas", func(t *testing.T) {
		c, err := l.ToCanvas(nil)
		assert.Nil(t, c)
		assert.Error(t, err)
	})
	t.Run("ToImage", func(t *testing.T) {
		img, err := l.ToImage(nil)
		assert.Nil(t, img)
		assert.Error(t, err)
	})
	t.Run("ToBMPReader", func(t *testing.T) {
		rdr, err := l.ToBMPReader(nil)
		assert.Nil(t, rdr)
		assert.Error(t, err)
	})
}
