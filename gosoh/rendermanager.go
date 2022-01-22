package gosoh

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func ProcessRenderables(screen *ebiten.Image, vpX, vpY, vpOffset float64) {
	for _, result := range drawView.Get() {
		// TODO: Handle loading / unloading different Enty's as the player comes near
		img := result.Components[renderableComp].(*Renderable)
		crtr := result.Components[creatureComp].(*Creature)
		pos := result.Components[positionComp].(*Position)
		cInfo := Creatures[crtr.CreatureId]

		img.Image = Tiles[cInfo.Images[crtr.Facing]]

		op := &ebiten.DrawImageOptions{}
		// Position, in an Entity's case, indicates the center of their bounding box
		op.GeoM.Translate(pos.X-(float64(TileWidth)/2)-vpX+vpOffset, pos.Y-(float64(TileHeight)/2)-vpY+vpOffset)
		screen.DrawImage(img.Image, op)
	}
}
