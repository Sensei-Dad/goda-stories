package main

// Functions to extract and process the data from .DTA resources

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghostiam/binstruct"
)

type ZoneInfo struct {
	Id        int
	Type      string
	Width     int
	Height    int
	Flags     string
	LayerData struct {
		Terrain []int
		Objects []int
		Overlay []int
	}
	TileTriggers []TileTrigger
	Izax         []byte
	Izx2         []byte
	Izx3         []byte
	Izx4         []byte
	Iact         [][]byte
}

type TileInfo struct {
	// TODO: need to process flags in separate groups (TypeFlags, ItemFlags, etc...)?
	Id         string
	Flags      string
	Type       string
	IsWalkable bool
}

type TileTrigger struct {
	Type string
	X    int
	Y    int
	Arg  int
}

func processYodaFile(fileName string) ([]TileInfo, []ZoneInfo) {
	yodaFilePath := "data/" + fileName
	tileFlags := []TileInfo{}
	zones := []ZoneInfo{}

	file, err := os.Open(yodaFilePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := binstruct.NewReader(file, binary.LittleEndian, false)

	defer file.Close()
	fmt.Printf("[%s] Opened file\n", fileName)

	outputs := make(map[string]interface{})

	// Parse the different sections
	for {
		// Grab section header
		_, section, err := reader.ReadBytes(4)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			// fmt.Println("Done.")
			break
		}

		fmt.Printf("[%s] Reading section: %s\n", fileName, section)
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
			zones = make([]ZoneInfo, int(zoneCount))
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
					continue
				} else {
					_, tileBytes, _ := reader.ReadBytes(0x400)
					err = saveByteSliceToPNG(tilename, tileBytes)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			fmt.Printf("    %d tiles extracted, %d skipped\n", numTiles-skipped, skipped)
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
	mapInfo, err := os.Create(mapInfoText)
	if err != nil {
		log.Fatal(err)
	}

	// Save map info to HTML and images to PNGs
	fmt.Printf("[%s] Stitching maps...\n", fileName)
	numZones := len(outputs["ZONE"].([]ZoneInfo))
	skipped := 0
	mapsHtml := htmlStarter

	for zId, zData := range outputs["ZONE"].([]ZoneInfo) {
		mapNum := fmt.Sprintf("%03d", zId)
		mapFilePath := "assets/maps/map_" + mapNum + ".png"

		mapsHtml += getZoneHTML(zData)

		_, err := os.Stat(mapFilePath)
		if err == nil {
			// Skip creating the map if it's already done
			skipped++
		} else {
			// Otherwise, stitch the map together and save it
			err = saveMapToPNG(mapFilePath, zData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("    %d maps extracted, %d skipped.\n", numZones-skipped, skipped)
	mapsHtml += "\n</body>\n</html>\n"

	fmt.Printf("    Dumping output to %s...\n", tileInfoFile)
	spew.Fdump(tileInfo, tileFlags)

	fmt.Printf("    Dumping output to %s...\n", mapInfoHtml)
	spew.Fprint(mapLayers, mapsHtml)

	fmt.Printf("    Dumping output to %s...\n", mapInfoText)
	// Cut the layers, so output is cleaner
	shorterInfo := make([]ZoneInfo, len(zones))
	for i, z := range zones {
		zon := z
		zon.LayerData.Terrain = nil
		zon.LayerData.Objects = nil
		zon.LayerData.Overlay = nil
		shorterInfo[i] = zon
	}
	spew.Fdump(mapInfo, shorterInfo)

	fmt.Printf("[%s] Processed data file.\n", yodaFile)

	return tileFlags, zones
}

func processTileData(tileId int, flags uint32) TileInfo {
	t := TileInfo{}
	t.Id = fmt.Sprintf("%04d", tileId)
	t.Flags = reverse(fmt.Sprintf("%032b", flags))

	// The first 9 bits let us break down what kind of tile this is
	// For now, this just affects collisions
	category := t.Flags[:9]
	switch category {
	case "010000000":
		t.Type = "Terrain"
		t.IsWalkable = true
	case "101000000":
		t.Type = "Object"
		t.IsWalkable = false
	case "001000000":
		t.Type = "Terrain"
		t.IsWalkable = false
	case "101100000":
		t.Type = "Block"
		t.IsWalkable = false
	case "100010000", "000010000":
		t.Type = "Overlay"
		t.IsWalkable = true
	case "100000001":
		t.Type = "Creature"
		t.IsWalkable = false
	case "100000010":
		t.Type = "Item"
		t.IsWalkable = false
	case "100000100":
		t.Type = "Weapon"
		t.IsWalkable = false
	default:
		t.IsWalkable = true
	}
	// Locator minimap tiles, we grab by ID: tiles 817-837
	if (tileId >= 817) && (tileId <= 837) {
		t.Type = "Locator"
	}

	return t
}

func processZoneData(zData []byte, tiles []TileInfo) ZoneInfo {
	z := new(ZoneInfo)

	// Sanity check
	zoneHeader := string(zData[2:6])
	z.Id = int(binary.LittleEndian.Uint16(zData[0:]))
	if zoneHeader != "IZON" {
		log.Fatal(fmt.Sprintf("IZON header not found: cannot parse zoneData for zoneId %s", fmt.Sprint(z.Id)))
	}

	// fmt.Printf("    Processing map_%03d: ", z.Id)

	// Populate a ZoneInfo for this map
	z.Width = int(binary.LittleEndian.Uint16(zData[10:]))
	z.Height = int(binary.LittleEndian.Uint16(zData[12:]))
	z.Flags = reverse(fmt.Sprintf("%08b", zData[14]))

	p := int(zData[20])
	switch p {
	case 1:
		z.Type = "desert"
	case 2:
		z.Type = "snow"
	case 3:
		z.Type = "forest"
	case 5:
		z.Type = "swamp"
	default:
		z.Type = "UNKNOWN"
	}

	// Grab tiles starting at byte 22: each one has 3x two-byte ints, for 3 tiles / cell
	z.LayerData.Terrain = make([]int, z.Width*z.Height)
	z.LayerData.Objects = make([]int, z.Width*z.Height)
	z.LayerData.Overlay = make([]int, z.Width*z.Height)
	for j := 0; j < (z.Width * z.Height); j++ {
		z.LayerData.Terrain[j] = int(binary.LittleEndian.Uint16(zData[6*j+22:]))
		z.LayerData.Objects[j] = int(binary.LittleEndian.Uint16(zData[6*j+24:]))
		z.LayerData.Overlay[j] = int(binary.LittleEndian.Uint16(zData[6*j+26:]))
	}

	// Parse entries for object info
	triggerTypes := []string{
		"trigger_location",
		"spawn_location",
		"force_location",
		"vehicle_to_secondary_map",
		"vehicle_to_primary_map",
		"object_gives_locator",
		"crate_with_item",
		"puzzle_NPC",
		"crate_with_weapon",
		"map_entrance",
		"map_exit",
		"unused",
		"lock",
		"teleporter",
		"xwing_from_dagobah",
		"xwing_to_dagobah",
		"UNKNOWN",
	}
	offset := (6 * z.Width * z.Height) + 22
	numTriggers := int(binary.LittleEndian.Uint16(zData[offset:]))
	if numTriggers > 0 {
		z.TileTriggers = make([]TileTrigger, numTriggers)
		for k := 0; k < numTriggers; k++ {
			z.TileTriggers[k].Type = triggerTypes[int(binary.LittleEndian.Uint16(zData[offset+2:]))]
			z.TileTriggers[k].X = int(binary.LittleEndian.Uint16(zData[offset+6:]))
			z.TileTriggers[k].Y = int(binary.LittleEndian.Uint16(zData[offset+8:]))
			z.TileTriggers[k].Arg = int(binary.LittleEndian.Uint16(zData[offset+12:]))
			offset += 12
		}
	}

	// Advance past the IZAX header and grab action data
	offset += 6
	sectionLength := int(binary.LittleEndian.Uint16(zData[offset:]))
	offset += 2
	z.Izax = zData[offset : offset+sectionLength-6]

	// Advance past the IZX2 header
	offset += sectionLength - 2
	sectionLength = int(binary.LittleEndian.Uint16(zData[offset:]))
	offset += 2
	z.Izx2 = zData[offset : offset+sectionLength-6]

	// ...And again for IZX3
	offset += sectionLength - 2
	sectionLength = int(binary.LittleEndian.Uint16(zData[offset:]))
	offset += 2
	z.Izx3 = zData[offset : offset+sectionLength-6]

	// ...And so on
	offset += sectionLength - 2
	z.Izx4 = zData[offset : offset+8]
	offset += 8

	// Parse actions, if there are any
	for len(zData) >= offset+4 {
		// IACT + sectionLength = 8 bytes
		sectionLength = int(binary.LittleEndian.Uint32(zData[offset+4:]))
		offset += 8
		z.Iact = append(z.Iact, zData[offset:offset+sectionLength])
		offset += sectionLength
	}

	return *z
}

func reverse(str string) (result string) {
	// Given a string, return it in reverse order
	for _, v := range str {
		result = string(v) + result
	}
	return
}
