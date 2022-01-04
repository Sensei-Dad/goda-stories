package main

import (
	"os"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessInput(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// TODO: save/quit functions
		os.Exit(0)
	}
	var pos *Position
	var crtr *Creature

	deltaX := 0
	deltaY := 0
	ml := g.World.GetLayers()

	for _, result := range playerView.Get() {
		// Eventually, do mouse input as well
		// fmt.Printf("Attempting to process input on %d components\n", len(result.Components))
		pos = result.Components[positionComp].(*Position)
		crtr = result.Components[creatureComp].(*Creature)

		// Detect which command is being sent
		if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad7) {
			deltaY = -1
			deltaX = -1
			crtr.Facing = UpLeft
			crtr.State = Walking
		} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyUp)) || ebiten.IsKeyPressed(ebiten.KeyNumpad9) {
			deltaY = -1
			deltaX = 1
			crtr.Facing = UpRight
			crtr.State = Walking
		} else if (ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad1) {
			deltaY = 1
			deltaX = -1
			crtr.Facing = DownLeft
			crtr.State = Walking
		} else if (ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyDown)) || ebiten.IsKeyPressed(ebiten.KeyNumpad3) {
			deltaY = 1
			deltaX = 1
			crtr.Facing = DownRight
			crtr.State = Walking
		} else if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyNumpad8) {
			deltaY = -1
			crtr.Facing = Up
			crtr.State = Walking
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyNumpad2) {
			deltaY = 1
			crtr.Facing = Down
			crtr.State = Walking
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyNumpad4) {
			deltaX = -1
			crtr.Facing = Left
			crtr.State = Walking
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyNumpad6) {
			deltaX = 1
			crtr.Facing = Right
			crtr.State = Walking
		}

		newX := pos.X + deltaX
		newY := pos.Y + deltaY

		// Stay within the map edges
		if newX == Clamp(newX, 0, ml.Width-1) && newY == Clamp(newY, 0, ml.Height-1) {
			// Check the tile we're moving to
			tIndex := ml.GetTileIndex(newX, newY)
			if (tIndex == Clamp(tIndex, 0, len(ml.Tiles)-1)) && (ml.Tiles[tIndex].IsWalkable) {
				g.MoveCreature(result, newX, newY)
				// TODO: Update the viewport to keep us on screen
			}
		}
		// TODO: detect moving to other screens, check collision on the next screen over, etc.
	}
}

func (g *Game) MoveCreature(crtr *ecs.QueryResult, newX, newY int) {
	// Move the creature to the specified coords
	pComp := crtr.Components[positionComp].(*Position)
	pComp.X = newX
	pComp.Y = newY
}
