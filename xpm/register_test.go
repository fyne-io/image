package xpm_test

import (
	"image"
	_ "image/png"
	"os"
	"testing"

	_ "github.com/fyne-io/image/xpm" // xpm image parser

	"github.com/stretchr/testify/assert"
)

func TestRegister_Read(t *testing.T) {
	r, err := os.Open("testdata/blarg.xpm")
	assert.Nil(t, err)
	img, fmt, err := image.Decode(r)
	assert.Nil(t, err)
	assert.Equal(t, "xpm", fmt)

	assert.Equal(t, 0, img.Bounds().Min.X)
	assert.Equal(t, 0, img.Bounds().Min.Y)
	assert.Equal(t, 16, img.Bounds().Max.X)
	assert.Equal(t, 7, img.Bounds().Max.Y)
}
