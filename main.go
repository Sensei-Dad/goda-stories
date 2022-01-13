package main

import (
	"log"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/hajimehoshi/ebiten/v2"
)

const yodaFile = "YODESK.DTA"

func main() {
	// Process the input file, grab tiles and maps
	// TODO:
	//  - some more tweaking to not repeat this step too many times
	//  - action scripts
	//  - worldgen rules?
	var tileInfo []gosoh.TileInfo
	var zoneInfo []gosoh.ZoneInfo
	var itemInfo []gosoh.ItemInfo
	var puzzleInfo []gosoh.PuzzleInfo
	var creatureInfo []gosoh.CreatureInfo

	tileInfo, zoneInfo, itemInfo, puzzleInfo, creatureInfo = processYodaFile(yodaFile, false)

	// Init the game...
	g := NewGame(tileInfo, zoneInfo, itemInfo, puzzleInfo, creatureInfo)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Goda Stories")

	// ...and run it
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
