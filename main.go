package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghostiam/binstruct"
	// "github.com/hajimehoshi/ebiten/v2"
)

const tileWidth, tileHeight = 32, 32
const tileInfoFile = "assets/text/tileInfo.txt"
const mapInfoHtml = "assets/text/mapInfo.html"

func main() {
	filename := "YODESK.DTA"
	path := "data/" + filename
	tileFlags := []TileInfo{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	reader := binstruct.NewReader(file, binary.LittleEndian, false)

	defer file.Close()
	fmt.Printf("[%s] Opened file\n", filename)

	outputs := make(map[string]interface{})

	// Parse the different sections
	for {
		// Grab section header
		_, section, err := reader.ReadBytes(4)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			fmt.Println("Done.")
			break
		}

		fmt.Printf("[%s] Reading section: %s\n", filename, section)
		switch s := string(section); s {
		case "VERS":
			major, _ := reader.ReadUint16()
			minor, _ := reader.ReadUint16()
			outputs[s] = fmt.Sprint(int(major)) + "." + fmt.Sprint(int(minor))
		case "STUP", "SNDS", "PUZ2", "CHAR", "CHWP", "CAUX", "TNAM":
			// Basically, just skip all these sections
			sectionLength, _ := reader.ReadUint32()
			_, _, err := reader.ReadBytes(int(sectionLength))
			// _, sectionData, err := reader.ReadBytes(int(sectionLength))
			// outputs[s] = sectionData

			if err != nil {
				fmt.Printf("Error reading section %s\n", section)
				log.Fatal(err)
			}
		case "ZONE":
			zoneCount, _ := reader.ReadUint16()
			zones := make([]ZoneInfo, int(zoneCount))
			for i := 0; i < int(zoneCount); i++ {
				// dunno what this does
				_, _ = reader.ReadUint16()

				zoneLength, _ := reader.ReadUint32()

				_, zoneData, _ := reader.ReadBytes(int(zoneLength))

				zones[i] = processZoneData(zoneData, tileFlags)
			}
			outputs[s] = zones
		case "TILE":
			// Each tile has 4 bytes for the tile data, plus 32x32 px (0x400)
			sectionLength, _ := reader.ReadUint32()
			numTiles := int(sectionLength) / 0x404
			tileFlags = make([]TileInfo, numTiles)
			skipped := 0

			// Extract tile bits into images
			for i := 0; i < numTiles; i++ {
				// Pad number with leading zeroes for filename
				tilename := "assets/tiles/tile_" + fmt.Sprintf("%04d", i) + ".png"
				flags, _ := reader.ReadUint32()
				tileFlags[i] = processTileData(i, flags)

				// Skip creating the tile image if it's already there
				_, err := os.Stat(tilename)
				if err == nil {
					skipped++
					reader.ReadBytes(0x400)
					fmt.Printf(".")
					continue
				} else {
					_, tileBytes, _ := reader.ReadBytes(0x400)
					err = saveByteSliceToPNG(tilename, tileBytes)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("*")
				}
			}
			fmt.Printf("]\n    %d tiles extracted, %d skipped.\n", numTiles-skipped, skipped)
		case "ENDF":
			// Read whatever odd bytes are left?
			_, err = reader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Printf("UNHANDLED CASE: %s\n", section)
			log.Fatal("Unhandled case")
		}
	}

	// create various output files
	tileInfo, err := os.Create(tileInfoFile)
	if err != nil {
		log.Fatal(err)
	}
	mapLayers, err := os.Create(mapInfoHtml)
	if err != nil {
		log.Fatal(err)
	}

	// output map info to file
	fmt.Printf("\n")
	fmt.Printf("[%s] Stitching maps... \n    [", filename)
	numZones := len(outputs["ZONE"].([]ZoneInfo))
	skipped := 0

	// Construct map layer markdown file
	mapsHtml := htmlStarter
	for zId, zData := range outputs["ZONE"].([]ZoneInfo) {
		// Save map image and fill output with zone data
		mapNum := fmt.Sprintf("%03d", zId)
		mapFilePath := "assets/maps/map_" + mapNum + ".png"

		mapsHtml += getZoneHTML(zData)

		_, err := os.Stat(mapFilePath)
		if err == nil {
			// Skip creating the map if it's already done
			skipped++
			fmt.Printf(".")
		} else {
			// Otherwise, stitch the map together and save it
			err = saveMapToPNG(mapFilePath, zData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("*")
		}

		// Prune the map layers, so the output is cleaner
		outputs["ZONE"].([]ZoneInfo)[zId].LayerData.Terrain = nil
		outputs["ZONE"].([]ZoneInfo)[zId].LayerData.Objects = nil
		outputs["ZONE"].([]ZoneInfo)[zId].LayerData.Overlay = nil
	}
	fmt.Println("]")
	fmt.Printf("    %d maps extracted, %d skipped.\n", numZones-skipped, skipped)
	mapsHtml += "\n</body>\n</html>\n"

	fmt.Printf("Dumping output to %s...\n", tileInfoFile)
	spew.Fdump(tileInfo, tileFlags)

	fmt.Printf("Dumping output to %s...\n", mapInfoHtml)
	spew.Fprint(mapLayers, mapsHtml)
}
