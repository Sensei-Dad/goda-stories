package main

import (
	"github.com/MasterShizzle/goda-stories/gosoh"
)

// World holds all the various bits of the game that we generate
type GameWorld struct {
	Name        string
	SubAreas    []*gosoh.MapArea
	CurrentArea int
}

// Coordinates of the Viewport; all measurements are in pixels
type ViewCoords struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func NewWorld() *GameWorld {
	gw := GameWorld{
		Name: "Goda Stories",
	}
	gw.SubAreas = make([]*gosoh.MapArea, 0)

	// Make a new Overworld
	// Place the player on Dagobah
	world := gosoh.NewOverworld(10, 10)
	gw.SubAreas = append(gw.SubAreas, world)
	gw.CurrentArea = world.Id
	world.PrintMap()

	return &gw
}

func (gw *GameWorld) GetCurrentArea() *gosoh.MapArea {
	return gw.SubAreas[gw.CurrentArea]
}
