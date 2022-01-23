package main

import (
	"image/color"
	"math"

	"github.com/MasterShizzle/goda-stories/ebitileui"
	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World *GameWorld
	Gui   *ebitileui.BitmapInterface
	View  ViewCoords
	tick  int64
}

func NewGame(tileInfo []gosoh.TileInfo, zoneInfo []gosoh.ZoneInfo, itemInfo []gosoh.ItemInfo, puzzleInfo []gosoh.PuzzleInfo, creatureInfo []gosoh.CreatureInfo, soundList []string) *Game {
	// TODO: Distinguish between "init game" and "new game"
	g := &Game{}
	g.Gui = ebitileui.NewBitmapInterface("assets/font_16x20.png", 16, 20)

	// Viewport is 12:10 ratio
	vHeight := float64(ebitileui.WindowHeight - (2 * ebitileui.ElementBuffer))
	vWidth := math.Round(vHeight * ebitileui.ViewAspectRatio)

	g.View = ViewCoords{
		X:      0.0,
		Y:      0.0,
		Width:  vWidth,
		Height: vHeight,
	}

	gosoh.LoadAllTiles(tileInfo)

	gosoh.Zones = zoneInfo
	gosoh.Items = itemInfo
	gosoh.Puzzles = puzzleInfo
	gosoh.Creatures = creatureInfo
	gosoh.Sounds = soundList

	g.World = NewWorld()

	// ECS!
	gosoh.InitializeECS()

	g.tick = 0

	return g
}

var currentArea *gosoh.MapArea

func (g *Game) Update() error {
	// g.tick++
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
	// Draw the Viewport
	currentArea.DrawLayer(gosoh.TerrainLayer, screen, g.View.X, g.View.Y, g.View.Width, g.View.Height, float64(ebitileui.ElementBuffer))
	// TODO: Walls and Renderables need to be drawn at the same time
	currentArea.DrawLayer(gosoh.WallsLayer, screen, g.View.X, g.View.Y, g.View.Width, g.View.Height, float64(ebitileui.ElementBuffer))
	gosoh.ProcessRenderables(screen, g.View.X, g.View.Y, float64(ebitileui.ElementBuffer))
	currentArea.DrawLayer(gosoh.OverlayLayer, screen, g.View.X, g.View.Y, g.View.Width, g.View.Height, float64(ebitileui.ElementBuffer))

	splash := g.Gui.GetText("Hello, World!!", color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 1})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	screen.DrawImage(splash, op)

	// Show player stuff
	gosoh.ShowDebugInfo(screen, g.View.X, g.View.Y)
	gosoh.DrawEntityBoxes(screen, g.View.X, g.View.Y, float64(ebitileui.ElementBuffer))
}

func (g *Game) Layout(w, h int) (int, int) {
	// 640x360 internal dimensions, by default
	// 16:9 aspect ratio, with plenty of scaling
	return ebitileui.WindowWidth, ebitileui.WindowHeight
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
