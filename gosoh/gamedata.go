package gosoh

import (
	"fmt"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// All the constants a person could ever want...
// Define "some" globals
var playerSpeed float64 = 2.0

const TileWidth, TileHeight int = 32, 32

var ECSManager *ecs.Manager
var ECSTags map[string]ecs.Tag
var Tiles []*ebiten.Image
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
	TileTriggers   []TileTrigger
	ActionTriggers []ActionTrigger
	ZoneActors     []ZoneActor
	RewardItems    []int // IZX2
	QuestNPCs      []int // IZX3
	Izx4a          int
	Izx4b          string
}

type ZoneActor struct {
	CreatureId int
	ZoneX      int
	ZoneY      int
	Args       []byte
	Unknown    []byte
}

// Tile triggers
type TriggerConditionType byte
type TriggerActionType byte

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
		ret = "the Player enters this Zone for the first time"
	case Enter:
		ret = "the Player enters this Zone"
	case BumpTile:
		ret = fmt.Sprintf("the Player interacts with tile_%04d at (%d, %d)", t.Args[2], t.Args[0], t.Args[1])
	case UseItem:
		ret = fmt.Sprintf("the Player uses item_%03d on tile_%04d at (%d, %d, %d)", t.Args[4], t.Args[3], t.Args[0], t.Args[1], t.Args[2])
	case Walk:
		ret = fmt.Sprintf("the Player walks onto tile (%d, %d)", t.Args[0], t.Args[1])
	case TempVarEq:
		ret = fmt.Sprintf("TempVar is equal to %d", t.Args[0])
	case RandVarEq:
		ret = fmt.Sprintf("RandVar is equal to %d", t.Args[0])
	case RandVarGt:
		ret = fmt.Sprintf("RandVar is greater than %d", t.Args[0])
	case RandVarLt:
		ret = fmt.Sprintf("RandVar is less than %d", t.Args[0])
	case EnterVehicle:
		ret = "the Player enters a vehicle"
	case CheckTile:
		tnam := fmt.Sprintf("tile_%04d", t.Args[0])
		if t.Args[0] == 65535 {
			tnam = "empty"
		}
		ret = fmt.Sprintf("the tile at (%d, %d, %d) is %s", t.Args[1], t.Args[2], t.Args[3], tnam)
	case EnemyDead:
		ret = fmt.Sprintf("the enemy #%02d is dead", t.Args[0])
	case AllEnemiesDead:
		ret = "all enemies are dead"
	case HasItem:
		ret = fmt.Sprintf("the Player has item_%03d", t.Args[0])
	case CheckQuestItem1:
		ret = fmt.Sprintf("the first Quest item is item_%03d", t.Args[0])
	case CheckQuestItem2:
		ret = fmt.Sprintf("the second Quest item is item_%03d", t.Args[0])
	case Unknown10:
		ret = "Unknown (10)..."
	case GameInProgress:
		ret = "the Player has not finished the game"
	case GameCompleted:
		ret = "the Player has finished the game"
	case HealthLt:
		ret = fmt.Sprintf("the Player has less than %d health", t.Args[0])
	case HealthGt:
		ret = fmt.Sprintf("the Player has more than %d health", t.Args[0])
	case Unknown15:
		ret = "Unknown (15)..."
	case Unknown16:
		ret = "Unknown (16)..."
	case UseWrongItem:
		ret = fmt.Sprintf("the Player uses the wrong item on tile_%04d at (%d, %d, %d)", t.Args[3], t.Args[0], t.Args[1], t.Args[2])
	case PlayerAtPos:
		ret = fmt.Sprintf("the Player is at zone coords (%d, %d)", t.Args[0], t.Args[1])
	case GlobalVarEq:
		ret = fmt.Sprintf("GlobalVar is equal to %d", t.Args[0])
	case GlobalVarLt:
		ret = fmt.Sprintf("GlobalVar is less than %d", t.Args[0])
	case GlobalVarGt:
		ret = fmt.Sprintf("GlobalVar is greater than %d", t.Args[0])
	case ExperienceEq:
		ret = fmt.Sprintf("the Player's XP is equal to %d", t.Args[0])
	case Unknown1D:
		ret = "Unknown (1D)..."
	case Unknown1E:
		ret = "Unknown (1E)..."
	case TempVarNe:
		ret = fmt.Sprintf("TempVar is not equal to %d", t.Args[0])
	case RandVarNe:
		ret = fmt.Sprintf("RandVar is not equal to %d", t.Args[0])
	case GlobalVarNe:
		ret = fmt.Sprintf("GlobalVar is not equal to %d", t.Args[0])
	case CheckTileVar:
		ret = fmt.Sprintf("the TileVar stored at (%d, %d, %d) is %d", t.Args[1], t.Args[2], t.Args[3], t.Args[0])
	case ExperienceGt:
		ret = fmt.Sprintf("the Player's XP is greater than %d", t.Args[0])
	}
	return ret
}

