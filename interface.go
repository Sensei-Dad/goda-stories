package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
)

// Globals related to the UI
const ElementBuffer int = 5
const WindowWidth, WindowHeight int = 1280, 720
const ViewAspectRatio float64 = 1.2
const GuiFontFile string = "assets/Cloude_Regular_Bold_1.02.ttf"

func buildGui() *ebitenui.UI {
	// Build UI elements
	buttonImg, err := loadButtonImage()
	if err != nil {
		log.Fatal(err)
	}

	guiFont, err := loadFont(GuiFontFile, 32)
	if err != nil {
		log.Fatal(err)
	}
	defer guiFont.Close()

	// Root UI container
	gameContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{0, 0, 0, 0})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top:    ElementBuffer,
				Bottom: ElementBuffer,
			}),
		)),
	)

	// Make a BUTTON!
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
		),
		widget.ButtonOpts.Image(buttonImg),
		widget.ButtonOpts.Text("Menu", guiFont, &widget.ButtonTextColor{
			Idle: color.RGBA{0xd2, 0xdb, 0xe0, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  2 * ElementBuffer,
			Right: 2 * ElementBuffer,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			fmt.Println("Button clicked")
		}),
	)

	gameContainer.AddChild(button)

	return &ebitenui.UI{
		Container: gameContainer,
	}
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle, err := loadNineSlice("assets/button_idle.png", [3]int{20, 20, 20}, [3]int{13, 4, 13})
	if err != nil {
		return nil, err
	}

	hover, err := loadNineSlice("assets/button_hover.png", [3]int{20, 20, 20}, [3]int{13, 4, 13})
	if err != nil {
		return nil, err
	}

	pressed, err := loadNineSlice("assets/button_pressed.png", [3]int{20, 20, 20}, [3]int{13, 4, 13})
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
