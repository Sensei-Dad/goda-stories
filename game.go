package main

import (
	"image/color"

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

	// TODO: Center viewport function
	g.View = ViewCoords{
		X:          0.0,
		Y:          0.0,
		Width:      gosoh.ViewportWidth,
		Height:     gosoh.ViewportHeight,
		CurrentMap: 0,
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

func (g *Game) Update() error {
	// g.tick++
	gosoh.ProcessInput()
	// TODO: Handle AI, randomly move critters around, etc.
	// ProcessCreatures(g)
	gosoh.ProcessMovement(g.World.Maps[g.View.CurrentMap])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each layer, starting at the bottom
	g.World.DrawTerrain(screen, g.View)
	g.World.DrawObjects(screen, g.View)
	gosoh.ProcessRenderables(screen)
	g.World.DrawOverlay(screen, g.View)

	splash := g.Gui.GetText("Hello, World!!", color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 1})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	screen.DrawImage(splash, op)

	// Show player stuff
	gosoh.ShowDebugInfo(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	// for now, return the map with nothing else around it
	return (g.View.Width * gosoh.TileWidth), (g.View.Height * gosoh.TileHeight)
}
