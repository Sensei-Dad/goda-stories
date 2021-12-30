package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MapTile struct {
	PixelX     int
	PixelY     int
	IsWalkable bool
	Image      *ebiten.Image
}

type MapLayers struct {
	Terrain []MapTile
	Objects []MapTile
	Overlay []MapTile
}

func loadAllTiles(tiles []TileInfo) []*ebiten.Image {
	// return one big-ass slice of all tile images in the game
	ret := make([]*ebiten.Image, len(tiles))
	for x := 0; x < len(tiles); x++ {
		tFile := fmt.Sprintf("assets/tiles/tile_%04d.png", x)
		tImage, _, err := ebitenutil.NewImageFromFile(tFile)
		if err != nil {
			log.Fatal(err)
		}

		ret[x] = tImage
		fmt.Printf("      Loaded %s...\n", tFile)
	}
	return ret
}

func (g *Game) LoadMap(zone int) (ret MapLayers) {
	z := g.getZone(zone)
	ret.Terrain = make([]MapTile, len(z.LayerData.Terrain))
	ret.Objects = make([]MapTile, len(z.LayerData.Objects))
	ret.Overlay = make([]MapTile, len(z.LayerData.Overlay))

	for y := 0; y < z.Height; y++ {
		for x := 0; x < z.Width; x++ {
			// Assemble the map from layer data
			tNum := (y * z.Width) + x
			terNum := z.LayerData.Terrain[tNum]
			ter := MapTile{
				PixelX:     x * tileWidth,
				PixelY:     y * tileHeight,
				IsWalkable: true,
				Image:      g.getTile(terNum),
			}
			objNum := z.LayerData.Objects[tNum]
			obj := MapTile{
				PixelX:     x * tileWidth,
				PixelY:     y * tileHeight,
				IsWalkable: true,
				Image:      g.getTile(objNum),
			}
			ovrNum := z.LayerData.Overlay[tNum]
			ovr := MapTile{
				PixelX:     x * tileWidth,
				PixelY:     y * tileHeight,
				IsWalkable: true,
				Image:      g.getTile(ovrNum),
			}

			ret.Terrain[tNum] = ter
			ret.Objects[tNum] = obj
			ret.Overlay[tNum] = ovr
		}
	}

	return
}

func (g *Game) getTile(tNum int) *ebiten.Image {
	if tNum != 65535 {
		return g.Tiles[tNum]
	} else {
		// return a blank tile
		return ebiten.NewImage(tileWidth, tileHeight)
	}
}

func (g *Game) getCurrentZone() (z ZoneInfo) {
	return g.Zones[g.CurrentZone]
}

func (g *Game) getZone(zoneNum int) (z ZoneInfo) {
	return g.Zones[zoneNum]
}
