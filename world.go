package main

import (
	"math/rand"
)

// World holds all the various bits of the game that we generate
type GameWorld struct {
	Name       string
	Maps       []MapScreen
	CurrentMap int
}

func NewWorld() GameWorld {
	// For now, choose a random screen from those available and start there
	// TODO: Detect where the starting screens are, find the level branching logic
	msNum := rand.Intn(len(zoneInfo))
	ms := NewMapScreen(msNum)

	maps := make([]MapScreen, 0)
	maps = append(maps, ms)
	gw := GameWorld{Name: "Goda Stories", Maps: maps, CurrentMap: 0}

	return gw
}

func (gw *GameWorld) GetLayers() MapLayers {
	// return the map of whatever the current level is
	return gw.Maps[gw.CurrentMap].Layers
}
