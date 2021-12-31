package main

import "github.com/hajimehoshi/ebiten/v2"

func ProcessRenderables(g *Game, ml MapLayers, screen *ebiten.Image) {
	for _, result := range g.ECSManager.Query(g.ECSTags["renderables"]) {
		pos := result.Components[position].(*Position)
		img := result.Components[renderable].(Renderable).Image

		tIndex := ml.GetTileIndex(pos.X, pos.Y)
		tile := ml.Objects[tIndex]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
		screen.DrawImage(img, op)
	}
}
