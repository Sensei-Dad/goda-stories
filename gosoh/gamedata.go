package gosoh

import (
	"fmt"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// All the constants a person could ever want...
// Define "some" globals
var playerSpeed float64 = 2.0

const TileWidth, TileHeight int = 32, 32
const TilesetColumns int = 20

var TilesetImage *ebiten.Image

var ECSManager *ecs.Manager
var ECSTags map[string]ecs.Tag
var TileInfos []TileInfo
var Zones []ZoneInfo
var Creatures []CreatureInfo
var Puzzles []PuzzleInfo
var Items []ItemInfo
var Sounds []string

var Up CardinalDirection = CardinalDirection{Name: "Up", DeltaX: 0, DeltaY: -1}
var Down CardinalDirection = CardinalDirection{Name: "Down", DeltaX: 0, DeltaY: 1}
var Left CardinalDirection = CardinalDirection{Name: "Left", DeltaX: -1, DeltaY: 0}
var Right CardinalDirection = CardinalDirection{Name: "Right", DeltaX: 1, DeltaY: 0}
var UpLeft CardinalDirection = CardinalDirection{Name: "UpLeft", DeltaX: -1, DeltaY: -1}
var DownLeft CardinalDirection = CardinalDirection{Name: "DownLeft", DeltaX: -1, DeltaY: 1}
var UpRight CardinalDirection = CardinalDirection{Name: "UpRight", DeltaX: 1, DeltaY: -1}
var DownRight CardinalDirection = CardinalDirection{Name: "DownRight", DeltaX: 1, DeltaY: 1}
var NoMove CardinalDirection = CardinalDirection{Name: "None", DeltaX: 0, DeltaY: 0}

var ClockwiseFrom = map[string]CardinalDirection{
	"Up":        UpRight,
	"UpRight":   Right,
	"Right":     DownRight,
	"DownRight": Down,
	"Down":      DownLeft,
	"DownLeft":  Left,
	"Left":      UpLeft,
	"UpLeft":    Up,
}

var WiddershinsFrom = map[string]CardinalDirection{
	"Up":        UpLeft,
	"UpRight":   Up,
	"Right":     UpRight,
	"DownRight": Right,
	"Down":      DownRight,
	"DownLeft":  Down,
	"Left":      DownLeft,
	"UpLeft":    Left,
}

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

type CreatureState string

const (
	Standing  CreatureState = "Standing"
	Walking   CreatureState = "Walking"
	Jumping   CreatureState = "Jumping"
	Attacking CreatureState = "Attacking"
	Dragging  CreatureState = "Dragging"
)

type CardinalDirection struct {
	Name   string
	DeltaX int
	DeltaY int
}

// A contiguous area to be displayed, i.e. a collection of Zones
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
		Walls   []int
		Overlay []int
	}
	Hotspots       []ZoneHotspot
	ActionTriggers []ActionTrigger
	ZoneActors     []ZoneActor
	RewardItems    []int // IZX2
	QuestNPCs      []int // IZX3
	Izx4a          int
	Izx4b          string
}

type ZoneActor struct {
	Index      int
	CreatureId int
	ZoneX      int
	ZoneY      int
	Args       []byte
	Unknown    []byte
}

// Tile triggers
type TriggerConditionType byte
type TriggerActionType byte
type TriggerHotspotType int

