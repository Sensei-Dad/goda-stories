package gosoh

import (
	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// Global components
var playerView *ecs.View
var moveView *ecs.View
var drawView *ecs.View
var collideView *ecs.View

var playerComp *ecs.Component
var positionComp *ecs.Component
var renderableComp *ecs.Component
var creatureComp *ecs.Component
var movementComp *ecs.Component
var collideComp *ecs.Component

// Components
type PlayerInput struct {
	ShowDebug    bool
	ShowBoxes    bool
	ShowWalkable bool
}

type Creature struct {
	Name       string
	State      CreatureState
	Facing     CardinalDirection
	CanMove    bool
	CreatureId int
}

type PlayerInventory struct {
	Items []int
}

type Position struct {
	X     float64 // Position, in pixels
	Y     float64
	TileX int // Same, but measured with Tiles
	TileY int
}

type Renderable struct {
	Image *ebiten.Image
}

type PlayerAnimation struct {
}

type AnimatedTile struct {
	CurrentFrame int
	FrameDelay   int   // Number of ticks between frames
	Frames       []int // List of tile IDs, for drawing
}

// Defines the dimensions of the Entity's bounding box, in fractions of a tile
type Collidable struct {
	IsBlocking bool
	LeftEdge   float64
	RightEdge  float64
	TopEdge    float64
	BottomEdge float64
}

// Movables can move around the map, in a pixel-wise fashion
type Movable struct {
	Speed     float64
	Direction CardinalDirection
}
