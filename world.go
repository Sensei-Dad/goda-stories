package main

import (
	"fmt"

	"github.com/bytearena/ecs"
)

// Global vars
var playerSpeed float64

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

func (g *Game) AddCreature(cInfo CreatureInfo, x, y int) *ecs.Entity {
	// Add a creature to the entity pool
	crtr := g.ECSManager.NewComponent()

	fmt.Printf("[ECSMgr] Adding creature: %s", cInfo.Name)

	return g.ECSManager.NewEntity().
		AddComponent(crtr, &Creature{
			Name:   cInfo.Name,
			State:  Standing,
			Facing: Down,
		}).
		AddComponent(renderableComp, &Renderable{
			Image: tiles[2037], // ALL JAWAS, ALL THE TIME
		}).
		AddComponent(positionComp, &Position{
			X: x,
			Y: y,
		}).
		AddComponent(movementComp, &Movable{
			OldX:  x,
			OldY:  y,
			Speed: playerSpeed,
		}).
		AddComponent(collideComp, &Collidable{
			IsBlocking: true,
		})
}

func (g *Game) InitializeWorld() (*ecs.Manager, map[string]ecs.Tag) {
	// Initialize the world via the ECS
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	// Make the global components and add the Player
	player := manager.NewComponent()
	creatureComp = manager.NewComponent()
	renderableComp = manager.NewComponent()
	movementComp = manager.NewComponent()
	positionComp = manager.NewComponent()
	collideComp = manager.NewComponent()

	// TODO: actually try to place the player on a movable tile
	playerX := 4
	playerY := 7
	playerSpeed = 60

	manager.NewEntity().
		AddComponent(player, &Player{}).
		AddComponent(creatureComp, &Creature{
			Name:   creatureInfo[0].Name,
			State:  Standing,
			Facing: Down,
		}).
		AddComponent(renderableComp, &Renderable{
			Image:  tiles[799], // TODO: Suss out where Luke's sprites are => create animations
			PixelX: float64(playerX * tileWidth),
			PixelY: float64(playerY * tileHeight),
		}).
		AddComponent(movementComp, &Movable{
			OldX:  4,
			OldY:  7,
			Speed: playerSpeed,
		}).
		AddComponent(positionComp, &Position{
			X: 4,
			Y: 7,
		}).AddComponent(collideComp, &Collidable{
		IsBlocking: true,
	})

	players := ecs.BuildTag(player, creatureComp, positionComp)
	tags["players"] = players
	playerView = manager.CreateView(players)

	renderables := ecs.BuildTag(renderableComp)
	tags["renderables"] = renderables
	drawView = manager.CreateView(renderables)

	creatures := ecs.BuildTag(creatureComp, positionComp)
	tags["creatures"] = creatures

	movables := ecs.BuildTag(movementComp, positionComp, creatureComp, renderableComp)
	tags["movables"] = movables
	moveView = manager.CreateView(movables)

	collidables := ecs.BuildTag(collideComp, positionComp)
	tags["collidables"] = collidables
	collideView = manager.CreateView(collidables)

	return manager, tags
}
