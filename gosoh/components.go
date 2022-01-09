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
}

type Creature struct {
	Name       string
	State      CreatureState
	Facing     CardinalDirection
	InMotion   bool
	CreatureId int
}

type PlayerInventory struct {
	Items []int
}

type Position struct {
	X int
	Y int
}

type Renderable struct {
	Image  *ebiten.Image
	PixelX float64
	PixelY float64
}

type PlayerAnimation struct {
}

type AnimatedTile struct {
	CurrentFrame int
	FrameDelay   int   // Number of ticks between frames
	Frames       []int // List of tile IDs, for drawing
}

type Collidable struct {
	IsBlocking bool
}

// Movables can move around the map, in a pixel-wise fashion
type Movable struct {
	OldX      int
	OldY      int
	Speed     float64
	Direction CardinalDirection
}
