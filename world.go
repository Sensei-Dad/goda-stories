package main

import (
	"github.com/MasterShizzle/goda-stories/gosoh"
)

// World holds all the various bits of the game that we generate
type GameWorld struct {
	Name     string
	SubAreas map[int]gosoh.MapArea
}

type ViewCoords struct {
	X           float64
	Y           float64
	Width       int
	Height      int
	CurrentArea int
}

func NewWorld() *GameWorld {
	gw := GameWorld{
		Name: "Goda Stories",
	}
	gw.SubAreas = make(map[int]gosoh.MapArea)

	// Make a new Overworld
	// Place the player on Dagobah
	world := gosoh.NewOverworld(10, 10)
	gw.SubAreas[world.Id] = world
	world.PrintMap()

	return &gw
}

func (g *Game) GetCurrentArea() gosoh.MapArea {
	return g.World.SubAreas[g.View.CurrentArea]
}
