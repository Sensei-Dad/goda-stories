package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Tiles       []*ebiten.Image
	Zones       []ZoneInfo
	CurrentZone int
}

func NewGame(tiles []TileInfo, zones []ZoneInfo) *Game {
	g := &Game{}
	// for now, just load a static map
	g.Zones = zones
	g.Tiles = loadAllTiles(tiles)
	fmt.Printf("    Loaded %d tile images\n", len(g.Tiles))
	g.CurrentZone = 93
	return g
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each layer, starting at the bottom:
	//   = Overlay =
	//  == Objects ==
	// === Terrain ===
	m := g.LoadMap(g.CurrentZone)
	for y := 0; y < viewWidth; y++ {
		for x := 0; x < viewHeight; x++ {
			tNum := (y * g.getCurrentZone().Width) + x
			terrainTile := m.Terrain[tNum]
			objectTile := m.Objects[tNum]
			overlayTile := m.Overlay[tNum]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(terrainTile.PixelX), float64(terrainTile.PixelY))
			screen.DrawImage(terrainTile.Image, op)
			screen.DrawImage(objectTile.Image, op)
			screen.DrawImage(overlayTile.Image, op)
		}
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	// for now, return the map with nothing else around it
	return viewWidth * tileWidth, viewHeight * tileHeight
}
