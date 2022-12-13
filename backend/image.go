package backend

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/tiff"
)

func NewImageCompress(mineType string) (ImageCompress, error) {
	mineType = strings.ToLower(mineType)
	switch mineType {
	case "jpg", "jpeg", ".jpg", ".jpeg":
		return &JpegCompress{}, nil
	case "png", ".png":
		return &PngCompress{}, nil
	}
	return nil, errors.New("no match image compress: " + mineType)
}

type Options struct {
	Quality int `json:"quality"`
}

type ImageCompress interface {
	Decode(r io.Reader) (image.Image, error)
	Encode(i image.Image, o *Options) (buf bytes.Buffer, err error)
}

type PngCompress struct {
}

func (comp PngCompress) Decode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}

func (comp PngCompress) Encode(i image.Image, o *Options) (buf bytes.Buffer, err error) {
	err = tiff.Encode(&buf, i, &tiff.Options{
		Compression: tiff.Deflate,
	})
	return buf, err
}

type JpegCompress struct {
}

func (comp JpegCompress) Decode(r io.Reader) (image.Image, error) {
	return jpeg.Decode(r)
}

func (comp JpegCompress) Encode(i image.Image, o *Options) (buf bytes.Buffer, err error) {
	err = jpeg.Encode(&buf, i, &jpeg.Options{Quality: o.Quality})
	return buf, err
}
