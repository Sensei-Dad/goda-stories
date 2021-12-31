package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func TryMovePlayer(g *Game) {
	players := g.ECSTags["players"]
	x := 0
	y := 0

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x = 1
	}

	ml := g.World.GetLayers()

	// TODO: Start the player on a walkable square?

	for _, result := range g.ECSManager.Query(players) {
		pos := result.Components[position].(*Position)
		// Stay within the map edges
		if pos.X+x == Clamp(pos.X+x, 0, ml.Width-1) && pos.Y+y == Clamp(pos.Y+y, 0, ml.Height-1) {
			// Check the tile we're moving to
			tIndex := ml.GetTileIndex(pos.X+x, pos.Y+y)
			if CheckIsWalkable(ml.Objects[tIndex].Id) {
				pos.X += x
				pos.Y += y
			} else {
				fmt.Printf("[inputmanager] Player at (%d, %d) colliding with tile or map edge at map_%04d(%d, %d)\n", pos.X, pos.Y, ml.MapId, pos.X+x, pos.Y+y)
			}
		}
		// TODO: detect moving to other screens, check collision on the next screen over, etc.
	}
}
