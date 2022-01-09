package gosoh

import (
	"fmt"
	"image/color"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/**
TODO:
- Render the entire map layer as an Image, and store it to the MapArea
- New struct: MapArea
	- Dagobah - make Dagobah (z93-96) to test rendering / loading multiple Zones, Yoda's hut for entrances, etc.
	Overworld - the big 10x10 map
	Subareas - inside buildings, etc.
	- LocatorImage bool - for subareas, all we need to do is point to the parent area's Locator map
- New struct: LocatorMap (see above)
- Refactor MapArea / Worldgen:
	- generate multiple screens and tile them together into a single *image ahead of time
	- SubZones for door entrances, etc
	- load the SubZone, stash the Entities
	- when you exit, the ECS stashes / restores them
	- for smaller maps, just shove some black around the edges until it's Viewport-sized
	- one big Image each for Terrain / Objects / Overlay
	- Add functions / getters / Neighbors, so they act like Nodes
	- save Objects/Triggers to ECS (LoadZoneObjects)
	- basically EVERY pushblock, chest, NPC, alternate NPC, enemy, etc. will be counted and added to the Entity pool
	- new "Processible" component will limit which Entities are tracked at any time
	- rename Objects to Walls layer => Obj to ECS, everything else is a Wall
- Refactor Draw:
	- instead of tile-by-tile, just draw a Subimage of each layer and use the Viewport to make a Rect
- Refactor Actions:
	- do movement, collisions, etc. by pixel / bounding boxes, filter with the Processible comp, and ignore the Viewport entirely
	- in another manager, go through each frame and activate / "de-spawn" each area's entities ahead of / behind the player as they move
	- Keep a 3x3 square of "maps" active, kinda like how chunks in Minecraft get loaded (but WAY smaller scale; we want this to be efficient)
	- Can also have push-blocks reset when de-spawned, instead of keeping their coords
**/

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
			Name:       Creatures[0].Name,
			State:      Standing,
			Facing:     Down,
			CanMove:    true,
			CreatureId: 0,
		}).
		AddComponent(renderableComp, &Renderable{
			Image: Tiles[Creatures[0].Images[Down]], // TODO: test the walking animations
		}).
		AddComponent(movementComp, &Movable{
			Speed:     playerSpeed,
			Direction: NoMove,
		}).
		AddComponent(positionComp, &Position{
			X:     float64(playerX*TileWidth) + 0.5, // Start center-tile
			Y:     float64(playerY*TileHeight) + 0.5,
			TileX: playerX,
			TileY: playerY,
		}).
		AddComponent(collideComp, &Collidable{
			IsBlocking: true,
		})

	players := ecs.BuildTag(playerComp, renderableComp, movementComp, creatureComp, positionComp)
	ECSTags["players"] = players
	playerView = ECSManager.CreateView(players)

	renderables := ecs.BuildTag(creatureComp, renderableComp, positionComp)
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

// Worldgen!
func NewOverworld(w, h int) MapArea {
	// Make Dagobah
	dago := NewMapArea(2, 2)

	dago.AddZoneToArea(DAGOBAH_BL, 0, 1)
	dago.AddZoneToArea(DAGOBAH_TL, 0, 0)
	dago.AddZoneToArea(DAGOBAH_TR, 1, 0)
	dago.AddZoneToArea(DAGOBAH_BR, 1, 1)

	return dago
}

func NewMapArea(w, h int) MapArea {
	ret := MapArea{
		Width:  w,
		Height: h,
	}
	tw := w * TileWidth * 18
	th := h * TileHeight * 18
	ret.Terrain = ebiten.NewImage(tw, th)
	ret.Walls = ebiten.NewImage(tw, th)
	ret.Overlay = ebiten.NewImage(tw, th)

	// Black background
	ret.Terrain.Fill(color.Black)

	ret.Zones = make([]int, w*h)

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

func (ms *MapArea) CoordsAreInBounds(x, y float64) bool {
	if x >= 0 && x <= float64((ViewportWidth-1)*TileWidth) {
		if y >= 0 && y <= float64((ViewportHeight-1)*TileHeight) {
			return true
		}
	}
	return false
}

func (a *MapArea) AddZoneToArea(zoneId, x, y int) {
	zInfo := Zones[zoneId]
	zIndex := (a.Width * y) + x
	areaX := x * 18 * TileWidth // pixel coords of the top-left corner we start in
	areaY := y * 18 * TileHeight

	// Draw the zone tiles onto this area's map images
	for j := 0; j < zInfo.Height; j++ {
		for i := 0; i < zInfo.Width; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(areaX+(i*TileWidth)), float64(areaY+(j*TileHeight)))

			tIndex := (zInfo.Width * j) + i
			terImg := GetTileImage(zInfo.LayerData.Terrain[tIndex])
			objImg := GetTileImage(zInfo.LayerData.Objects[tIndex])
			ovrImg := GetTileImage(zInfo.LayerData.Overlay[tIndex])

			a.Terrain.DrawImage(terImg, op)
			a.Walls.DrawImage(objImg, op)
			a.Overlay.DrawImage(ovrImg, op)
		}
	}

	// Save a ref to which zone number this is, so we can grab ZoneInfo later
	a.Zones[zIndex] = zoneId
}

// Pass in X,Y pixel coords => get the Tile info at those coords
// func (a *MapArea) GetTileAt(x, y float64) MapTile {
// 	// Find where this Zone is, within the Area...
// 	ax := int(x/float64(TileWidth)) / 18
// 	ay := int(y/float64(TileHeight)) / 18
// 	zIndex := (a.Width * ay) + ax
// 	zone := Zones[a.Zones[zIndex]]

// 	// ...then find where this tile is within the Zone...
// 	zx := int(x/float64(TileWidth)) - (18 * ax)
// 	zy := int(y/float64(TileHeight)) - (18 * ay)
// 	tIndex := (zone.Width * zy) + zx

// 	return TileInfos[tIndex]
// 	return error
// }

// func (ms *MapArea) PrintMap() {
// 	fmt.Printf("Map of zone_%03d:\n", ms.ZoneId)
// 	for y := 0; y < ms.Height; y++ {
// 		line1 := ""
// 		line2 := ""
// 		line3 := ""
// 		for x := 0; x < ms.Width; x++ {
// 			tile := ms.GetTileAt(x, y)
// 			if tile.TerrainImage == Clamp(tile.TerrainImage, 0, len(Tiles)-1) {
// 				line1 += fmt.Sprintf("%04d  ", tile.TerrainImage)
// 			} else {
// 				line1 += "      "
// 			}
// 			if tile.ObjectsImage == Clamp(tile.ObjectsImage, 0, len(Tiles)-1) {
// 				line2 += fmt.Sprintf("%04d  ", tile.ObjectsImage)
// 			} else {
// 				line2 += "      "
// 			}
// 			if tile.OverlayImage == Clamp(tile.OverlayImage, 0, len(Tiles)-1) {
// 				line3 += fmt.Sprintf("%04d  ", tile.OverlayImage)
// 			} else {
// 				line3 += "      "
// 			}
// 		}
// 		fmt.Println(line1)
// 		fmt.Println(line2)
// 		fmt.Println(line3)
// 		fmt.Println()
// 	}
// }
