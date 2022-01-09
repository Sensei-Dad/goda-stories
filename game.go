package main

import (
	"math"

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
		X:           0.0,
		Y:           0.0,
		Width:       gosoh.ViewportWidth,
		Height:      gosoh.ViewportHeight,
		CurrentArea: 0,
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
	ms := g.World.SubAreas[g.View.CurrentArea]
	gosoh.ProcessMovement(ms)
	g.CenterViewport(ms.Width, ms.Height)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the map and entities
	g.World.DrawTerrain(screen, g.View)
	g.World.DrawWalls(screen, g.View)
	gosoh.ProcessRenderables(screen)
	g.World.DrawOverlay(screen, g.View)

	// splash := g.Gui.GetText("Hello, World!!", color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 1})
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(10, 10)
	// screen.DrawImage(splash, op)

	// Show player stuff
	gosoh.ShowDebugInfo(screen, g.View.X, g.View.Y)
}

func (g *Game) Layout(w, h int) (int, int) {
	// for now, return the map with nothing else around it
	return (g.View.Width * gosoh.TileWidth), (g.View.Height * gosoh.TileHeight)
}

func (g *Game) CenterViewport(mapWidth, mapHeight int) {
	pX, pY := gosoh.GetPlayerCoords()
	vw := float64(g.View.Width * gosoh.TileWidth)
	vh := float64(g.View.Height * gosoh.TileHeight)

	halfWidth := vw / 2
	halfHeight := vh / 2

	maxX := float64(mapWidth*gosoh.TileWidth) - vw
	maxY := float64(mapHeight*gosoh.TileHeight) - vh

	// Center on the player wherever possible...
	g.View.X = math.Min(math.Max(pX-halfWidth, 0), maxX)
	g.View.Y = math.Min(math.Max(pY-halfHeight, 0), maxY)

	// ...or if the map is smaller than the viewport, just center it
	if mapWidth < g.View.Width {
		g.View.X = halfWidth - float64(mapWidth*gosoh.TileWidth/2)
	}
	if mapHeight < g.View.Height {
		g.View.Y = halfHeight - float64(mapHeight*gosoh.TileHeight/2)
	}
}
