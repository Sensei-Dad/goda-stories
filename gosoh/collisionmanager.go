package gosoh

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// returns true if these two CollisionBoxes overlap each other
func (a *CollisionBox) Overlaps(b CollisionBox) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// construct a box from the given point and edges
// box size is a fraction of tile size, so e.g. everything at 0.5 would give a collision box of 1 tile
func (c *Collidable) GetBox(centerX, centerY float64) CollisionBox {
	ret := CollisionBox{
		X:      centerX - (c.LeftEdge * float64(TileWidth)),
		Y:      centerY - (c.TopEdge * float64(TileHeight)),
		Width:  (c.LeftEdge + c.RightEdge) * float64(TileWidth),
		Height: (c.LeftEdge + c.RightEdge) * float64(TileHeight),
	}

	return ret
}

func DrawBox(screen *ebiten.Image, box CollisionBox) {
	var clr color.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	ebitenutil.DrawLine(screen, box.X, box.Y, box.X+box.Width, box.Y, clr)
	ebitenutil.DrawLine(screen, box.X+box.Width, box.Y, box.X+box.Width, box.Y+box.Height, clr)
	ebitenutil.DrawLine(screen, box.X+box.Width, box.Y+box.Height, box.X, box.Y+box.Height, clr)
	ebitenutil.DrawLine(screen, box.X, box.Y+box.Height, box.X, box.Y, clr)

	ebitenutil.DrawLine(screen, box.X+0.4*box.Width, box.Y+0.5*box.Height, box.X+0.6*box.Width, box.Y+0.5*box.Height, clr)
	ebitenutil.DrawLine(screen, box.X+0.5*box.Width, box.Y+0.4*box.Height, box.X+0.5*box.Width, box.Y+0.6*box.Height, clr)
}
