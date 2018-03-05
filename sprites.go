package sprites

import (
	"image"
	"image/color"
	"io"

	"github.com/pkg/errors"
)

const (
	spriteWidth = 8

	headerSize     = 16
	spriteSize     = 16
	prgBankSize    = 16 * 1024
	chrBankSize    = 8 * 1024
	spritesPerBank = 512
)

// A Sprite is represented by 16 bytes.
type Sprite []byte

// Colors returns the colors from 0 to 3 in an byte array.
// Each row is defined by summing the nth and the n+8th byte.
// n:   11110000b
// n+8: 11001100b
//      33112200d
func (s Sprite) Colors() []byte {
	cs := make([]byte, spriteWidth*spriteWidth)
	for i := uint(0); i < spriteWidth; i++ {
		c1 := s[i]
		c2 := s[i+spriteWidth]
		for j := uint(0); j < spriteWidth; j++ {
			v := (c1 >> j & 1) + (c2>>j)&1<<1
			x := spriteWidth - j - 1
			y := i * spriteWidth
			cs[x+y] = v
		}
	}
	return cs
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
