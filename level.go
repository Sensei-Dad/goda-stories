package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MapScreen struct {
	Layers MapLayers
	Items  []ItemInfo
	// Creatures []CreatureInfo
	// PushBlocks []PushBlockInfo
}

type MapTile struct {
	IsWalkable   bool
	TerrainImage *ebiten.Image
	ObjectsImage *ebiten.Image
	OverlayImage *ebiten.Image
}

type MapLayers struct {
	Width  int
	Height int
	MapId  int
	Tiles  []MapTile
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
	z := zoneInfo[zone]
	ret.Tiles = make([]MapTile, len(z.LayerData.Terrain))
	ret.Width = z.Width
	ret.Height = z.Height
	ret.MapId = zone

	for y := 0; y < z.Height; y++ {
		for x := 0; x < z.Width; x++ {
			// Assemble the map from layer data
			// TODO: account for entities (pushblocks, creatures, etc.)
			// TODO: load triggers into ECS
			tIndex := (y * z.Width) + x
			terNum := z.LayerData.Terrain[tIndex]
			objNum := z.LayerData.Objects[tIndex]
			ovrNum := z.LayerData.Overlay[tIndex]

			// Check for objects
			for _, item := range itemInfo {
				// If it's an object, remove it from the map and save its coordinates
				if terNum == item.Id {
					terNum = 65535 // indicates a blank tile for the map rendering
					thing := item
					thing.MapX = x
					thing.MapY = y
					ms.Items = append(ms.Items, thing)
				}
				if objNum == item.Id {
					objNum = 65535
					thing := item
					thing.MapX = x
					thing.MapY = y
					ms.Items = append(ms.Items, thing)
				}
				if ovrNum == item.Id {
					ovrNum = 65535
					thing := item
					thing.MapX = x
					thing.MapY = y
					ms.Items = append(ms.Items, thing)
				}
			}

			tile := MapTile{
				IsWalkable:   CheckIsWalkable(objNum),
				TerrainImage: ms.GetTileImage(terNum),
				ObjectsImage: ms.GetTileImage(objNum),
				OverlayImage: ms.GetTileImage(ovrNum),
			}

			ret.Tiles[tIndex] = tile
		}
	}

	return
}

func (l *MapLayers) GetTileIndex(x, y int) int {
	// Helper function: input tile coords, return tile index
	if x != Clamp(x, 0, l.Width-1) || y != Clamp(y, 0, l.Height-1) {
		return -1
	}
	return (y * l.Width) + x
}

type LayerName string

const (
	Terrain LayerName = "terrain"
	Objects LayerName = "objects"
	Overlay LayerName = "overlay"
)

func (l *MapLayers) DrawLayer(layer LayerName, screen *ebiten.Image) {
	// Render the appropriate layer
	for y := 0; y < ViewportHeight; y++ {
		for x := 0; x < ViewportWidth; x++ {
			tIndex := l.GetTileIndex(x, y)
			var tile *ebiten.Image
			// Skip the tile if it's out of bounds
			if tIndex == -1 {
				tile = blankTile
			} else {
				switch layer {
				case Terrain:
					tile = l.Tiles[tIndex].TerrainImage
				case Objects:
					tile = l.Tiles[tIndex].ObjectsImage
				case Overlay:
					tile = l.Tiles[tIndex].OverlayImage
				default:
					log.Fatal("Unrecognized layer name")
				}
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileWidth), float64(y*tileHeight))
			screen.DrawImage(tile, op)
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

func (l *MapLayers) CanMove(x, y int) bool {
	// TODO: check collision on the next screen over
	tIndex := l.GetTileIndex(x, y)
	return (x != Clamp(x, 0, l.Width-1) || y != Clamp(y, 0, l.Height-1)) &&
		tIndex == Clamp(tIndex, 0, len(l.Tiles)-1) &&
		l.Tiles[tIndex].IsWalkable
}
