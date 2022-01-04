package main

import "github.com/hajimehoshi/ebiten/v2"

func ProcessRenderables(g *Game, ml MapLayers, screen *ebiten.Image) {
	for _, result := range drawView.Get() {
		// TODO: Make sure something is actually within the Viewport
		img := result.Components[renderableComp].(*Renderable)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(img.PixelX), float64(img.PixelY))
		screen.DrawImage(img.Image, op)
	}
}
