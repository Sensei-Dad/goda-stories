package main

import (
	"math"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/hajimehoshi/ebiten/v2"
)

// World holds all the various bits of the game that we generate
type GameWorld struct {
	Name string
	Maps []gosoh.MapScreen
}

type ViewCoords struct {
	X          float64
	Y          float64
	Width      int
	Height     int
	CurrentMap int
}

func NewWorld() *GameWorld {
	// For now, choose a random screen from those available and start there
	// TODO: actual worldgen
	gw := GameWorld{
		Name: "Goda Stories",
	}
	zoneNum := gosoh.RandomInt(len(gosoh.Zones))
	gw.Maps = make([]gosoh.MapScreen, 0)
	ms := gosoh.NewMapScreen(zoneNum)

	gw.Maps = append(gw.Maps, ms)

	return &gw
}

// Render the Terrain layer, which goes on the bottom
func (gw *GameWorld) DrawTerrain(screen *ebiten.Image, vp ViewCoords) {
	ms := gw.Maps[vp.CurrentMap]
	// ms.PrintMap()
	for y := 0; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			// Make sure coords are within the map bounds
			mapX := int(math.Round(vp.X)) + x
			mapY := int(math.Round(vp.Y)) + y

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*gosoh.TileWidth), float64(y*gosoh.TileHeight))
			var img *ebiten.Image

			if ms.CoordsAreInBounds(mapX, mapY) {
				tile := ms.GetTileAt(mapX, mapY)
				img = gosoh.GetTileImage(tile.TerrainImage)
			} else {
				// If the tile isn't on the map, we don't need to do anything
				// fmt.Printf("  [GetMapCoords] V(%d,%d) => M(%d,%d) is not on the map\n", x, y, mapX, mapY)
				img = gosoh.BlankTile
			}

			screen.DrawImage(img, op)
		}
	}
}

func (gw *GameWorld) DrawObjects(screen *ebiten.Image, vp ViewCoords) {
	ms := gw.Maps[vp.CurrentMap]
	// ms.PrintMap()
	for y := 0; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			// Make sure coords are within the map bounds
			mapX := int(math.Round(vp.X)) + x
			mapY := int(math.Round(vp.Y)) + y

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*gosoh.TileWidth), float64(y*gosoh.TileHeight))
			var img *ebiten.Image

			if ms.CoordsAreInBounds(mapX, mapY) {
				tile := ms.GetTileAt(mapX, mapY)
				img = gosoh.GetTileImage(tile.ObjectsImage)
			} else {
				// If the tile isn't on the map, we don't need to do anything
				// fmt.Printf("  [GetMapCoords] V(%d,%d) => M(%d,%d) is not on the map\n", x, y, mapX, mapY)
				img = gosoh.BlankTile
			}

			screen.DrawImage(img, op)
		}
	}
}

func (gw *GameWorld) DrawOverlay(screen *ebiten.Image, vp ViewCoords) {
	ms := gw.Maps[vp.CurrentMap]
	// ms.PrintMap()
	for y := 0; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			// Make sure coords are within the map bounds
			mapX := int(math.Round(vp.X)) + x
			mapY := int(math.Round(vp.Y)) + y

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*gosoh.TileWidth), float64(y*gosoh.TileHeight))
			var img *ebiten.Image

			if ms.CoordsAreInBounds(mapX, mapY) {
				tile := ms.GetTileAt(mapX, mapY)
				img = gosoh.GetTileImage(tile.OverlayImage)
			} else {
				// If the tile isn't on the map, we don't need to do anything
				// fmt.Printf("  [GetMapCoords] V(%d,%d) => M(%d,%d) is not on the map\n", x, y, mapX, mapY)
				img = gosoh.BlankTile
			}

			screen.DrawImage(img, op)
		}
	}
}
