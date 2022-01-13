package main

import (
	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World *GameWorld
	Gui   *BitmapInterface
	View  ViewCoords
	tick  int64
}

func NewGame(tileInfo []gosoh.TileInfo, zoneInfo []gosoh.ZoneInfo, itemInfo []gosoh.ItemInfo, puzzleInfo []gosoh.PuzzleInfo, creatureInfo []gosoh.CreatureInfo) *Game {
	// TODO: Distinguish between "init game" and "new game"
	g := &Game{}
	g.Gui = NewBitmapInterface("assets/font_16x20.png", 16, 20)

	g.View = ViewCoords{
		X:      0.0,
		Y:      0.0,
		Width:  gosoh.ViewportWidth,
		Height: gosoh.ViewportHeight,
	}

	gosoh.LoadAllTiles(tileInfo)

	gosoh.Zones = zoneInfo
	gosoh.Items = itemInfo
	gosoh.Puzzles = puzzleInfo
	gosoh.Creatures = creatureInfo

	g.World = NewWorld()

	// ECS!
	gosoh.InitializeECS()

	g.tick = 0

	return g
}

var currentArea gosoh.MapArea

func (g *Game) Update() error {
	// g.tick++
	gosoh.ProcessInput()
	// TODO: Handle AI, randomly move critters around, etc.
	// ProcessCreatures(g)
	currentArea = g.GetCurrentArea()
	gosoh.ProcessMovement(currentArea)
	g.CenterViewport(currentArea)
	// if the player has moved, then check loading / unloading Entities
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the map and entities
	currentArea = g.GetCurrentArea()
	currentArea.DrawTerrain(screen, g.View.X, g.View.Y, g.View.Width, g.View.Height)
	// g.World.DrawWalls(screen, g.View)
	gosoh.ProcessRenderables(screen, g.View.X, g.View.Y)
	// g.World.DrawOverlay(screen, g.View)

	// splash := g.Gui.GetText("Hello, World!!", color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 1})
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(10, 10)
	// screen.DrawImage(splash, op)

	// Show player stuff
	gosoh.ShowDebugInfo(screen, g.View.X, g.View.Y)
	gosoh.DrawEntityBoxes(screen, g.View.X, g.View.Y)
}

func (g *Game) Layout(w, h int) (int, int) {
	// for now, return the map with nothing else around it
	return (g.View.Width * gosoh.TileWidth), (g.View.Height * gosoh.TileHeight)
}

func (g *Game) CenterViewport(a gosoh.MapArea) {
	pX, pY := gosoh.GetPlayerCoords()

	halfWidth := float64(g.View.Width * (gosoh.TileWidth / 2))
	halfHeight := float64(g.View.Height * (gosoh.TileHeight / 2))

	maxX := float64(a.Width * 18 * gosoh.TileWidth)
	maxY := float64(a.Height * 18 * gosoh.TileHeight)

	// Center on the player wherever possible
	g.View.X = pX - halfWidth
	g.View.X = gosoh.ClampFloat(g.View.X, 0, maxX)
	g.View.Y = pY - halfHeight
	g.View.Y = gosoh.ClampFloat(g.View.Y, 0, maxY)
}
