package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"git.192k.pw/bake/nes-sprites"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: nes-sprites [rom.nes]")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sprs, err := sprites.Read(f)
	if err != nil {
		log.Fatal(err)
	}

	spriteSize := 8
	spritesPerRow := 16
	w := spritesPerRow * spriteSize
	h := len(sprs) / spritesPerRow * spriteSize
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	for i, s := range sprs {
		x := i % spritesPerRow * spriteSize
		y := i / spritesPerRow * spriteSize
		r := image.Rect(x, y, x+spriteSize, y+spriteSize)
		draw.Draw(im, r, s.Image(), image.Point{}, draw.Src)
	}

	f, err = os.Create("sprites.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	png.Encode(f, im)
}
