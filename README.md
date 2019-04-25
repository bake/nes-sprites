# nes-sprites

Extract sprites from iNES ROM files. Sprites are represented by 16 bytes which are splitted in two channels (the first 8 bytes and the second 8 bytes). Adding the bits from the *n*th byte of each channel results in a pixels color index (from 0 to 3).

An example implementation can be found in [/cmd/nes-sprites/main.go](/cmd/nes-sprites/main.go). The following sprites are extracted from Super Mario Bros.

```bash
$ go install github.com/bake/nes-sprites/...
$ nes-sprites mariobros1.nes
$
```

![Sprites](/sprites.png)

The format is described at [Sadistech](https://sadistech.com/nesromtool/romdoc.html) and [N3S](https://n3s.io/index.php?title=How_It_Works). Note that this package only reads CHR ROM and not CHR RAM.
