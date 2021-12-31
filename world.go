package main

import (
	"github.com/bytearena/ecs"
)

// Global components
var position *ecs.Component
var renderable *ecs.Component

// World holds all the various bits of the game that we generate
type GameWorld struct {
	Name       string
	Maps       []MapScreen
	CurrentMap int
}

func NewWorld() GameWorld {
	// For now, choose a random screen from those available and start there
	// TODO: actual worldgen
	msNum := RandomInt(len(zoneInfo))
	ms := NewMapScreen(msNum)

	maps := make([]MapScreen, 0)
	maps = append(maps, ms)
	gw := GameWorld{
		Name:       "Goda Stories",
		Maps:       maps,
		CurrentMap: 0,
	}

	return gw
}

func (gw *GameWorld) GetLayers() MapLayers {
	// return the map of whatever the current level is
	return gw.Maps[gw.CurrentMap].Layers
}

func (gw *GameWorld) GetScreen() MapScreen {
	// return the current MapScreen
	return gw.Maps[gw.CurrentMap]
}

func InitializeWorld() (*ecs.Manager, map[string]ecs.Tag) {
	// Initialize the world via the ECS
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	// Make stuff!
	player := manager.NewComponent()
	// TODO: Suss out where Luke's sprites are => create animations
	playerImg := tiles[799]
	position = manager.NewComponent()
	renderable = manager.NewComponent()
	movable := manager.NewComponent()

	manager.NewEntity().
		AddComponent(player, Player{}).
		AddComponent(renderable, Renderable{
			Image: playerImg,
		}).
		AddComponent(movable, Movable{}).
		AddComponent(position, &Position{
			X: 4,
			Y: 7,
		})

	players := ecs.BuildTag(player, position)
	tags["players"] = players

	renderables := ecs.BuildTag(renderable, position)
	tags["renderables"] = renderables

	return manager, tags
}
