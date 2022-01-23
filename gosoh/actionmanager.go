package gosoh

import "fmt"

// Action manager:
// - move stuff around the map

func ProcessMovement(a *MapArea) {
	// Check all entities with a Movement comp
	for _, result := range moveView.Get() {
		moves := result.Components[movementComp].(*Movable)
		pos := result.Components[positionComp].(*Position)
		crtr := result.Components[creatureComp].(*Creature)
		col := result.Components[collideComp].(*Collidable)

		// Update the tile we're considered to be "in"
		pos.TileX = int(pos.X / float64(TileWidth))
		pos.TileY = int(pos.Y / float64(TileHeight))

		if moves.Direction.IsDirection() && crtr.CanMove {
			nudgeX := float64(moves.Direction.DeltaX) * moves.Speed
			nudgeY := float64(moves.Direction.DeltaY) * moves.Speed

			newX := pos.X + nudgeX
			newY := pos.Y + nudgeY

			var pBox, tBox CollisionBox
			pBox = col.GetBox(newX, newY)

			// Check which tiles overlap at the new location
			topLeft, topRight, bottomLeft, bottomRight := a.CheckCorners(pBox)
			var horizOk, vertOk bool

			// There's probably a more elegant way to loop this
			// TODO: Check for diagonals, allow "jumps"
			if moves.Direction.DeltaX > 0 { // Look rightward
				if moves.Direction.IsDiagonal() {
					if moves.Direction.DeltaY < 0 {
						horizOk = !topRight
					} else {
						horizOk = !bottomRight
					}
				} else {
					horizOk = !(topRight || bottomRight)
				}
			} else { // Leftward
				if moves.Direction.IsDiagonal() {
					if moves.Direction.DeltaY < 0 {
						horizOk = !topLeft
					} else {
						horizOk = !bottomLeft
					}
				} else {
					horizOk = !(topLeft || bottomLeft)
				}
			}
			if moves.Direction.DeltaY > 0 { // Upward
				if moves.Direction.IsDiagonal() {
					if moves.Direction.DeltaX < 0 {
						vertOk = !topLeft
					} else {
						vertOk = !topRight
					}
				} else {
					vertOk = !(topRight || topLeft)
				}
			} else { // Downward
				if moves.Direction.IsDiagonal() {
					if moves.Direction.DeltaX < 0 {
						vertOk = !bottomLeft
					} else {
						vertOk = !bottomRight
					}
				} else {
					vertOk = !(bottomRight || bottomLeft)
				}
			}

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

			fmt.Printf("HMove: %t, VMove: %t\n", horizOk, vertOk)

			// Move, if possible
			if pBox.X == ClampFloat(pBox.X, 0, float64(a.Width*TileWidth*18)) && horizOk {
				pos.X = newX
			}
			if pBox.Y == ClampFloat(pBox.Y, 0, float64(a.Height*TileHeight*18)) && vertOk {
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
