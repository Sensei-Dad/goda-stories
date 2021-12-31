package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessInput(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// TODO: save/quit functions
		os.Exit(0)
	}

	players := g.ECSTags["players"]
	deltaX := 0
	deltaY := 0

	// Detect which command is being sent
	// Eventually, do mouse input as well
	if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad7) {
		deltaY = -1
		deltaX = -1
	} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad9) {
		deltaY = -1
		deltaX = 1
	} else if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad1) {
		deltaY = 1
		deltaX = -1
	} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad3) {
		deltaY = 1
		deltaX = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyNumpad8) {
		deltaY = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyNumpad2) {
		deltaY = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyNumpad4) {
		deltaX = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyNumpad6) {
		deltaX = 1
	}

	ml := g.World.GetLayers()

	for _, result := range g.ECSManager.Query(players) {
		pos := result.Components[position].(*Position)
		newX := pos.X + deltaX
		newY := pos.Y + deltaY
		// Stay within the map edges
		if newX == Clamp(newX, 0, ml.Width-1) && newY == Clamp(newY, 0, ml.Height-1) {
			// Check the tile we're moving to
			tIndex := ml.GetTileIndex(newX, newY)
			if (tIndex == Clamp(tIndex, 0, len(ml.Tiles)-1)) && (ml.Tiles[tIndex].IsWalkable) {
				pos.X = newX
				pos.Y = newY

				// TODO: Update the viewport to keep us on screen
			}
		}
		// TODO: detect moving to other screens, check collision on the next screen over, etc.
	}
}
