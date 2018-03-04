package sprites

import (
	"fmt"
	"image"
	"image/color"
	"io"

	fcolor "github.com/fatih/color"
	"github.com/pkg/errors"
)

const (
	headerSize     = 16
	spriteSize     = 16
	prgBankSize    = 16 * 1024
	chrBankSize    = 8 * 1024
	spritesPerBank = 512
)

// A Sprite is represented by 16 bytes.
type Sprite []byte

// Colors returns the colors from 0 to 3 in an byte array.
func (s Sprite) Colors() []byte {
	cs := []byte{}
	for i := 0; i < 8; i++ {
		c1 := s[i]
		c2 := s[i+8]
		// Use an unsigned integer and check for <= 7 instead of using an integer
		// and casting it on every iteration.
		for j := uint(7); j <= 7; j-- {
			v := (c1 >> j & 1) + (c2>>j)&1<<1
			cs = append(cs, v)
		}
	}
	return cs
}

var colorFuncs = []func(a ...interface{}){
	fcolor.New(fcolor.BgBlack).PrintFunc(),
	fcolor.New(fcolor.BgRed).PrintFunc(),
	fcolor.New(fcolor.BgGreen).PrintFunc(),
	fcolor.New(fcolor.BgBlue).PrintFunc(),
}

// Print the sprite to the terminal.
// Each row is defined by summing the nth and the n+8th byte.
// n:   11110000
// n+8: 11001100
//      33112200
func (s Sprite) Print() {
	for i, c := range s.Colors() {
		if i > 0 && i%8 == 0 {
			fmt.Println()
		}
		colorFuncs[c](c)
	}
	fmt.Println()
}

var colors = []color.Color{
	color.RGBA{},
	color.RGBA{R: 255, A: 255},
	color.RGBA{G: 255, A: 255},
	color.RGBA{B: 255, A: 255},
}

// Image returns an image representation of the sprite.
func (s Sprite) Image() *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for i, v := range s.Colors() {
		im.Set(i%8, i/8, colors[v])
	}
	return im
}

// Read a iNES ROM.
func Read(r io.ReadSeeker) ([]Sprite, error) {
	h := make([]byte, headerSize)
	if _, err := r.Read(h); err != nil {
		return nil, err
	}

	if h[0] != 0x4e || h[1] != 0x45 || h[2] != 0x53 || h[3] != 0x1a {
		return nil, errors.New("first 4 bites are expected to be NES\\0x01a")
	}
	pc := int(h[4])
	cc := int(h[5])
	if cc == 0 {
		return nil, errors.New("CHR RAM is used")
	}

	// Seek to the first CHR ROM bank.
	r.Seek(headerSize+int64(pc)*prgBankSize, io.SeekStart)
	sprs := []Sprite{}
	for i := 0; i < cc; i++ {
		for j := 0; j < chrBankSize/spriteSize; j++ {
			s := make(Sprite, spriteSize)
			if _, err := r.Read(s); err != nil {
				if err == io.EOF {
					break
				}
				return nil, errors.Wrapf(err, "could not read sprite %d", j)
			}
			sprs = append(sprs, s)
		}
	}
	return sprs, nil
}
