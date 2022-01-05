package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) ProcessInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// TODO: save/quit functions
		os.Exit(0)
	}

	// Without inputs, assume we're standing still
	dir := NoMove

	// Detect inputs
	if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad7) {
		dir = UpLeft
	} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad9) {
		dir = UpRight
	} else if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad1) {
		dir = DownLeft
	} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad3) {
		dir = DownRight
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyNumpad8) {
		dir = Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyNumpad2) {
		dir = Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyNumpad4) {
		dir = Left
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyNumpad6) {
		dir = Right
	}

	for _, result := range playerView.Get() {
		// Eventually, do mouse input as well
		// fmt.Printf("Attempting to process input on %d components\n", len(result.Components))
		mov := result.Components[movementComp].(*Movable)
		crtr := result.Components[creatureComp].(*Creature)

		mov.Direction = dir

		// Change facing and attempt a move, if we're not already moving
		if crtr.State != InMotion && !dir.NoDirection() {
			crtr.Facing = dir
			crtr.State = Walking
		}
	}
}
