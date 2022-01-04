package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const yodaFile = "YODESK.DTA"

var tileInfo []TileInfo
var zoneInfo []ZoneInfo
var itemInfo []ItemInfo
var puzzleInfo []PuzzleInfo
var creatureInfo []CreatureInfo

func main() {
	// Process the input file, grab tiles and maps
	// TODO:
	//  - some more tweaking to not repeat this step too many times
	//  - action scripts
	//  - worldgen rules?
	tileInfo, zoneInfo, itemInfo, puzzleInfo, creatureInfo = processYodaFile(yodaFile)

	// Init the game...
	g := NewGame(tileInfo, zoneInfo)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Goda Stories")

	// ...and run it
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
