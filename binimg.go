package binimg

import (
	"image"
	"image/color"
	"io"
	"math"
)

//BinImg an image with one-byte HSL colors (see type HSL)
//when decoded from a file the first two bytes are the width and height
type BinImg struct {
	h, w uint8
	pix  []uint8
}

//At returns the color of the pixel at (x, y)
func (b BinImg) At(x, y int) color.Color {
	colour := b.pix[x+y*int(b.w)]

	return FromByte(colour)
}

//Bounds returns the domain for which At is valid.
func (b BinImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(b.w), int(b.h))
}

//ColorModel returns the Image's color model.
func (b BinImg) ColorModel() color.Model {
	return HSLModel
}

func mDecode(r io.Reader) (image.Image, error) {
	sz := make([]uint8, 2)
	_, err := r.Read(sz)
	if err != nil {
		return nil, err
	}

	img := BinImg{
		w: sz[0],
		h: sz[1],
	}

	img.pix = make([]uint8, int(img.w)*int(img.h))
	_, err = r.Read(img.pix)
	return img, err
}

func mEncode(w io.Writer, img image.Image) error {
	b := img.Bounds()
	data := []uint8{uint8(math.Min(float64(b.Dx()), 255)), uint8(math.Min(float64(b.Dy()), 255))}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := HSLModel.Convert(img.At(x, y)).(HSL)
			data = append(data, c.Byte())
		}
	}
	_, err := w.Write(data)
	return err
}
