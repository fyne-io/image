package xpm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"
)

// maxPixels is the maximum number of pixels that the parser supports,
// to protect against prohibitively large memory allocations.  XPM
// pixmaps tend to be small, it should not be possible to run into
// this limit with real-life data.
const maxPixels = 1024 * 1024 * 1024

// ErrInvalidFormat indicates that the input image was malformatted.
var ErrInvalidFormat = errors.New("invalid format")

func parseXPM(data io.Reader) (image.Image, error) {
	// Specification: https://www.xfree86.org/current/xpm.pdf
	var colCount, charSize int
	colors := make(map[string]color.Color)
	var img *image.NRGBA

	rowNum := 0
	scan := bufio.NewScanner(data)
	for scan.Scan() {
		row := scan.Text()
		if row == "" || row[0] != '"' {
			continue
		}
		row = stripQuotes(row)

		if rowNum == 0 {
			w, h, cols, size, err := parseDimensions(row)
			if err != nil {
				return nil, err
			}
			img = image.NewNRGBA(image.Rectangle{image.Point{}, image.Point{w, h}})
			colCount = cols
			charSize = size
		} else if rowNum <= colCount {
			id, c, err := parseColor(row, charSize)
			if err != nil {
				return nil, err
			}

			if id != "" {
				colors[id] = c
			}
		} else {
			err := parsePixels(row, charSize, rowNum-colCount-1, colors, img)
			if err != nil {
				return nil, err
			}
		}
		rowNum++
	}
	return img, scan.Err()
}

func parseColor(data string, charSize int) (id string, c color.Color, err error) {
	if len(data) < charSize {
		return "", nil, fmt.Errorf("%w: missing color specification", ErrInvalidFormat)
	}

	id = data[:charSize]
	parts := strings.Fields(data[charSize:])
	for len(parts) >= 2 {
		key := parts[0]
		parts = parts[1:]

		nki := nextKeyIndex(parts)
		color := strings.Join(parts[:nki], " ")
		parts = parts[nki:]

		if color == "" {
			return "", nil, fmt.Errorf("%w: missing color specification", ErrInvalidFormat)
		}

		switch key {
		case "c":
			c, err := stringToColor(color)
			return id, c, err
		case "m", "s", "g4", "g":
			// We don't support mono, symbolic, and
			// grayscale visuals.
			continue
		default:
			return "", nil, fmt.Errorf("unknown visual %q", key)
		}
	}
	return "", nil, fmt.Errorf("%w: missing color specification", ErrInvalidFormat)
}

// nextKeyIndex returns the index of the next "c", "m", s", "g4", or
// "g", or otherwise len(parts).
func nextKeyIndex(parts []string) int {
	for i, p := range parts {
		switch p {
		case "c", "m", "s", "g4", "g":
			return i
		}
	}
	return len(parts)
}

func parseDimensions(data string) (w, h, ncolors, cpp int, err error) {
	if len(data) == 0 {
		return
	}
	parts := strings.Split(data, " ")
	if len(parts) != 4 {
		return
	}

	w, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	h, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	if w*h <= 0 {
		err = fmt.Errorf("%w: empty or negative-sized image (%v x %v)", ErrInvalidFormat, w, h)
		return
	}
	if w*h >= maxPixels {
		err = fmt.Errorf("%w: too many pixels (%v x %v), want < %v", ErrInvalidFormat, w, h, maxPixels)
		return
	}
	ncolors, err = strconv.Atoi(parts[2])
	if err != nil {
		return
	}
	if ncolors <= 0 {
		err = fmt.Errorf("%w: ncolors <= 0: missing color palette", ErrInvalidFormat)
		return
	}
	cpp, err = strconv.Atoi(parts[3])
	if err != nil {
		return
	}
	if cpp <= 0 {
		err = fmt.Errorf("%w: characters per pixel <= 0", ErrInvalidFormat)
		return
	}
	return
}

func parsePixels(row string, charSize int, pixRow int, colors map[string]color.Color, img *image.NRGBA) error {
	if len(row) < charSize*(img.Stride/4) {
		return fmt.Errorf("%w: missing pixel data", ErrInvalidFormat)
	}
	off := pixRow * img.Stride
	if len(img.Pix) < off+img.Stride {
		return fmt.Errorf("%w: too much pixel data", ErrInvalidFormat)
	}
	chPos := 0
	for i := 0; i < img.Stride/4; i++ {
		id := row[chPos : chPos+charSize]
		c, ok := colors[id]
		if !ok {
			c = color.Transparent
		}

		pos := off + (i * 4)
		r, g, b, a := c.RGBA()
		img.Pix[pos] = uint8(r)
		img.Pix[pos+1] = uint8(g)
		img.Pix[pos+2] = uint8(b)
		img.Pix[pos+3] = uint8(a)
		chPos += charSize
	}
	return nil
}

func stringToColor(data string) (color.Color, error) {
	if strings.EqualFold("none", data) {
		return color.Transparent, nil
	}

	switch data[0] {
	case '#':
		var (
			r, g, b uint8
			err     error
		)
		switch len(data) {
		case 7:
			_, err = fmt.Sscanf(data, "#%02x%02x%02x", &r, &g, &b)
		case 4:
			_, err = fmt.Sscanf(data, "#%01x%01x%01x", &r, &g, &b)
			// In X11 color specs, #fff == #f0f0f0
			// See https://gitlab.freedesktop.org/xorg/lib/libxpm/-/issues/7
			r, g, b = 0x10*r, 0x10*g, 0x10*b
		default:
			return nil, fmt.Errorf("%w: invalid hex color %q", ErrInvalidFormat, data)
		}
		return color.NRGBA{r, g, b, 0xff}, err
	default:
		c, ok := x11colors[data]
		if !ok {
			return nil, fmt.Errorf("%w: invalid X11 color %q", ErrInvalidFormat, data)
		}
		return c, nil
	}
}

func stripQuotes(data string) string {
	if len(data) == 0 || data[0] != '"' {
		return data
	}

	end := strings.Index(data[1:], "\"")
	if end == -1 {
		return data[1:]
	}
	return data[1 : end+1]
}
