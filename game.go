package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const tileWidth, tileHeight = 32, 32         // Tile width and height, in pixels
const ViewportWidth, ViewportHeight = 18, 18 // Viewport width and height, in tiles
// const uiPadding = 5                  // Padding between UI elements, in pixels

type Game struct {
	CurrentScreen int
	World         GameWorld
}

var blankTile = ebiten.NewImage(tileWidth, tileHeight)
var tiles = []*ebiten.Image{}

func NewGame(tInfo []TileInfo, zones []ZoneInfo) *Game {
	tiles = LoadAllTiles(tInfo)
	g := &Game{}
	g.World = NewWorld()

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
	layers := g.World.GetLayers()
	for y := 0; y < ViewportWidth; y++ {
		for x := 0; x < ViewportHeight; x++ {
			tNum := layers.GetTileNum(x, y)
			terrainTile := layers.Terrain[tNum]
			objectTile := layers.Objects[tNum]
			overlayTile := layers.Overlay[tNum]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(terrainTile.PixelX), float64(terrainTile.PixelY))
			screen.DrawImage(terrainTile.Image, op)
			screen.DrawImage(objectTile.Image, op)
			screen.DrawImage(overlayTile.Image, op)
		}
	}

	// Show FPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(w, h int) (int, int) {
	// for now, return the map with nothing else around it
	return ViewportWidth * tileWidth, ViewportHeight * tileHeight
}

func GetZone(zoneNum int) (z ZoneInfo) {
	return zoneInfo[zoneNum]
}
