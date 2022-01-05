package main

import (
	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// Global components
var playerView *ecs.View
var moveView *ecs.View
var drawView *ecs.View
var collideView *ecs.View

var positionComp *ecs.Component
var renderableComp *ecs.Component
var creatureComp *ecs.Component
var movementComp *ecs.Component
var collideComp *ecs.Component

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

var Up CardinalDirection = CardinalDirection{Name: "Up", DeltaX: 0, DeltaY: -1}
var Down CardinalDirection = CardinalDirection{Name: "Down", DeltaX: 0, DeltaY: 1}
var Left CardinalDirection = CardinalDirection{Name: "Left", DeltaX: -1, DeltaY: 0}
var Right CardinalDirection = CardinalDirection{Name: "Right", DeltaX: 1, DeltaY: 0}
var UpLeft CardinalDirection = CardinalDirection{Name: "UpLeft", DeltaX: -1, DeltaY: -1}
var DownLeft CardinalDirection = CardinalDirection{Name: "DownLeft", DeltaX: -1, DeltaY: 1}
var UpRight CardinalDirection = CardinalDirection{Name: "UpRight", DeltaX: 1, DeltaY: -1}
var DownRight CardinalDirection = CardinalDirection{Name: "DownRight", DeltaX: 1, DeltaY: 1}
var NoMove CardinalDirection = CardinalDirection{Name: "None", DeltaX: 0, DeltaY: 0}

func (d *CardinalDirection) NoDirection() bool {
	// Return true if movement is zero
	return (d.DeltaX == 0 && d.DeltaY == 0)
}

type Player struct {
}

type Creature struct {
	Name       string
	State      CreatureState
	Facing     CardinalDirection
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
