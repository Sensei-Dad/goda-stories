package gosoh

import "fmt"

// Action manager:
// - move stuff around the map

func ProcessMovement(a MapArea) {
	// Check all entities with a Movement comp
	for _, result := range moveView.Get() {
		moves := result.Components[movementComp].(*Movable)
		pos := result.Components[positionComp].(*Position)
		crtr := result.Components[creatureComp].(*Creature)
		col := result.Components[collideComp].(*Collidable)

		if moves.Direction.IsDirection() && crtr.CanMove {
			nudgeX := float64(moves.Direction.DeltaX) * moves.Speed
			nudgeY := float64(moves.Direction.DeltaY) * moves.Speed

			// Do we even need to check for the map edges...?
			newX := pos.X + nudgeX
			newY := pos.Y + nudgeY

			var pBox, tBox CollisionBox
			pBox = col.GetBox(newX, newY)

			// Check all the collidables for common destinations, except for itself
			for _, thing := range collideView.Get() {
				// This is ugly, but manageable since we're only ever checking against one pool of stuff
				pos2 := thing.Components[positionComp].(*Position)
				col2 := thing.Components[collideComp].(*Collidable)
				tBox = col2.GetBox(pos2.X, pos2.Y)
				if col2.IsBlocking && pBox.Overlaps(tBox) && thing.Entity.ID != result.Entity.ID {
					fmt.Println("Found blocking Entity")
				}
			}

			// Move, if possible
			if pBox.X == ClampFloat(pBox.X, 0, float64(a.Width*TileWidth*18)-pBox.Width) {
				pos.X = newX
			}
			if pBox.Y == ClampFloat(pBox.Y, 0, float64(a.Height*TileHeight*18)-pBox.Height) {
				pos.Y = newY
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
func GetPlayerCoords() (X, Y float64) {
	for _, result := range playerView.Get() {
		pos := result.Components[positionComp].(*Position)
		X = pos.X
		Y = pos.Y
	}
	return
}
