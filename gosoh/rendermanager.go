package gosoh

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessRenderables(screen *ebiten.Image) {
	for _, result := range drawView.Get() {
		// TODO: Make sure something is actually within the Viewport
		img := result.Components[renderableComp].(*Renderable)
		crtr := result.Components[creatureComp].(*Creature)
		cInfo := Creatures[crtr.CreatureId]

		img.Image = Tiles[cInfo.Images[crtr.Facing]]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(img.PixelX, img.PixelY)
		screen.DrawImage(img.Image, op)
	}
}
