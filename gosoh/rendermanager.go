package gosoh

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessRenderables(screen *ebiten.Image) {
	for _, result := range drawView.Get() {
		// TODO: Handle loading / unloading different Enty's as the player comes near
		img := result.Components[renderableComp].(*Renderable)
		crtr := result.Components[creatureComp].(*Creature)
		pos := result.Components[positionComp].(*Position)
		cInfo := Creatures[crtr.CreatureId]

		img.Image = Tiles[cInfo.Images[crtr.Facing]]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(img.Image, op)
	}
}
