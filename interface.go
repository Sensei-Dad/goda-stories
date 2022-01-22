package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Globals related to the UI
const ElementBuffer int = 5
const WindowWidth, WindowHeight int = 640, 360

// If the order of characters in the first part of the image changes, this string needs to be updated
const fontChars string = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

type BitmapInterface struct {
	FontImage       *ebiten.Image
	FontImageWidth  int
	FontImageHeight int
	TileWidth       int
	TileHeight      int
	Charmap         map[rune]int
}

func NewBitmapInterface(fontImagePath string, tileW, tileH int) *BitmapInterface {
	ret := BitmapInterface{
		TileWidth:  tileW,
		TileHeight: tileH,
	}

	tImage, _, err := ebitenutil.NewImageFromFile(fontImagePath)
	if err != nil {
		log.Fatal(err)
	}

	ret.FontImage = tImage
	w, h := tImage.Size()
	ret.FontImageWidth = w / tileW
	ret.FontImageHeight = h / tileH
	ret.Charmap = make(map[rune]int, len(fontChars))

	for i := 0; i < len(fontChars); i++ {
		c := rune(fontChars[i])
		ret.Charmap[c] = i
	}

	return &ret
}

// TODO: Implement a target font size
func (ui *BitmapInterface) GetText(s string, clr color.RGBA) *ebiten.Image {
	ret := ebiten.NewImage(ui.TileWidth*len(s), ui.TileHeight)

	for i, char := range s {
		charNum := ui.Charmap[char]
		op := &ebiten.DrawImageOptions{}

		r := float64(clr.R) / 0xff
		g := float64(clr.G) / 0xff
		b := float64(clr.B) / 0xff
		op.ColorM.Scale(r, g, b, 1)
		op.GeoM.Translate(float64(i*ui.TileWidth), 0)

		srcX := (charNum % ui.FontImageWidth) * ui.TileWidth
		srcY := (charNum / ui.FontImageWidth) * ui.TileHeight
		ret.DrawImage(ui.FontImage.SubImage(image.Rect(srcX, srcY, srcX+ui.TileWidth, srcY+ui.TileHeight)).(*ebiten.Image), op)
	}

	return ret
}
