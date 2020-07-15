package binimg

import (
	"image/color"
	"math"
)

var (
	//LegalAngles are possible HSL H angles
	//note: these angles were picked somewhat arbitrarily by the author
	LegalAngles = []float64{0, 21, 41, 62, 83, 103, 124, 145, 165, 186, 207, 227, 248, 269, 289, 310}

	//HSLModel is the HSL color model
	HSLModel = color.ModelFunc(model)
)

//HSL Describes a one-byte "HSL" color
//the least significant four bits decide L
//if L = 0, L is instead decided by the most significant four bits to give a grayscale option with H & S = 0
//if L > 0, the most significant four bits decide H and S = 1
//H is limited to angles in LegalAngles and conversion uses the closest one
type HSL struct {
	h, s, l float64
}

//FromByte converts a byte to an HSL color
func FromByte(b byte) HSL {
	if (b & 0x0F) == 0 {
		return HSL{
			l: float64((b&0xF0)>>4) / 15.,
		}
	}
	return HSL{
		l: float64(b&0x0F) / 17.,
		s: 1,
		h: math.Round(float64((b&0xF0)>>4) * 310. / 15.),
	}
}

//Byte converts an HSL color into a byte
func (h HSL) Byte() (b byte) {
	if h.s == 1 {
		b |= byte(h.l*17) & 0x0F
		for i, v := range LegalAngles {
			if v == h.h {
				b |= byte((i) << 4)
				break
			}
		}
		// b |= (byte(h.h*15./310.) & 0x0F) << 4
	} else {
		b |= (byte(h.l*15.) & 0x0F) << 4
	}
	return
}

//RGBA returns the alpha-premultiplied red, green, blue and alpha values
func (h HSL) RGBA() (r, g, b, a uint32) {
	C := (1 - math.Abs(2*h.l-1)) * h.s
	X := C * (1 - math.Abs(math.Mod(h.h/60, 2)-1))
	m := h.l - C/2

	if h.h < 60 {
		r = uint32((C + m) * 255)
		g = uint32((X + m) * 255)
		b = uint32(m * 255)
	} else if h.h < 120 {
		r = uint32((X + m) * 255)
		g = uint32((C + m) * 255)
		b = uint32(m * 255)
	} else if h.h < 180 {
		r = uint32(m * 255)
		g = uint32((C + m) * 255)
		b = uint32((X + m) * 255)
	} else if h.h < 240 {
		r = uint32(m * 255)
		g = uint32((X + m) * 255)
		b = uint32((C + m) * 255)
	} else if h.h < 300 {
		r = uint32((X + m) * 255)
		g = uint32(m * 255)
		b = uint32((C + m) * 255)
	} else {
		r = uint32((C + m) * 255)
		g = uint32(m * 255)
		b = uint32((X + m) * 255)
	}
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a = 0xFFFF

	return
}

func model(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	R := float64(r>>8) / 255.
	G := float64(g>>8) / 255.
	B := float64(b>>8) / 255.

	cmax := math.Max(R, math.Max(G, B))
	cmin := math.Min(R, math.Min(G, B))
	delta := cmax - cmin

	hsl := HSL{
		l: (cmax + cmin) / 2.,
	}
	hsl.s = 1
	if delta == 0 {
		hsl.h = 0
		hsl.s = 0
		return hsl
	} else if cmax == R {
		hsl.h = 60 * math.Mod((G-B)/delta, 6)
	} else if cmax == G {
		hsl.h = 60 * ((B-R)/delta + 2)
	} else if cmax == B {
		hsl.h = 60 * ((R-G)/delta + 4)
	}

	minf := 2.
	mind := 2.
	for i := 1; i < 17; i++ {
		f := float64(i) / 17.
		d := math.Abs(hsl.l - f)
		if d < mind {
			minf = f
			mind = d
		}
	}
	hsl.l = minf

	if hsl.h < 0 {
		hsl.h += 360.
	}
	min := 100.
	mini := -1
	for i, c := range LegalAngles {
		d := math.Abs(hsl.h - c)
		if d < min {
			min = d
			mini = i
		}
	}
	hsl.h = LegalAngles[mini]

	return hsl
}
