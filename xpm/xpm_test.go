package icon

import (
	"image"
	"image/color"
	_ "image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	r, err := os.Open("testdata/blarg.xpm")
	assert.Nil(t, err)
	defer r.Close()
	img, err := parseXPM(r)
	assert.Nil(t, err)
	assert.Equal(t, 0, img.Bounds().Min.X)
	assert.Equal(t, 0, img.Bounds().Min.Y)
	assert.Equal(t, 16, img.Bounds().Max.X)
	assert.Equal(t, 7, img.Bounds().Max.Y)

	r, err = os.Open("testdata/blarg.png")
	if err != nil {
		t.Error(err)
	}

	golden, _, err := image.Decode(r)
	if err != nil {
		t.Error(err)
	}

	pixCount := len(golden.(*image.RGBA).Pix)
	assert.Equal(t, pixCount, len(img.(*image.NRGBA).Pix))
	for i := 0; i < pixCount; i++ {
		assert.Equal(t, golden.(*image.RGBA).Pix[i], img.(*image.NRGBA).Pix[i])
	}
}

func TestParseColor(t *testing.T) {
	id, c, err := parseColor(". c #000000", 1)
	assert.Nil(t, err)
	assert.Equal(t, ".", id)
	assert.Equal(t, &color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff}, c)

	id, c, err = parseColor("  c #000000", 1) // special case, id is spaces
	assert.Nil(t, err)
	assert.Equal(t, " ", id)
	assert.Equal(t, &color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff}, c)
}

func TestParseDimensions(t *testing.T) {
	w, h, i, j, err := parseDimensions("5 10 2 1")
	assert.Nil(t, err)
	assert.Equal(t, 5, w)
	assert.Equal(t, 10, h)
	assert.Equal(t, 2, i)
	assert.Equal(t, 1, j)
}

func TestStringToColor(t *testing.T) {
	c, err := stringToColor("None")
	assert.Nil(t, err)
	assert.Equal(t, color.Transparent, c)
	c, err = stringToColor("#000000")
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{A: 0xff}, c)
	c, err = stringToColor("#ffffff")
	assert.Nil(t, err)
	assert.Equal(t, &color.NRGBA{0xff, 0xff, 0xff, 0xff}, c)
}

func TestStripQuotes(t *testing.T) {
	assert.Equal(t, "hello", stripQuotes("\"hello\""))
	assert.Equal(t, "hello", stripQuotes("\"hello\","))
	assert.Equal(t, "", stripQuotes("\"\""))
}
