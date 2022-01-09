package gosoh

import (
	"fmt"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var BlankTile *ebiten.Image = ebiten.NewImage(32, 32)

func InitializeECS() {
	// Initialize the world via the ECS
	ECSTags = make(map[string]ecs.Tag)
	ECSManager = ecs.NewManager()

	// Make the global components
	playerComp = ECSManager.NewComponent()
	creatureComp = ECSManager.NewComponent()
	renderableComp = ECSManager.NewComponent()
	movementComp = ECSManager.NewComponent()
	positionComp = ECSManager.NewComponent()
	collideComp = ECSManager.NewComponent()

	// TODO: actually try to place the player on a movable tile
	playerX := 4
	playerY := 7

	ECSManager.NewEntity().
		AddComponent(playerComp, &PlayerInput{}).
		AddComponent(creatureComp, &Creature{
			Name:   Creatures[0].Name,
			State:  Standing,
			Facing: Down,
			// TODO: Move this flag to the Movement component
			InMotion:   false,
			CreatureId: 0,
		}).
		AddComponent(renderableComp, &Renderable{
			Image:  Tiles[Creatures[0].Images[Down]], // TODO: test the walking animations
			PixelX: float64(playerX * TileWidth),
			PixelY: float64(playerY * TileHeight),
		}).
		AddComponent(movementComp, &Movable{
			OldX:      4,
			OldY:      7,
			Speed:     playerSpeed,
			Direction: NoMove,
		}).
		AddComponent(positionComp, &Position{
			X: 4,
			Y: 7,
		}).AddComponent(collideComp, &Collidable{
		IsBlocking: true,
	})

	players := ecs.BuildTag(playerComp, movementComp, creatureComp, positionComp)
	ECSTags["players"] = players
	playerView = ECSManager.CreateView(players)

	renderables := ecs.BuildTag(creatureComp, renderableComp)
	ECSTags["renderables"] = renderables
	drawView = ECSManager.CreateView(renderables)

	creatures := ecs.BuildTag(creatureComp, positionComp)
	ECSTags["creatures"] = creatures

	movables := ecs.BuildTag(movementComp, positionComp, creatureComp, renderableComp)
	ECSTags["movables"] = movables
	moveView = ECSManager.CreateView(movables)

	collidables := ecs.BuildTag(collideComp, positionComp)
	ECSTags["collidables"] = collidables
	collideView = ECSManager.CreateView(collidables)
}

func LoadAllTiles(tiles []TileInfo) {
	// return one big-ass slice of all tile images in the game,
	// because Go and I are lazy like that
	Tiles = make([]*ebiten.Image, len(tiles))
	TileInfos = tiles
	for x := 0; x < len(tiles); x++ {
		tFile := fmt.Sprintf("assets/tiles/tile_%04d.png", x)
		tImage, _, err := ebitenutil.NewImageFromFile(tFile)
		if err != nil {
			log.Fatal(err)
		}

		Tiles[x] = tImage
	}
}

func NewMapScreen(zone int) MapScreen {
	z := Zones[zone]
	fmt.Printf("Pulling zone data for zone %03d\n", zone)
	ret := MapScreen{
		Width:  z.Width,
		Height: z.Height,
		ZoneId: zone,
	}
	ret.Tiles = make([]MapTile, 0)

	for y := 0; y < z.Height; y++ {
		for x := 0; x < z.Width; x++ {
			// Assemble the map from layer data
			// TODO: account for entities (pushblocks, creatures, etc.)
			// TODO: load triggers into ECS
			tIndex := (y * z.Width) + x

			tile := MapTile{}

			tile.TerrainImage = z.LayerData.Terrain[tIndex]
			tile.ObjectsImage = z.LayerData.Objects[tIndex]
			tile.OverlayImage = z.LayerData.Overlay[tIndex]

			tile.IsWalkable = CheckIsWalkable(z.LayerData.Objects[tIndex])

			ret.Tiles = append(ret.Tiles, tile)
		}
	}

	return ret
}

func GetTileImage(tNum int) *ebiten.Image {
	if tNum != 65535 {
		return Tiles[tNum]
	} else {
		// 65535 indicates a blank tile
		return BlankTile
	}
}

func CheckIsWalkable(tNum int) bool {
	if tNum >= len(Tiles) || tNum < 0 {
		return true
	}
	return TileInfos[tNum].IsWalkable
}

func (ms *MapScreen) CoordsAreInBounds(x, y int) bool {
	if x == Clamp(x, 0, ms.Width-1) && y == Clamp(y, 0, ms.Height-1) {
		return true
	}
	return false
}

// Pass in X,Y coords => get the Tile Id at those coords
func (ms *MapScreen) GetTileAt(x, y int) MapTile {
	tId := (y * ms.Width) + x
	return ms.Tiles[tId]
}

func (ms *MapScreen) PrintMap() {
	fmt.Printf("Map of zone_%03d:\n", ms.ZoneId)
	for y := 0; y < ms.Height; y++ {
		line1 := ""
		line2 := ""
		line3 := ""
		for x := 0; x < ms.Width; x++ {
			tile := ms.GetTileAt(x, y)
			if tile.TerrainImage == Clamp(tile.TerrainImage, 0, len(Tiles)-1) {
				line1 += fmt.Sprintf("%04d  ", tile.TerrainImage)
			} else {
				line1 += "      "
			}
			if tile.ObjectsImage == Clamp(tile.ObjectsImage, 0, len(Tiles)-1) {
				line2 += fmt.Sprintf("%04d  ", tile.ObjectsImage)
			} else {
				line2 += "      "
			}
			if tile.OverlayImage == Clamp(tile.OverlayImage, 0, len(Tiles)-1) {
				line3 += fmt.Sprintf("%04d  ", tile.OverlayImage)
			} else {
				line3 += "      "
			}
		}
		fmt.Println(line1)
		fmt.Println(line2)
		fmt.Println(line3)
		fmt.Println()
	}
}
