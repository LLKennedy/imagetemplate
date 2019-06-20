package imagetemplate

import (
	"testing"

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

func TestFromBuilder(t *testing.T) {

}
