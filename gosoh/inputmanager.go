package gosoh

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func ProcessInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// TODO: save/quit functions
		os.Exit(0)
	}

	var toggleDebug bool = false
	// F2 shows Debug info
	if ebiten.IsKeyPressed(ebiten.KeyF2) {
		toggleDebug = true
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
		plyr := result.Components[playerComp].(*PlayerInput)

		// Update facing, if we're able to move
		if crtr.CanMove {
			if dir.IsDirection() {
				crtr.Facing = dir
				crtr.State = Walking
			} else {
				crtr.State = Standing
			}
			// Direction or not, tell the ActionManager what we're trying to do
			mov.Direction = dir
		}

		if toggleDebug {
			if plyr.ShowDebug {
				plyr.ShowDebug = false
			} else {
				plyr.ShowDebug = true
			}
		}
	}
}

func ShowDebugInfo(screen *ebiten.Image, viewX, viewY float64) {
	out := ""
	out += fmt.Sprintf("Viewport: (%0.2f, %0.2f)\n", viewX, viewY)

	// Player info
	px, py, tx, ty := GetPlayerCoords()
	out += fmt.Sprintf("Player: %0.2f, %0.2f (X: %d, Y: %d)\n", px, py, tx, ty)
	for _, result := range playerView.Get() {
		crtr := result.Components[creatureComp].(*Creature)
		out += fmt.Sprintf("State:  %s\nFacing: %s\n", crtr.State, crtr.Facing.Name)
		out += fmt.Sprintf("CanMove:  %t\n", crtr.CanMove)
	}

	ebitenutil.DebugPrint(screen, out)
}

func DrawEntityBoxes(screen *ebiten.Image, viewX, viewY, viewOffset float64) {
	for _, result := range collideView.Get() {
		var box CollisionBox
		col := result.Components[collideComp].(*Collidable)
		pos := result.Components[positionComp].(*Position)
		box = col.GetBox(pos.X-viewX+viewOffset, pos.Y-viewY+viewOffset)
		DrawBox(screen, box)
	}
}
