package gosoh

import (
	"fmt"
	"math"

	"github.com/bytearena/ecs"
)

func ProcessMovement(ms MapScreen) {
	// Check all entities with a Movement comp
	for _, result := range moveView.Get() {
		moves := result.Components[movementComp].(*Movable)
		pos := result.Components[positionComp].(*Position)
		crtr := result.Components[creatureComp].(*Creature)
		img := result.Components[renderableComp].(*Renderable)

		// Creature wants to start a move, and isn't
		// currently in the middle of another one
		if !crtr.InMotion && !moves.Direction.NoDirection() {
			newX := pos.X + moves.Direction.DeltaX
			newY := pos.Y + moves.Direction.DeltaY

			// Check the map tile first
			newTile := ms.GetTileAt(newX, newY)
			thingCanMove := newTile.IsWalkable

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
				moves.OldX = pos.X
				moves.OldY = pos.Y
				pos.X = newX
				pos.Y = newY
			} else {
				// No, you can't move there
				crtr.InMotion = false
				crtr.State = Standing
				pos.X = moves.OldX
				pos.Y = moves.OldY
				moves.Direction = NoMove
			}
		}
		if crtr.InMotion { // Move in-progress
			// Nudge the thing in its chosen direction, according to its speed
			// TODO: Set a global game speed
			nudgeX := (float64(pos.X-moves.OldX) * moves.Speed)
			nudgeY := (float64(pos.Y-moves.OldY) * moves.Speed)
			img.PixelX += nudgeX
			img.PixelY += nudgeY

			// Detect how far we've left in the move
			distanceX := math.Abs(float64(pos.X*TileWidth) - img.PixelX)
			distanceY := math.Abs(float64(pos.Y*TileHeight) - img.PixelY)

			// Once we've got less than one move left, the move is completed
			if distanceX <= moves.Speed && distanceY <= moves.Speed {
				img.PixelX = float64(pos.X * TileWidth)
				img.PixelY = float64(pos.Y * TileHeight)
				moves.OldX = pos.X
				moves.OldY = pos.Y
				crtr.InMotion = false
				moves.Direction = NoMove
			}
		}
	}
}

// TODO: Update animation images

func MoveCreature(crtr *ecs.QueryResult, newX, newY int) {
	// Move the creature to the specified coords
	pComp := crtr.Components[positionComp].(*Position)
	pComp.X = newX
	pComp.Y = newY
}
