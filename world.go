package main

import (
	"math/rand"

	"github.com/bytearena/ecs"
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

func InitializeWorld() (*ecs.Manager, map[string]ecs.Tag) {
	// Initialize the world via the ECS
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	// Make stuff!
	player := manager.NewComponent()
	position := manager.NewComponent()
	renderable := manager.NewComponent()
	movable := manager.NewComponent()

	manager.NewEntity().
		AddComponent(player, Player{}).
		AddComponent(renderable, Renderable{}).
		AddComponent(movable, Movable{}).
		AddComponent(position, &Position{
			X: 5,
			Y: 5,
		})

	players := ecs.BuildTag(player, position)
	tags["players"] = players

	return manager, tags
}
