package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World *GameWorld
	Gui   *ebitenui.UI
	View  ViewCoords
	tick  int64
}

func NewGame(tileset *ebiten.Image, tileInfo []gosoh.TileInfo, zoneInfo []gosoh.ZoneInfo, itemInfo []gosoh.ItemInfo, puzzleInfo []gosoh.PuzzleInfo, creatureInfo []gosoh.CreatureInfo, soundList []string) *Game {
	// TODO: Distinguish between "init game" and "new game"
	g := &Game{}

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
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Make a BUTTON!
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImg),
		widget.ButtonOpts.Text("This button does nothing!", guiFont, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  10,
			Right: 10,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			fmt.Println("Button clicked")
		}),
	)

	gameContainer.AddChild(button)

	ui := ebitenui.UI{
		Container: gameContainer,
	}

	g.Gui = &ui

	// Viewport is 12:10 ratio
	vHeight := float64(WindowHeight - (2 * ElementBuffer))
	vWidth := math.Round(vHeight * ViewAspectRatio)

	g.View = ViewCoords{
		X:      0.0,
		Y:      0.0,
		Width:  vWidth,
		Height: vHeight,
	}

	gosoh.Zones = zoneInfo
	gosoh.Items = itemInfo
	gosoh.Puzzles = puzzleInfo
	gosoh.Creatures = creatureInfo
	gosoh.Sounds = soundList
	gosoh.TilesetImage = tileset
	gosoh.TileInfos = tileInfo

	g.World = NewWorld()

	// ECS!
	gosoh.InitializeECS()

	g.tick = 0

	return g
}

var currentArea *gosoh.MapArea

func (g *Game) Update() error {
	// g.tick++
	g.Gui.Update()
	gosoh.ProcessInput()
	// TODO: Handle AI, randomly move critters around, etc.
	// ProcessCreatures(g)
	currentArea = g.World.GetCurrentArea()
	gosoh.ProcessMovement(currentArea)
	g.CenterViewport(currentArea)
	// if the player has moved, then check loading / unloading Entities
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Gui.Draw(screen)
	// Draw the Viewport
	currentArea.DrawLayer(gosoh.TerrainLayer, screen, g.View.X, g.View.Y, g.View.Width, g.View.Height, float64(ElementBuffer))
	// TODO: Walls and Renderables need to be interleaved and drawn at the same time
	currentArea.DrawLayer(gosoh.WallsLayer, screen, g.View.X, g.View.Y, g.View.Width, g.View.Height, float64(ElementBuffer))
	gosoh.ProcessRenderables(screen, g.View.X, g.View.Y, float64(ElementBuffer))
	currentArea.DrawLayer(gosoh.OverlayLayer, screen, g.View.X, g.View.Y, g.View.Width, g.View.Height, float64(ElementBuffer))

	// Show player stuff
	gosoh.ShowDebugInfo(screen, g.View.X, g.View.Y)
	gosoh.DrawEntityBoxes(screen, g.View.X, g.View.Y, float64(ElementBuffer))
}

func (g *Game) Layout(w, h int) (int, int) {
	// 640x360 internal dimensions, by default
	// 16:9 aspect ratio, with plenty of scaling
	return WindowWidth, WindowHeight
}

func (g *Game) CenterViewport(a *gosoh.MapArea) {
	pX, pY, _, _ := gosoh.GetPlayerCoords()

	halfWidth := g.View.Width / 2
	halfHeight := g.View.Height / 2

	maxX := float64(a.Width * 18 * gosoh.TileWidth)
	maxY := float64(a.Height * 18 * gosoh.TileHeight)

	// Center on the player wherever possible
	g.View.X = pX - halfWidth
	g.View.X = gosoh.ClampFloat(g.View.X, 0, maxX)
	g.View.Y = pY - halfHeight
	g.View.Y = gosoh.ClampFloat(g.View.Y, 0, maxY)
}
