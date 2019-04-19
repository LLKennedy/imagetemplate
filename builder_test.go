package imagetemplate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	dimensions := []int{-10000, -500, -33, -1, 0, 1, 33, 500, 10000}
	for _, x := range dimensions {
		for _, y := range dimensions {
			t.Run(fmt.Sprintf("Creating builder with dimensions %d by %d", x, y), func(t *testing.T) {
				var newBuilder Builder
				newBuilder, err := NewBuilder(x, y)
				if x <= 0 && y <= 0 {
					assert.Nil(t, newBuilder)
					assert.Equal(t, "Invalid width and height", err.Error())
					return
				}
				if x <= 0 {
					assert.Nil(t, newBuilder)
					assert.Equal(t, "Invalid width", err.Error())
					return
				}
				if y <= 0 {
					assert.Nil(t, newBuilder)
					assert.Equal(t, "Invalid height", err.Error())
					return
				}
				assert.NotNil(t, newBuilder)
				realBuilder, ok := newBuilder.(*ImageBuilder)
				assert.True(t, ok)
				assert.NotNil(t, realBuilder)
				imageBounds := realBuilder.GetUnderlyingImage().Bounds()
				assert.Equal(t, imageBounds.Size().X, x)
				assert.Equal(t, imageBounds.Size().Y, y)
			})
		}
	}

}