func (a *TriggerAction) ToString() string {
	ret := ""
	switch a.Action {
	case SetTile:
		ret = fmt.Sprintf("Set the tile at (%d, %d, %d) to %d", a.Args[0], a.Args[1], a.Args[2], a.Args[3])
	case ClearTile:
		ret = fmt.Sprintf("Clear the tile at (%d, %d, %d)", a.Args[0], a.Args[1], a.Args[2])
	case MoveTile:
		ret = fmt.Sprintf("Move the tile at (%d, %d, %d) to (%d, %d, %d)", a.Args[0], a.Args[1], a.Args[2], a.Args[3], a.Args[4], a.Args[2])
	case DrawOverlayTile:
		ret = fmt.Sprintf("Draw tile_%04d over the top of (%d, %d)", a.Args[2], a.Args[0], a.Args[1])
	case PlayerSay:
		ret = fmt.Sprintf("The player says: \"%s\"", a.Text)
	case CreatureSay:
		ret = fmt.Sprintf("The creature at (%d, %d) says: \"%s\"", a.Args[0], a.Args[1], a.Text)
	case RedrawTile:
		ret = fmt.Sprintf("Redraw the tile at (%d, %d)", a.Args[0], a.Args[1])
	case RedrawRect:
		ret = fmt.Sprintf("Redraw all the tiles from (%d, %d) to (%d, %d)", a.Args[0], a.Args[1], a.Args[2], a.Args[3])
	case RenderChanges:
		ret = "Redraw all the things"
	case WaitTicks:
		ret = fmt.Sprintf("Wait for %d ticks", a.Args[0])
	case PlaySound:
		ret = fmt.Sprintf("Play \"%s\"", Sounds[a.Args[0]])
	case FadeIn:
		ret = "Do the \"Screen-Wipe In\" animation"
	case RandomNum:
		ret = fmt.Sprintf("Set RandVar to a random value between 0 and %d", a.Args[0])
	case SetTempVar:
		ret = fmt.Sprintf("Set TempVar to %d", a.Args[0])
	case AddTempVar:
		ret = fmt.Sprintf("Add %d to TempVar", a.Args[0])
	case SetTileVar:
		ret = fmt.Sprintf("Set the TileVar at (%d, %d, %d) to %d", a.Args[0], a.Args[1], a.Args[2], a.Args[3])
	case ReleaseCamera:
		ret = "Un-anchor the camera from the Player"
	case LockCamera:
		ret = "Anchor the camera to the Player's position"
	case SetPlayerPos:
		ret = fmt.Sprintf("Teleport the Player to zone coords (%d, %d)", a.Args[0], a.Args[1])
	case MoveCamera:
		ret = fmt.Sprintf("Pan the camera from (%d, %d) to (%d, %d) over the next %d ticks", a.Args[0], a.Args[1], a.Args[2], a.Args[3], a.Args[4])
	case RunOnlyOnce:
		ret = "Destroy this trigger after it finishes"
	case ShowObject:
		ret = fmt.Sprintf("Show object #%d", a.Args[0])
	case HideObject:
		ret = fmt.Sprintf("Hide object #%d", a.Args[0])
	case ShowEntity:
		ret = fmt.Sprintf("Show creature #%d", a.Args[0])
	case HideEntity:
		ret = fmt.Sprintf("Hide creature #%d", a.Args[0])
	case ShowAllEntities:
		ret = "Show all creatures"
	case HideAllEntities:
		ret = "Hide all creatures"
	case SpawnItem:
		ret = fmt.Sprintf("Spawn item_%03d at (%d, %d)", a.Args[0], a.Args[1], a.Args[2])
	case GiveToPlayer:
		ret = fmt.Sprintf("Give item_%03d to the Player", a.Args[0])
	case TakeFromPlayer:
		ret = fmt.Sprintf("Take item_%03d from the Player", a.Args[0])
	case OpenOrShow:
		ret = "Set a bunch of values to 1...?"
	case Unknown1f:
		ret = "Unknown (1F)..."
	case Unknown20:
		ret = "Unknown (20)... (never used?)"
	case GoToZone:
		ret = fmt.Sprintf("Take the player to (%d, %d) in Zone %03d", a.Args[1], a.Args[2], a.Args[0])
	case SetGlobalVar:
		ret = fmt.Sprintf("Set GlobalVar to %d", a.Args[0])
	case AddGlobalVar:
		ret = fmt.Sprintf("Add %d to GlobalVar", a.Args[0])
	case SetRandVar:
		ret = fmt.Sprintf("Set RandVar to %d", a.Args[0])
	case AddToHealth:
		ret = fmt.Sprintf("Give the player %d health", a.Args[0])
	}
	return ret
}
