package main

import (
	"log"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/hajimehoshi/ebiten/v2"
)

const yodaFile = "YODESK.DTA"

func main() {
	// TODO:
	//  - Figure out how to grab YODESK.dta and sounds and junk
	//      via extracting from the ISO directly
	//  	- Make an Installer
	//  - action scripts
	//  - worldgen rules
	var tileset *ebiten.Image
	var tileInfo []gosoh.TileInfo
	var zoneInfo []gosoh.ZoneInfo
	var itemInfo []gosoh.ItemInfo
	var puzzleInfo []gosoh.PuzzleInfo
	var creatureInfo []gosoh.CreatureInfo
	var soundList []string

	tileset, tileInfo, zoneInfo, itemInfo, puzzleInfo, creatureInfo, soundList = processYodaFile(yodaFile)

	// Init the game
	g := NewGame(tileset, tileInfo, zoneInfo, itemInfo, puzzleInfo, creatureInfo, soundList)

	// Create various output files
	// TODO: This should probably be in the installer, or optional
	if true {
		saveTiledMaps(g)
	}

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Goda Stories")
	ebiten.MaximizeWindow()

	// Run it!
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