const (
	FirstEnter      TriggerConditionType = 0x00
	Enter           TriggerConditionType = 0x01
	BumpTile        TriggerConditionType = 0x02
	UseItem         TriggerConditionType = 0x03
	Walk            TriggerConditionType = 0x04
	TempVarEq       TriggerConditionType = 0x05
	RandVarEq       TriggerConditionType = 0x06
	RandVarGt       TriggerConditionType = 0x07
	RandVarLt       TriggerConditionType = 0x08
	EnterVehicle    TriggerConditionType = 0x09
	CheckTile       TriggerConditionType = 0x0A
	EnemyDead       TriggerConditionType = 0x0B
	AllEnemiesDead  TriggerConditionType = 0x0C
	HasItem         TriggerConditionType = 0x0D
	CheckQuestItem1 TriggerConditionType = 0x0E
	CheckQuestItem2 TriggerConditionType = 0x0F
	Unknown10       TriggerConditionType = 0x10
	GameInProgress  TriggerConditionType = 0x11
	GameCompleted   TriggerConditionType = 0x12
	HealthLt        TriggerConditionType = 0x13
	HealthGt        TriggerConditionType = 0x14
	Unknown15       TriggerConditionType = 0x15
	Unknown16       TriggerConditionType = 0x16
	UseWrongItem    TriggerConditionType = 0x17
	PlayerAtPos     TriggerConditionType = 0x18
	GlobalVarEq     TriggerConditionType = 0x19
	GlobalVarLt     TriggerConditionType = 0x1A
	GlobalVarGt     TriggerConditionType = 0x1B
	ExperienceEq    TriggerConditionType = 0x1C
	Unknown1D       TriggerConditionType = 0x1D
	Unknown1E       TriggerConditionType = 0x1E
	TempVarNe       TriggerConditionType = 0x1F
	RandVarNe       TriggerConditionType = 0x20
	GlobalVarNe     TriggerConditionType = 0x21
	CheckTileVar    TriggerConditionType = 0x22
	ExperienceGt    TriggerConditionType = 0x23
)

const (
	SetTile         TriggerActionType = 0x00
	ClearTile       TriggerActionType = 0x01
	MoveTile        TriggerActionType = 0x02
	DrawOverlayTile TriggerActionType = 0x03
	PlayerSay       TriggerActionType = 0x04
	CreatureSay     TriggerActionType = 0x05
	RedrawTile      TriggerActionType = 0x06
	RedrawRect      TriggerActionType = 0x07
	RenderChanges   TriggerActionType = 0x08
	WaitTicks       TriggerActionType = 0x09
	PlaySound       TriggerActionType = 0x0a
	FadeIn          TriggerActionType = 0x0b
	RandomNum       TriggerActionType = 0x0c
	SetTempVar      TriggerActionType = 0x0d
	AddTempVar      TriggerActionType = 0x0e
	SetTileVar      TriggerActionType = 0x0f
	ReleaseCamera   TriggerActionType = 0x10
	LockCamera      TriggerActionType = 0x11
	SetPlayerPos    TriggerActionType = 0x12
	MoveCamera      TriggerActionType = 0x13
	RunOnlyOnce     TriggerActionType = 0x14
	ShowObject      TriggerActionType = 0x15
	HideObject      TriggerActionType = 0x16
	ShowEntity      TriggerActionType = 0x17
	HideEntity      TriggerActionType = 0x18
	ShowAllEntities TriggerActionType = 0x19
	HideAllEntities TriggerActionType = 0x1a
	SpawnItem       TriggerActionType = 0x1b
	GiveToPlayer    TriggerActionType = 0x1c
	TakeFromPlayer  TriggerActionType = 0x1d
	OpenOrShow      TriggerActionType = 0x1e
	Unknown1f       TriggerActionType = 0x1f
	Unknown20       TriggerActionType = 0x20
	GoToZone        TriggerActionType = 0x21
	SetGlobalVar    TriggerActionType = 0x22
	AddGlobalVar    TriggerActionType = 0x23
	SetRandVar      TriggerActionType = 0x24
	AddToHealth     TriggerActionType = 0x25
)

const (
	TriggerSpot        TriggerHotspotType = 0
	SpawnLocation      TriggerHotspotType = 1
	ForceLocation      TriggerHotspotType = 2
	VehicleToSubarea   TriggerHotspotType = 3
	VehicleToOverworld TriggerHotspotType = 4
	LocatorSpot        TriggerHotspotType = 5
	ItemSpot           TriggerHotspotType = 6
	QuestNPCSpot       TriggerHotspotType = 7
	WeaponSpot         TriggerHotspotType = 8
	ZoneEntrance       TriggerHotspotType = 9
	ZoneExit           TriggerHotspotType = 10
	UNUSED             TriggerHotspotType = 11
	LockSpot           TriggerHotspotType = 12
	TeleportSpot       TriggerHotspotType = 13
	XWingFromDagobah   TriggerHotspotType = 14
	XWingToDagobah     TriggerHotspotType = 15
	UNKNOWNHOTSPOT     TriggerHotspotType = 16
)

