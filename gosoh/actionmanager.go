package gosoh

import (
	"fmt"
	"math"
)

// Action manager:
// - move stuff around the map

func ProcessMovement(a *MapArea) {
	// Check all entities with a Movement comp
	for _, result := range moveView.Get() {
		moves := result.Components[movementComp].(*Movable)
		pos := result.Components[positionComp].(*Position)
		crtr := result.Components[creatureComp].(*Creature)
		// col := result.Components[collideComp].(*Collidable)

		if moves.Direction.IsDirection() && crtr.CanMove {
			// Calculate new position
			newX := pos.TileX + moves.Direction.DeltaX
			newY := pos.TileY + moves.Direction.DeltaY

			// Check the map
			tileIsOpen := a.Tiles[newX][newY].IsWalkable
			// Check all the collidables for common destinations, except for itself
			// TODO: Only check the active ones
			for _, thing := range collideView.Get() {
				pos2 := thing.Components[positionComp].(*Position)
				if pos2.TileX == newX && pos2.TileY == newY {
					tileIsOpen = false
				}
			}
			if !tileIsOpen {
				// TODO: Send "Bump" event
				fmt.Printf("Tile (%d, %d) is blocked", newX, newY)
				crtr.CanMove = true
				crtr.State = Standing
			} else {
				// Lock the creature against further actions, until it finishes its move
				crtr.CanMove = false
				crtr.State = Walking
				fmt.Printf("Starting walk: %s walking from (%d,%d) to (%d,%d)\n", crtr.Name, pos.TileX, pos.TileY, newX, newY)
				pos.TileX = newX
				pos.TileY = newY
			}
		}

		if crtr.State == Walking {
			// Move in-progress: Nudge the thing toward its destination, according to its speed
			// Higher speed => fewer ticks to complete a move
			// Detect how far we've left in the move
			distanceX := math.Abs(pos.X - float64(pos.TileX*TileWidth) - float64(TileWidth/2))
			distanceY := math.Abs(pos.Y - float64(pos.TileY*TileHeight) - float64(TileHeight/2))

			// If we've got less than one move left, finish the move
			if distanceX <= moves.Speed && distanceY <= moves.Speed {
				// TODO: Send "entered tile" event?
				fmt.Printf("Finished walk: %s enters (%d, %d)\n", crtr.Name, pos.TileX, pos.TileY)
				pos.X = float64(pos.TileX*TileWidth) + float64(TileWidth/2)
				pos.Y = float64(pos.TileY*TileHeight) + float64(TileHeight/2)
				crtr.CanMove = true
			} else {
				// If not, nudge the thing closer
				pos.X += float64(moves.Direction.DeltaX) * moves.Speed
				pos.Y += float64(moves.Direction.DeltaY) * moves.Speed
			}
		}
	}
}

// Chess term: adjust the "piece" within its own tile, without moving it
// func Jadoube(tPos *Position) {
// 	moves := result.Components[movementComp].(*Movable)
// 	crtr := result.Components[creatureComp].(*Creature)
// 	img := result.Components[renderableComp].(*Renderable)
// 	// Detect how far we've left in the move
// 	distanceX := math.Abs(float64(tPos.X*TileWidth) - img.PixelX)
// 	distanceY := math.Abs(float64(tPos.Y*TileHeight) - img.PixelY)

// 	// Once we've got less than one move left, the move is completed
// 	if distanceX <= moves.Speed && distanceY <= moves.Speed {
// 		img.PixelX = float64(pos.X * TileWidth)
// 		img.PixelY = float64(pos.Y * TileHeight)
// 		moves.OldX = pos.X
// 		moves.OldY = pos.Y
// 		crtr.InMotion = false
// 		moves.Direction = NoMove
// 	}
// }

// TODO: Update animation images
func GetPlayerCoords() (X, Y float64, tX, tY int) {
	for _, result := range playerView.Get() {
		pos := result.Components[positionComp].(*Position)
		X = pos.X
		Y = pos.Y
		tX = pos.TileX
		tY = pos.TileY
	}
	return
}
