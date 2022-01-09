package gosoh

import (
	"fmt"

	"github.com/bytearena/ecs"
)

func AddCreature(cInfo CreatureInfo, x, y int) *ecs.Entity {
	// Add a creature to the entity pool
	crtr := ECSManager.NewComponent()

	fmt.Printf("[ECSMgr] Adding creature: %s", cInfo.Name)

	return ECSManager.NewEntity().
		AddComponent(crtr, &Creature{
			Name:   cInfo.Name,
			State:  Standing,
			Facing: Down,
		}).
		AddComponent(renderableComp, &Renderable{
			Image: Tiles[cInfo.Images[Down]], // ALL JAWAS, ALL THE TIME
		}).
		AddComponent(positionComp, &Position{
			X:     float64(x*TileWidth) + 0.5, // Spawn in the center of the tile
			Y:     float64(y*TileHeight) + 0.5,
			TileX: x,
			TileY: y,
		}).
		AddComponent(movementComp, &Movable{
			Speed: 2.0,
		}).
		AddComponent(collideComp, &Collidable{
			IsBlocking: true,
		})
}
