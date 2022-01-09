package main

import (
	"image"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/hajimehoshi/ebiten/v2"
)

// World holds all the various bits of the game that we generate
type GameWorld struct {
	Name     string
	SubAreas map[int]gosoh.MapArea
}

type ViewCoords struct {
	X           float64
	Y           float64
	Width       int
	Height      int
	CurrentArea int
}

func NewWorld() *GameWorld {
	gw := GameWorld{
		Name: "Goda Stories",
	}
	gw.SubAreas = make(map[int]gosoh.MapArea)

	// Make a new Overworld
	// Place the player on Dagobah
	world := gosoh.NewOverworld(10, 10)
	gw.SubAreas[world.Id] = world

	return &gw
}

// ALL this shizz needs to move to the ZoneManager
func (gw *GameWorld) DrawTerrain(screen *ebiten.Image, vp ViewCoords) {
	a := gw.SubAreas[vp.CurrentArea]
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(a.Terrain.SubImage(image.Rect(int(vp.X), int(vp.Y), int(vp.X+float64(vp.Width*gosoh.TileWidth)), int(vp.Y+float64(vp.Height*gosoh.TileHeight)))).(*ebiten.Image), op)
}

func (gw *GameWorld) DrawWalls(screen *ebiten.Image, vp ViewCoords) {
	a := gw.SubAreas[vp.CurrentArea]
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(a.Walls.SubImage(image.Rect(int(vp.X), int(vp.Y), int(vp.X+float64(vp.Width*gosoh.TileWidth)), int(vp.Y+float64(vp.Height*gosoh.TileHeight)))).(*ebiten.Image), op)
}

func (gw *GameWorld) DrawOverlay(screen *ebiten.Image, vp ViewCoords) {
	a := gw.SubAreas[vp.CurrentArea]
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(a.Overlay.SubImage(image.Rect(int(vp.X), int(vp.Y), int(vp.X+float64(vp.Width*gosoh.TileWidth)), int(vp.Y+float64(vp.Height*gosoh.TileHeight)))).(*ebiten.Image), op)
}