type TriggerCondition struct {
	Condition TriggerConditionType
	Args      []int
}

type TriggerAction struct {
	Action TriggerActionType
	Args   []int
	Text   string
}

type ActionTrigger struct {
	Conditions []TriggerCondition
	Actions    []TriggerAction
}

type TileInfo struct {
	// TODO: need to process flags in separate groups (TypeFlags, ItemFlags, etc...)?
	Id         int
	Flags      string
	Type       string
	IsWalkable bool
}

type ZoneHotspot struct {
	Id   int
	Type TriggerHotspotType
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
func GetItemName(tNum int) (ret string) {
	ret = "UNKNOWN"
	for _, i := range Items {
		if i.Id == tNum {
			ret = i.Name
		}
	}

	return
}

func GetCreatureInfo(cNum int) CreatureInfo {
	for _, c := range Creatures {
		if c.Id == cNum {
			return c
		}
	}
	return CreatureInfo{
		Id:   -1,
		Name: "UNKNOWN",
	}
}

func (a *ActionTrigger) ToString() string {
	ret := ""
	for i, c := range a.Conditions {
		if i == 0 {
			ret += "When "
		} else {
			ret += "\n and "
		}
		ret += c.ToString()
	}
	ret += "...\n"
	for _, a := range a.Actions {
		ret += "   - " + a.ToString() + "\n"
	}
	return ret
}

func (t *TriggerCondition) ToString() string {
	ret := ""
	switch t.Condition {
	case FirstEnter:
		ret = "FirstEnter"
	case Enter:
		ret = "ZoneEnter"
	case BumpTile:
		ret = "BumpTile"
	case UseItem:
		ret = "UseItem"
	case Walk:
		ret = "TileWalk"
	case TempVarEq:
		ret = "TVar_EQ"
	case RandVarEq:
		ret = "RVar_EQ"
	case RandVarGt:
		ret = "RVar_GT"
	case RandVarLt:
		ret = "RVar_LT"
	case EnterVehicle:
		ret = "EnterVehicle"
	case CheckTile:
		ret = "CheckTile"
	case EnemyDead:
		ret = "CrtrDead"
	case AllEnemiesDead:
		ret = "AllDead"
	case HasItem:
		ret = "HasItem"
	case CheckQuestItem1:
		ret = "Item1Is"
	case CheckQuestItem2:
		ret = "Item2Is"
	case Unknown10:
		ret = "Unkwn10"
	case GameInProgress:
		ret = "MainQuestOpen"
	case GameCompleted:
		ret = "MainQuestDone"
	case HealthLt:
		ret = "Life_LT"
	case HealthGt:
		ret = "Life_GT"
	case Unknown15:
		ret = "Unkwn15"
	case Unknown16:
		ret = "Unkwn16"
	case UseWrongItem:
		ret = "WrongItem"
	case PlayerAtPos:
		ret = "PlyrAtPos"
	case GlobalVarEq:
		ret = "GVar_EQ"
	case GlobalVarLt:
		ret = "GVar_LT"
	case GlobalVarGt:
		ret = "GVar_GT"
	case ExperienceEq:
		ret = "Wins_EQ"
	case Unknown1D:
		ret = "Unkwn1d"
	case Unknown1E:
		ret = "Unkwn1e"
	case TempVarNe:
		ret = "TVar_NE"
	case RandVarNe:
		ret = "RVar_NE"
	case GlobalVarNe:
		ret = "GVar_NE"
	case CheckTileVar:
		ret = "CheckTileVar"
	case ExperienceGt:
		ret = "Wins_GT"
	}
	for _, arg := range t.Args {
		ret += fmt.Sprintf(",%d", arg)
	}
	return ret
}

func (a *TriggerAction) ToString() string {
	ret := ""
	switch a.Action {
	case SetTile:
		ret = "SetTile"
	case ClearTile:
		ret = "ClearTile"
	case MoveTile:
		ret = "MoveTile"
	case DrawOverlayTile:
		ret = "DrawOver"
	case PlayerSay:
		ret = "PlyrSez"
	case CreatureSay:
		ret = "CrtrSez"
	case RedrawTile:
		ret = "DrawTile"
	case RedrawRect:
		ret = "DrawRect"
	case RenderChanges:
		ret = "DrawAll"
	case WaitTicks:
		ret = "WaitFor"
	case PlaySound:
		ret = "PlaySound"
	case FadeIn:
		ret = "FadeIn"
	case RandomNum:
		ret = "RVarRange"
	case SetTempVar:
		ret = "SetTVar"
	case AddTempVar:
		ret = "AddTVar"
	case SetTileVar:
		ret = "SetTileVar"
	case ReleaseCamera:
		ret = "FreeCam"
	case LockCamera:
		ret = "LockCam"
	case SetPlayerPos:
		ret = "SetPlyrPos"
	case MoveCamera:
		ret = "MoveCam"
	case RunOnlyOnce:
		ret = "RunOnce"
	case ShowObject:
		ret = "ShowObj"
	case HideObject:
		ret = "HideObj"
	case ShowEntity:
		ret = "ShowCrtr"
	case HideEntity:
		ret = "HideCrtr"
	case ShowAllEntities:
		ret = "ShowAll"
	case HideAllEntities:
		ret = "HideAll"
	case SpawnItem:
		ret = "SpawnItem"
	case GiveToPlayer:
		ret = "GiveItem"
	case TakeFromPlayer:
		ret = "TakeItem"
	case OpenOrShow:
		ret = "OpenOrShow"
	case Unknown1f:
		ret = "Unkwn1f"
	case Unknown20:
		ret = "Unkwn20"
	case GoToZone:
		ret = "GoToZone"
	case SetGlobalVar:
		ret = "SetGVar"
	case AddGlobalVar:
		ret = "AddGVar"
	case SetRandVar:
		ret = "SetRVar"
	case AddToHealth:
		ret = "AddLife"
	}
	for _, arg := range a.Args {
		ret += fmt.Sprintf(",%d", arg)
	}
	ret += "," + a.Text
	return ret
}

func (hs *ZoneHotspot) ToString() string {
	// ret := fmt.Sprintf("%02d (%d, %d) ", hs.Id, hs.X, hs.Y)
	ret := ""
	switch hs.Type {
	case ZoneEntrance:
		ret += "EnterZone"
	case ZoneExit:
		ret += "ExitZone"
	case VehicleToSubarea, VehicleToOverworld:
		ret += "VehicleSpot"
	case XWingToDagobah, XWingFromDagobah:
		ret += "XwingSpot"
	case TriggerSpot:
		ret += "TriggerSpot"
	case SpawnLocation:
		ret += "NpcSpawnSpot"
	case ForceLocation:
		ret += "ForceSpot"
	case LocatorSpot:
		ret += "GetLocator"
	case ItemSpot:
		ret += "ItemSpot"
	case QuestNPCSpot:
		ret += "QuestNPC"
	case WeaponSpot:
		ret += "WeaponSpot"
	case LockSpot:
		ret += "LockSpot"
	case TeleportSpot:
		ret += "TPortSpot"
	case UNUSED:
		ret += "UNUSED"
	case UNKNOWNHOTSPOT:
		ret += "UNKNOWN"
	default:
		ret += fmt.Sprintf("Unhandled trigger type, arg %d", hs.Arg)
		log.Fatal("UNHANDLED HOTSPOT TYPE")
	}
	ret += fmt.Sprintf(",%d", hs.Arg)

	return ret
}
