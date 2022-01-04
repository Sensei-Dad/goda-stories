package main

import (
	"fmt"

	"github.com/bytearena/ecs"
)

func (g *Game) ProcessMovement() {
	// Check all entities with a Movement comp
	for _, result := range moveView.Get() {
		moves := result.Components[movementComp].(*Movable)
		pos := result.Components[positionComp].(*Position)
		crtr := result.Components[creatureComp].(*Creature)
		img := result.Components[renderableComp].(*Renderable)

		ml := g.World.GetLayers()

		deltaX := pos.X - moves.OldX
		deltaY := pos.Y - moves.OldY

		if crtr.State == Standing && (deltaX != 0 || deltaY != 0) { // Creature wants to move
			crtr.State = Walking          // Do this to avoid repeated inputs
			if ml.CanMove(pos.X, pos.Y) { // Check the map, first
				// Move to that tile if nobody else wants it
				thingCanMove := true
				for _, thing := range collideView.Get() {
					// Check all the collidables for common destinations, except itself
					// This is ugly, but manageable since we're only ever checking against one pool of stuff
					pos2 := thing.Components[positionComp].(*Position)
					col2 := thing.Components[collideComp].(*Collidable)
					if pos2.X == pos.X && pos2.Y == pos.Y && col2.IsBlocking && thing.Entity.ID != result.Entity.ID {
						thingCanMove = false
					}
				}
				if thingCanMove {
					// ...then move the thing!
					fmt.Printf("Starting walk: %s wants to go from (%d,%d) to (%d,%d)\n", crtr.Name, moves.OldX, moves.OldY, pos.X, pos.Y)
				} else {
					// ...otherwise, move it back
					fmt.Printf("%s wants to go from (%d,%d) to (%d,%d), but it's blocked\n", crtr.Name, moves.OldX, moves.OldY, pos.X, pos.Y)
					pos.X = moves.OldX
					pos.Y = moves.OldY
				}
			}
		}
		if crtr.State == Walking {
			// Nudge tile closer to its destination, according to its speed
			// TODO: Set a global game speed
			nudgeX := (float64(deltaX) * moves.Speed) / float64(tileWidth)
			nudgeY := (float64(deltaY) * moves.Speed) / float64(tileHeight)

			img.PixelX += nudgeX
			img.PixelY += nudgeY

			// Detect how far it is to our destination
			// These should approach zero as we get toward it
			distanceX := (float64(pos.X*tileWidth) - img.PixelX) * float64(deltaX)
			distanceY := (float64(pos.Y*tileHeight) - img.PixelY) * float64(deltaY)
			if distanceX <= 0 && distanceY <= 0 {
				// Set the values to be exactly on the tile
				fmt.Printf("Finished walk: %s went from (%d,%d) to (%d,%d)\n", crtr.Name, moves.OldX, moves.OldY, pos.X, pos.Y)
				crtr.State = Standing
				img.PixelX = float64(pos.X * tileWidth)
				img.PixelY = float64(pos.Y * tileHeight)
				moves.OldX = pos.X
				moves.OldY = pos.Y
			}
		}
	}
}

// TODO: Update animation images

func (g *Game) MoveCreature(crtr *ecs.QueryResult, newX, newY int) {
	// Move the creature to the specified coords
	pComp := crtr.Components[positionComp].(*Position)
	pComp.X = newX
	pComp.Y = newY
}
