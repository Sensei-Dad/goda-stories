package gosoh

import (
	"fmt"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/**
TODO:
- Render the entire map layer as an Image, and store it to the MapArea?
	- This results in an image that might be too large? Need a way to load => tile zones as-before
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

	// Add the Player Entity
	// TODO: actually try to place the player on a movable tile
	playerX := 4
	playerY := 14
	pX := float64(playerX*TileWidth + (TileWidth / 2))
	pY := float64(playerY*TileHeight + (TileHeight / 2))

	ECSManager.NewEntity().
		AddComponent(playerComp, &PlayerInput{
			ShowDebug:    false,
			ShowBoxes:    false,
			ShowWalkable: false,
		}).
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
			X:     pX,
			Y:     pY,
			TileX: playerX,
			TileY: playerY,
		}).
		AddComponent(collideComp, &Collidable{
			IsBlocking: true,
			LeftEdge:   0.3,
			RightEdge:  0.3,
			TopEdge:    0.2,
			BottomEdge: 0.5,
		})

	players := ecs.BuildTag(playerComp, renderableComp, movementComp, creatureComp, positionComp)
	ECSTags["players"] = players
	playerView = ECSManager.CreateView(players)

	renderables := ecs.BuildTag(creatureComp, renderableComp, positionComp)
	ECSTags["renderables"] = renderables
	drawView = ECSManager.CreateView(renderables)

	creatures := ecs.BuildTag(creatureComp, positionComp)
	ECSTags["creatures"] = creatures

	movables := ecs.BuildTag(movementComp, collideComp, positionComp, creatureComp, renderableComp)
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
func NewOverworld(w, h int) *MapArea {
	// Make Dagobah
	dago := NewMapArea(2, 2)

	dago.AddZoneToArea(DAGOBAH_BL, 0, 1)
	dago.AddZoneToArea(DAGOBAH_TL, 0, 0)
	dago.AddZoneToArea(DAGOBAH_TR, 1, 0)
	dago.AddZoneToArea(DAGOBAH_BR, 1, 1)

	return &dago
}

func NewMapArea(w, h int) MapArea {
	ret := MapArea{
		Width:  w,
		Height: h,
	}
	tw := w * 18
	th := h * 18

	ret.Zones = make([][]*ZoneInfo, w)
	ret.Tiles = make([][]MapTile, tw)

	for x := 0; x < w; x++ {
		ret.Zones[x] = make([]*ZoneInfo, h)
	}
	for x := 0; x < tw; x++ {
		ret.Tiles[x] = make([]MapTile, th)
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

// Check for which tile (or tiles) a CollisionBox may be overlapping
// A box overlaps 4 tiles at most (as it's roughly the size of a tile)
func (a *MapArea) CheckCorners(b CollisionBox) (tl, tr, bl, br bool) {
	// Find the tile coordinates of each corner
	tx1 := int(b.X / float64(TileWidth))
	ty1 := int(b.Y / float64(TileHeight))

	tx2 := int((b.X + b.Width) / float64(TileWidth))
	ty2 := int((b.Y + b.Height) / float64(TileHeight))

	// If all 4 corners are within the same tile, then it doesn't overlap any other tiles
	if tx1 == tx2 && ty1 == ty2 {
		tl, tr, bl, br = false, false, false, false
		return
	}

	// Check the tile at each corner, including IsWalkable
	// If the box overlaps the edges of the map, that "counts"
	if tx1 != Clamp(tx1, 0, len(a.Tiles)) {
		tl = true
	} else {
		tl = b.OverlapsTile(a.Tiles[tx1][ty1])
	}
	if tx2 != Clamp(tx2, 0, len(a.Tiles)) {
		tr = true
	} else {
		tr = b.OverlapsTile(a.Tiles[tx2][ty1])
	}
	if ty1 != Clamp(ty1, 0, len(a.Tiles[0])) {
		bl = true
	} else {
		bl = b.OverlapsTile(a.Tiles[tx1][ty2])
	}
	if ty2 != Clamp(ty2, 0, len(a.Tiles[0])) {
		br = true
	} else {
		br = b.OverlapsTile(a.Tiles[tx2][ty2])
	}

	return
}

func (a *MapArea) AddZoneToArea(zoneId, x, y int) {
	zInfo := Zones[zoneId]
	// Save a ref to which zone number this is, so we can grab ZoneInfo later
	// Also used to determine when enemies get loaded / unloaded
	a.Zones[x][y] = &zInfo

	// Copy the zone tile IDs onto this area's tiles, and set the collision box position
	for j := 0; j < zInfo.Height; j++ {
		for i := 0; i < zInfo.Width; i++ {
			t := zInfo.GetTileAt(i, j)
			tx := (x * zInfo.Width) + i
			ty := (y * zInfo.Height) + j
			t.Box.X = float64(tx * TileWidth)
			t.Box.Y = float64(ty * TileHeight)
			a.Tiles[tx][ty] = t
		}
	}

	fmt.Printf("[AddZoneToArea] Added zone %03d to MapArea starting at (%d,%d)\n", zoneId, x*18, y*18)
}

// Pass in X,Y coords => get the Tile info at those coords
func (z *ZoneInfo) GetTileAt(x, y int) MapTile {
	tIndex := (z.Width * y) + x
	ret := MapTile{}
	ret.TerrainTileId = z.TileMaps.Terrain[tIndex]
	ret.WallTileId = z.TileMaps.Objects[tIndex]
	ret.OverlayTileId = z.TileMaps.Overlay[tIndex]
	ret.IsWalkable = CheckIsWalkable(ret.WallTileId)
	ret.Box = CollisionBox{
		Width:  float64(TileWidth),
		Height: float64(TileHeight),
	}

	return ret
}

// Draw a layer to the screen
func (a *MapArea) DrawLayer(lyr LayerName, screen *ebiten.Image, viewX, viewY, viewWidth, viewHeight, viewOffset float64) {
	viewBox := CollisionBox{
		X:      viewX,
		Y:      viewY,
		Width:  viewWidth,
		Height: viewHeight,
	}

	for y := 0; y < a.Height*18; y++ {
		for x := 0; x < a.Width*18; x++ {
			// Only need to draw a tile if we're inside the bounds of the MapArea
			if viewBox.Overlaps(a.Tiles[x][y].Box) {
				tNum := 65535 // Draw the blank tile by default
				switch lyr {
				case TerrainLayer:
					tNum = a.Tiles[x][y].TerrainTileId
				case WallsLayer:
					tNum = a.Tiles[x][y].WallTileId
				case OverlayLayer:
					tNum = a.Tiles[x][y].OverlayTileId
				}
				tile := GetTileImage(tNum)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(a.Tiles[x][y].Box.X-viewX+viewOffset, a.Tiles[x][y].Box.Y-viewY+viewOffset)
				// Draw the box, if it's collidable
				if !a.Tiles[x][y].IsWalkable {
					DrawTileBox(screen, a.Tiles[x][y].Box, viewBox.X, viewBox.Y, viewOffset)
				}
				screen.DrawImage(tile, op)
			}
		}
	}
}

func (a *MapArea) PrintMap() {
	fmt.Printf("Map of MapArea %d:\n", a.Id)
	for y := 0; y < len(a.Zones); y++ {
		line1 := ""
		for x := 0; x < len(a.Zones[y]); x++ {
			line1 += fmt.Sprintf("%03d  ", a.Zones[x][y].Id)
		}
		fmt.Println(line1)
	}
}
