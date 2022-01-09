package gosoh

import (
	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type CreatureState string

const (
	Standing  CreatureState = "Standing"
	Walking   CreatureState = "Walking"
	InMotion  CreatureState = "InMotion"
	Attacking CreatureState = "Attacking"
	Dragging  CreatureState = "Dragging"
)

type CardinalDirection struct {
	Name   string
	DeltaX int
	DeltaY int
}

// Some globals
var ECSManager *ecs.Manager
var ECSTags map[string]ecs.Tag
var Tiles []*ebiten.Image
var TileInfos []TileInfo
var Zones []ZoneInfo
var Creatures []CreatureInfo
var Puzzles []PuzzleInfo
var Items []ItemInfo

var Up CardinalDirection = CardinalDirection{Name: "Up", DeltaX: 0, DeltaY: -1}
var Down CardinalDirection = CardinalDirection{Name: "Down", DeltaX: 0, DeltaY: 1}
var Left CardinalDirection = CardinalDirection{Name: "Left", DeltaX: -1, DeltaY: 0}
var Right CardinalDirection = CardinalDirection{Name: "Right", DeltaX: 1, DeltaY: 0}
var UpLeft CardinalDirection = CardinalDirection{Name: "UpLeft", DeltaX: -1, DeltaY: -1}
var DownLeft CardinalDirection = CardinalDirection{Name: "DownLeft", DeltaX: -1, DeltaY: 1}
var UpRight CardinalDirection = CardinalDirection{Name: "UpRight", DeltaX: 1, DeltaY: -1}
var DownRight CardinalDirection = CardinalDirection{Name: "DownRight", DeltaX: 1, DeltaY: 1}
var NoMove CardinalDirection = CardinalDirection{Name: "None", DeltaX: 0, DeltaY: 0}

var playerSpeed float64 = 2.0

const ViewportWidth, ViewportHeight int = 10, 10

// Distance from the viewport edge, in tiles, before the screen begins to scroll
const ViewportBuffer int = 4
const TileWidth, TileHeight int = 32, 32

func (d *CardinalDirection) NoDirection() bool {
	// Return true if movement is zero
	return (d.DeltaX == 0 && d.DeltaY == 0)
}

// Game data
type MapScreen struct {
	Width  int
	Height int
	ZoneId int
	Tiles  []MapTile
	Items  []ItemInfo
	// Creatures []CreatureInfo
	// PushBlocks []PushBlockInfo
}

type MapTile struct {
	IsWalkable   bool
	TerrainImage int
	ObjectsImage int
	OverlayImage int
}

type ZoneInfo struct {
	Id        int
	Type      string
	Width     int
	Height    int
	Flags     string
	LayerData struct {
		Terrain []int
		Objects []int
		Overlay []int
	}
	TileTriggers []TileTrigger
	Izax         []byte
	Izx2         []byte
	Izx3         []byte
	Izx4         []byte
	Iact         [][]byte
}

type TileInfo struct {
	// TODO: need to process flags in separate groups (TypeFlags, ItemFlags, etc...)?
	Id         int
	Flags      string
	Type       string
	IsWalkable bool
}

type TileTrigger struct {
	Type string
	X    int
	Y    int
	Arg  int
}

type PuzzleInfo struct {
	Id           int
	TextBytes    []byte
	LockItemId   int
	RewardItemId int
	RewardFlags  string
}

type ItemInfo struct {
	Id   int
	Name string
	MapX int
	MapY int
}

type CreatureInfo struct {
	Id     int
	Name   string
	Images map[CardinalDirection]int
}

// Kinda like CreatureInfo, but with everything you need
// to initialize the player entity
type PlayerInfo struct {
	Id     int
	Name   string
	Images map[CardinalDirection]int
	Speed  float64
}
