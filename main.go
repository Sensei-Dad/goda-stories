package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const tileWidth, tileHeight = 32, 32 // Tile width and height, in pixels
const viewWidth, viewHeight = 18, 18 // Viewport width and height, in tiles
// const uiPadding = 5                  // Padding between UI elements, in pixels
const tileInfoFile = "assets/text/tileInfo.txt"
const mapInfoHtml = "assets/text/mapInfo.html"

func main() {
	// Process the input file
	// TODO: some more tweaking to not repeat this step too many times
	yodaFile := "YODESK.DTA"
	tileInfo, zoneInfo := processYodaFile(yodaFile)
	// _, _ = processYodaFile(yodaFile)

	// Init the game
	g := NewGame(tileInfo, zoneInfo)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Goda Stories")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
