package main

import (
	"fmt"
	"math"

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

		// Creature wants to start a move, and isn't
		// currently in the middle of another one
		if !crtr.InMotion && !moves.Direction.NoDirection() {
			newX := pos.X + moves.Direction.DeltaX
			newY := pos.Y + moves.Direction.DeltaY

			// Check the map tile first
			thingCanMove := ml.CanMove(newX, newY)

			// Check all the collidables for common destinations, except for itself
			for _, thing := range collideView.Get() {
				// This is ugly, but manageable since we're only ever checking against one pool of stuff
				pos2 := thing.Components[positionComp].(*Position)
				col2 := thing.Components[collideComp].(*Collidable)
				if pos2.X == newX && pos2.Y == newY && col2.IsBlocking && thing.Entity.ID != result.Entity.ID {
					fmt.Println("Found blocking Entity")
					thingCanMove = false
				}
			}

			if thingCanMove {
				// Yes, start the move
				crtr.InMotion = true
				crtr.State = Walking
				fmt.Printf("Starting walk: %s walking from (%d,%d) to (%d,%d)\n", crtr.Name, pos.X, pos.Y, newX, newY)
				pos.X = newX
				pos.Y = newY
			} else {
				// No, you can't move there
				crtr.InMotion = false
				crtr.State = Standing
				pos.X = moves.OldX
				pos.Y = moves.OldY
				moves.Direction = NoMove
				fmt.Printf("Starting walk: %s wants to go from (%d,%d) to (%d,%d), but it's blocked\n", crtr.Name, pos.X, pos.Y, newX, newY)
			}
		}
		if crtr.InMotion { // Move in-progress
			// Nudge tile closer to its destination, according to its speed
			// TODO: Set a global game speed
			nudgeX := (float64(pos.X-moves.OldX) * moves.Speed)
			nudgeY := (float64(pos.Y-moves.OldY) * moves.Speed)
			img.PixelX += nudgeX
			img.PixelY += nudgeY

			// Detect how far we've left in the move
			distanceX := math.Abs(float64(pos.X*tileWidth) - img.PixelX)
			distanceY := math.Abs(float64(pos.Y*tileHeight) - img.PixelY)

			// Once we've got less than one move left, the move is completed
			if distanceX <= moves.Speed && distanceY <= moves.Speed {
				fmt.Printf("Finished walk: %s went from (%d,%d) to (%d,%d)\n", crtr.Name, moves.OldX, moves.OldY, pos.X, pos.Y)
				img.PixelX = float64(pos.X * tileWidth)
				img.PixelY = float64(pos.Y * tileHeight)
				moves.OldX = pos.X
				moves.OldY = pos.Y
				crtr.InMotion = false
				crtr.State = Standing
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
