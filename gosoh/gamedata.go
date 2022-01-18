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

func (d *CardinalDirection) IsHorizontal() bool {
	return d.DeltaX != 0
}

func (d *CardinalDirection) IsVertical() bool {
	return d.DeltaY != 0
}

func (d *CardinalDirection) IsDiagonal() bool {
	return d.DeltaX != 0 && d.DeltaY != 0
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

type LayerName string

const (
	TerrainLayer LayerName = "Terrain"
	WallsLayer   LayerName = "Walls"
	OverlayLayer LayerName = "Overlay"
)

type MapTile struct {
	IsWalkable    bool
	Box           CollisionBox
	TerrainTileId int
	WallTileId    int
	OverlayTileId int
}

type ZoneInfo struct {
	Id          int
	Biome       string
	Width       int
	Height      int
	Type        string
	IsOverworld bool
	TileMaps    struct {
		Terrain []int
		Objects []int
		Overlay []int
	}
	TileTriggers []TileTrigger
	ZoneActors   []ZoneActor
	RewardItems  []int // IZX2
	QuestNPCs    []int // IZX3
	Izx4a        int
	Izx4b        string
	Iact         [][]byte
}

type ZoneActor struct {
	CreatureId int
	ZoneX      int
	ZoneY      int
	Args       []byte
	Unknown    []byte
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

type ItemType int
type PuzzleText string

type PuzzleInfo struct {
	Id           int
	Type         string
	ItemType     string
	NeedText     string // "Hey, bring me a ___..."
	HaveText     string // "...and in return, I'll give you a ___..."
	DoneText     string // "...Thanks!" etc.
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

// Get a named item, by tile ID
func GetItemName(tNum int, iList []ItemInfo) (ret string) {
	ret = "UNKNOWN"
	for _, i := range iList {
		if i.Id == tNum {
			ret = i.Name
		}
	}

	return
}

func GetCreatureInfo(cNum int, cList []CreatureInfo) CreatureInfo {
	for _, c := range cList {
		if c.Id == cNum {
			return c
		}
	}
	return CreatureInfo{
		Id:   -1,
		Name: "UNKNOWN",
	}
}
