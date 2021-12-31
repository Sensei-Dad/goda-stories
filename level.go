package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MapScreen struct {
	Layers MapLayers
}

type MapTile struct {
	PixelX     int
	PixelY     int
	Id         int
	IsWalkable bool
	Image      *ebiten.Image
}

type MapLayers struct {
	Width   int
	Height  int
	MapId   int
	Terrain []MapTile
	Objects []MapTile
	Overlay []MapTile
}

func NewMapScreen(mapNum int) MapScreen {
	ms := MapScreen{}
	l := ms.LoadMap(mapNum)
	ms.Layers = l
	return ms
}

func LoadAllTiles(tiles []TileInfo) []*ebiten.Image {
	// return one big-ass slice of all tile images in the game,
	// because Go and I are lazy like that
	ret := make([]*ebiten.Image, len(tiles))
	for x := 0; x < len(tiles); x++ {
		tFile := fmt.Sprintf("assets/tiles/tile_%04d.png", x)
		tImage, _, err := ebitenutil.NewImageFromFile(tFile)
		if err != nil {
			log.Fatal(err)
		}

		ret[x] = tImage
	}
	return ret
}

func (ms *MapScreen) GetTileImage(tNum int) *ebiten.Image {
	if tNum != 65535 {
		return tiles[tNum]
	} else {
		// 65535 indicates a blank tile
		return blankTile
	}
}

func (ms *MapScreen) LoadMap(zone int) (ret MapLayers) {
	z := GetZone(zone)
	ret.Terrain = make([]MapTile, len(z.LayerData.Terrain))
	ret.Objects = make([]MapTile, len(z.LayerData.Objects))
	ret.Overlay = make([]MapTile, len(z.LayerData.Overlay))
	ret.Width = z.Width
	ret.Height = z.Height
	ret.MapId = zone

	for y := 0; y < z.Height; y++ {
		for x := 0; x < z.Width; x++ {
			// Assemble the map from layer data
			// TODO: account for entities (pushblocks, creatures, etc.)
			// TODO: probably lump the objImage, terImage, and ovrImage into the same MapTile struct
			tNum := (y * z.Width) + x
			terNum := z.LayerData.Terrain[tNum]
			ter := MapTile{
				Id:         terNum,
				PixelX:     x * tileWidth,
				PixelY:     y * tileHeight,
				IsWalkable: CheckIsWalkable(terNum),
				Image:      ms.GetTileImage(terNum),
			}
			objNum := z.LayerData.Objects[tNum]
			obj := MapTile{
				Id:         objNum,
				PixelX:     x * tileWidth,
				PixelY:     y * tileHeight,
				IsWalkable: CheckIsWalkable(objNum),
				Image:      ms.GetTileImage(objNum),
			}
			ovrNum := z.LayerData.Overlay[tNum]
			ovr := MapTile{
				Id:         ovrNum,
				PixelX:     x * tileWidth,
				PixelY:     y * tileHeight,
				IsWalkable: CheckIsWalkable(ovrNum),
				Image:      ms.GetTileImage(ovrNum),
			}

			ret.Terrain[tNum] = ter
			ret.Objects[tNum] = obj
			ret.Overlay[tNum] = ovr
		}
	}

	return
}

func (l *MapLayers) GetTileIndex(x, y int) int {
	// Helper function: input tile coords, return tile index
	return (y * l.Width) + x
}

func (l *MapLayers) DrawMap(screen *ebiten.Image) {
	for y := 0; y < ViewportWidth; y++ {
		for x := 0; x < ViewportHeight; x++ {
			tNum := l.GetTileIndex(x, y)
			terrainTile := l.Terrain[tNum]
			objectTile := l.Objects[tNum]
			overlayTile := l.Overlay[tNum]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(terrainTile.PixelX), float64(terrainTile.PixelY))
			screen.DrawImage(terrainTile.Image, op)
			screen.DrawImage(objectTile.Image, op)
			screen.DrawImage(overlayTile.Image, op)
		}
	}
}

func CheckIsWalkable(tNum int) bool {
	if tNum >= len(tileInfo) || tNum < 0 {
		return true
	}
	if tileInfo[tNum].IsWalkable {
		return true
	}
	return false
}
