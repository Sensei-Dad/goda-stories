package main

import "github.com/hajimehoshi/ebiten/v2"

func ProcessRenderables(g *Game, ml MapLayers, screen *ebiten.Image) {
	for _, result := range g.ECSManager.Query(g.ECSTags["renderables"]) {
		// TODO: Make sure something is actually within the Viewport
		pos := result.Components[position].(*Position)
		img := result.Components[renderable].(Renderable).Image

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pos.X*tileWidth), float64(pos.Y*tileHeight))
		screen.DrawImage(img, op)
	}
}
