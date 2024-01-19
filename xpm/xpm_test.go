package xpm

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
	// Syntax is: <chars> {<key> <color>}+ (XPM Manual Chapter 2)
	for _, tt := range []struct {
		input     string
		wantID    string
		wantColor color.Color
	}{
		{
			input:     ". c #000000",
			wantID:    ".",
			wantColor: color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
		},
		{
			input:     "  c #000000",
			wantID:    " ", // special case, id is spaces
			wantColor: color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
		},
		{
			// three-digit hex color
			input:     "O c #123",
			wantID:    "O",
			wantColor: color.NRGBA{R: 0x10, G: 0x20, B: 0x30, A: 0xff},
		},
		{
			// color referenced by X11 color name
			input:     "r c red",
			wantID:    "r",
			wantColor: color.NRGBA{R: 0xff, G: 0x0, B: 0x0, A: 0xff},
		},
		{
			// a multi-word color name
			input:     "g c dark slate grey",
			wantID:    "g",
			wantColor: color.NRGBA{R: 47, G: 79, B: 79, A: 0xff},
		},
		{
			// "c" visual is not the first
			input:     "r g gray g4 #888888 c red m black",
			wantID:    "r",
			wantColor: color.NRGBA{R: 0xff, G: 0x0, B: 0x0, A: 0xff},
		},
		{
			// "c" visual is not the first, some colors have multiple words
			input:     "g g dark slate gray c pale green m black",
			wantID:    "g", // special case, id is spaces
			wantColor: color.NRGBA{R: 152, G: 251, B: 152, A: 0xff},
		},
	} {
		id, c, err := parseColor(tt.input, 1)
		assert.Nil(t, err, "parseColor(%q, 1): error: %v", tt.input, err)
		assert.Equal(t, tt.wantID, id, "parseColor(%q, 1): id %q, want %q", tt.input, id, tt.wantID)
		assert.Equal(t, tt.wantColor, c, "parseColor(%q, 1): color %v, want %v", tt.input, c, tt.wantColor)
	}
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
	for _, tt := range []struct {
		input string
		want  color.Color
	}{
		{
			input: "None",
			want:  color.Transparent,
		},
		{
			input: "#000000",
			want:  color.NRGBA{A: 0xff},
		},
		{
			input: "#ffffff",
			want:  color.NRGBA{0xff, 0xff, 0xff, 0xff},
		},
		{
			input: "#fff",
			want:  color.NRGBA{0xf0, 0xf0, 0xf0, 0xff},
		},
		{
			input: "red",
			want:  color.NRGBA{0xff, 0x00, 0x00, 0xff},
		},
	} {
		c, err := stringToColor(tt.input)
		assert.Nil(t, err)
		assert.Equal(t, tt.want, c, "stringToColor(%q) = %+v, want %+v", tt.input, c, tt.want)
	}
}

func TestStripQuotes(t *testing.T) {
	assert.Equal(t, "hello", stripQuotes("\"hello\""))
	assert.Equal(t, "hello", stripQuotes("\"hello\","))
	assert.Equal(t, "", stripQuotes("\"\""))
}
