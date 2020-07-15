# binimg
A golang image library with one-byte "HSL" colors

## Format
An image could be any file, the first two bytes are width and height, the following bytes describe the pixels as "HSL" color one byte per pixel.

The conversion is decribed in the docs [here](https://godoc.org/github.com/Strosel/binimg#HSL) and a file of [0x10, 0x10, 0x00 .. 0xFF] gives this palette

![](https://raw.githubusercontent.com/Strosel/binimg/master/palette_big.png)
