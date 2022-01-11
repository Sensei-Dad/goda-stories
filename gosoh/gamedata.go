package gosoh

import (
	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// All the constants a person could ever want...
// Define "some" globals
var playerSpeed float64 = 2.0

const TileWidth, TileHeight int = 32, 32
const ViewportWidth, ViewportHeight int = 12, 10

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

// Special zone IDs to pay attention to
const (
	START_FACE int = 0
	WIN_FACE   int = 76
	LOSE_FACE  int = 77
	DAGOBAH_BL int = 93
	DAGOBAH_TL int = 94
	DAGOBAH_TR int = 95
	DAGOBAH_BR int = 96
)

func (d *CardinalDirection) IsDirection() bool {
	// Return true if movement is non-zero
	return !(d.DeltaX == 0 && d.DeltaY == 0)
}

// Game data
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

// A contiguous area to be displayed (Dagobah is 4 screens, the Overworld is 10x10, etc...)
type MapArea struct {
	Id     int
	Width  int
	Height int
	Zones  [][]*ZoneInfo
	Tiles  [][]MapTile
}

type MapTile struct {
	IsWalkable    bool
	TerrainTileId int
	WallTileId    int
	OverlayTileId int
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

// Collision detection for Creatures, Objects, etc.
type CollisionBox struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}
