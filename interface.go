package main

import (
	_ "image/png"
	"io/ioutil"

	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
)

// Globals related to the UI
const ElementBuffer int = 5
const WindowWidth, WindowHeight int = 640, 360
const ViewAspectRatio float64 = 1.2
const GuiFontFile string = "assets/Cloude_Regular_Bold_1.02.ttf"

func loadButtonImage() (*widget.ButtonImage, error) {
	idle, err := loadNineSlice("assets/button_idle.png", [3]int{16, 16, 16}, [3]int{16, 16, 16})
	if err != nil {
		return nil, err
	}

	hover, err := loadNineSlice("assets/button_hover.png", [3]int{16, 16, 16}, [3]int{16, 16, 16})
	if err != nil {
		return nil, err
	}

	pressed, err := loadNineSlice("assets/button_pressed.png", [3]int{16, 16, 16}, [3]int{16, 16, 16})
	if err != nil {
		return nil, err
	}

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadNineSlice(path string, w [3]int, h [3]int) (*image.NineSlice, error) {
	i, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}

	return image.NewNineSlice(i, w, h), nil
}

func loadFont(path string, size float64) (font.Face, error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
