package xpm

import (
	"image"
	"io"
)

func init() {
	image.RegisterFormat("xpm", "/* XPM */", Decode, DecodeConfig)
	image.RegisterFormat("xpm", "static char", Decode, DecodeConfig)
}

func Decode(r io.Reader) (image.Image, error) {
	return parseXPM(r)
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	return parseXPMConfig(r)
}
