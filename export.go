package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/davecgh/go-spew/spew"
	"github.com/salviati/go-tmx/tmx"
)

const tilesetFileName = "assets/yodatiles.tsx"

func dumpToFile(filepath string, foo ...interface{}) error {
	// Create the output file
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}

	spew.Fdump(file, foo)
	fmt.Printf("[dumpToFile] Saved to file: %s\n", filepath)
	return err
}

func saveXMLToFile(filepath string, foo interface{}) error {
	// Create the output file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}

	file, err := xml.MarshalIndent(foo, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	spew.Fprint(out, file)

	fmt.Printf("[saveXMLToFile] Saved to file: %s\n", filepath)
	return err
}

func saveTiledMaps(g *Game) {
	// Save Tileset and Zones to Tiled-compatible files
	fmt.Println("[saveTiledMaps] Stitching maps...")
	tilesetXML := tmx.Tileset{
		Name:       "Yoda Stories Tileset",
		TileWidth:  gosoh.TileWidth,
		TileHeight: gosoh.TileHeight,
		Spacing:    0,
		Margin:     0,
		Tilecount:  len(gosoh.TileInfos),
		Columns:    gosoh.TilesetColumns,
	}

	for zId, zData := range gosoh.Zones {
		mapNum := fmt.Sprintf("%03d", zId)
		mapFilePath := "assets/maps/map_" + mapNum + ".png"

		// Create a Tiled map for each Zone
		saveZoneToTiledMap(mapFilePath, zData, &tilesetXML)
	}
	fmt.Printf("    %d maps extracted.\n", len(gosoh.Zones))

	saveXMLToFile(tilesetFileName, tilesetXML)
}

func saveZoneToTiledMap(filepath string, zData gosoh.ZoneInfo, tileset *tmx.Tileset) {
	zoneXML := tmx.Map{
		Width:       zData.Width,
		Height:      zData.Height,
		TileWidth:   gosoh.TileWidth,
		TileHeight:  gosoh.TileHeight,
		Orientation: "orthogonal",
	}
	zoneXML.Tilesets = make([]tmx.Tileset, 0)

	zoneXML.Layers = make([]tmx.Layer, 0)
	// Add Terrain tiles

	// lyr := tmx.Layer{}
}
