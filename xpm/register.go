package icon

import (
	"image"
	"io"
	"io/ioutil"
)

func init() {
	image.RegisterFormat("xpm", "/* XPM */", Decode, nil)
	image.RegisterFormat("xpm", "static char", Decode, nil)
}

func Decode(r io.Reader) (image.Image, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parseXPM(data), nil
}
