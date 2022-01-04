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

	deltaX := 0
	deltaY := 0
	var dir CreatureDirection

	// Detect moves, if not moving
	if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad7) {
		deltaY = -1
		deltaX = -1
		dir = UpLeft
	} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad9) {
		deltaY = -1
		deltaX = 1
		dir = UpRight
	} else if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad1) {
		deltaY = 1
		deltaX = -1
		dir = DownLeft
	} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad3) {
		deltaY = 1
		deltaX = 1
		dir = DownRight
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyNumpad8) {
		deltaY = -1
		dir = Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyNumpad2) {
		deltaY = 1
		dir = Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyNumpad4) {
		deltaX = -1
		dir = Left
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyNumpad6) {
		deltaX = 1
		dir = Right
	}

	for _, result := range playerView.Get() {
		// Eventually, do mouse input as well
		// fmt.Printf("Attempting to process input on %d components\n", len(result.Components))
		pos := result.Components[positionComp].(*Position)
		crtr := result.Components[creatureComp].(*Creature)

		if crtr.State == Standing && (deltaX != 0 || deltaY != 0) {
			pos.X += deltaX
			pos.Y += deltaY
			crtr.Facing = dir
		}
	}
}
