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
type CreatureDirection string

const (
	Standing  CreatureState = "Standing"
	Walking   CreatureState = "Walking"
	Attacking CreatureState = "Attacking"
	Dragging  CreatureState = "Dragging"
)

const (
	Up        CreatureDirection = "Up"
	Down      CreatureDirection = "Down"
	Left      CreatureDirection = "Left"
	Right     CreatureDirection = "Right"
	UpLeft    CreatureDirection = "UpLeft"
	DownLeft  CreatureDirection = "DownLeft"
	UpRight   CreatureDirection = "UpRight"
	DownRight CreatureDirection = "DownRight"
)

type Player struct {
}

type Creature struct {
	Name   string
	State  CreatureState
	Facing CreatureDirection
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

type Movable struct {
	OldX  int
	OldY  int
	Speed float64
}
