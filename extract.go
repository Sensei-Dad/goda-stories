package main

// Functions to extract and process the data from .DTA resources

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/davecgh/go-spew/spew"
	"github.com/ghostiam/binstruct"
)

func processYodaFile(fileName string, dumpOutputs bool) ([]gosoh.TileInfo, []gosoh.ZoneInfo, []gosoh.ItemInfo, []gosoh.PuzzleInfo, []gosoh.CreatureInfo) {
	yodaFilePath := "data/" + fileName
	outTiles := make([]gosoh.TileInfo, 0)
	outZones := make([]gosoh.ZoneInfo, 0)
	outItems := make([]gosoh.ItemInfo, 0)
	outPuzzles := make([]gosoh.PuzzleInfo, 0)
	outCreatures := make([]gosoh.CreatureInfo, 0)

	file, err := os.Open(yodaFilePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := binstruct.NewReader(file, binary.LittleEndian, false)

	defer file.Close()
	fmt.Printf("[%s] Opened file\n", fileName)

	numTiles := 0

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
			_, _ = reader.ReadByte()
			major, _ := reader.ReadByte()
			_, _ = reader.ReadByte()
			minor, _ := reader.ReadByte()
			fmt.Printf("    Detected version: %d.%d\n", major, minor)
		case "STUP", "SNDS", "CHWP", "CAUX":
			// Basically, just skip all these sections
			sectionLength, _ := reader.ReadUint32()
			_, _, err := reader.ReadBytes(int(sectionLength))
			// _, sectionData, err := reader.ReadBytes(int(sectionLength))
			// outputs[s] = sectionData

			if err != nil {
				fmt.Printf("Error reading section %s\n", s)
				log.Fatal(err)
			}
		case "ZONE":
			zoneCount, _ := reader.ReadUint16()
			for i := 0; i < int(zoneCount); i++ {
				// dunno what this does
				_, _ = reader.ReadUint16()

				zoneLength, _ := reader.ReadUint32()

				_, zoneData, _ := reader.ReadBytes(int(zoneLength))

				zInfo := processZoneData(zoneData, outTiles)
				outZones = append(outZones, zInfo)
			}
		case "TILE":
			// Each tile has 4 bytes for the tile data, plus 32x32 px (0x400)
			sectionLength, _ := reader.ReadUint32()
			numTiles = int(sectionLength) / 0x404
			skipped := 0

			// Extract tile bits into images
			for i := 0; i < numTiles; i++ {
				// Pad number with leading zeroes for filename
				tilename := fmt.Sprintf("assets/tiles/tile_%04d.png", i)
				flags, _ := reader.ReadUint32()
				tData := processTileData(i, flags)
				outTiles = append(outTiles, tData)

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
		case "PUZ2":
			sectionLength, _ := reader.ReadUint32()
			_, puzzleData, err := reader.ReadBytes(int(sectionLength))
			puzzles := processPuzzleData(puzzleData, numTiles)
			if err != nil {
				fmt.Printf("Error reading section %s\n", s)
				log.Fatal(err)
			}
			outPuzzles = puzzles
		case "TNAM":
			sectionLength, _ := reader.ReadUint32()
			_, itemData, err := reader.ReadBytes(int(sectionLength))
			items := processItemList(itemData)
			if err != nil {
				fmt.Printf("Error reading section %s\n", s)
				log.Fatal(err)
			}
			outItems = items
		case "CHAR":
			sectionLength, _ := reader.ReadUint32()
			_, itemData, err := reader.ReadBytes(int(sectionLength))
			chars := processCharList(itemData)
			if err != nil {
				fmt.Printf("Error reading section %s\n", s)
				log.Fatal(err)
			}
			outCreatures = chars
		case "ENDF":
			// Read whatever odd bytes are left?
			_, err = reader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Printf("UNHANDLED CASE: %s\n", s)
			log.Fatal("Unhandled case")
		}
	}

	// create various output files
	if dumpOutputs {
		err = dumpToFile(tileInfoFile, outTiles)
		if err != nil {
			log.Fatal(err)
		}
		err = dumpToFile(itemInfoFile, outItems)
		if err != nil {
			log.Fatal(err)
		}
		err = dumpToFile(puzzleInfoFile, outPuzzles)
		if err != nil {
			log.Fatal(err)
		}
		err = dumpToFile(crtrInfoText, outCreatures)
		if err != nil {
			log.Fatal(err)
		}

		// Save action triggers
		outTriggers := ""
		for i, z := range outZones {
			for j, t := range z.ActionTriggers {
				outTriggers += fmt.Sprintf("%03d,%02d,%s \n", i, j, spew.Sdump(t))
			}
		}
		err = printToFile(mapInfoText, outTriggers)
		if err != nil {
			log.Fatal(err)
		}

		saveHTMLMaps(outZones, outItems, outCreatures)
	}

	fmt.Printf("[%s] Processed data file.\n", yodaFile)

	return outTiles, outZones, outItems, outPuzzles, outCreatures
}

func processTileData(tileId int, flags uint32) gosoh.TileInfo {
	t := gosoh.TileInfo{}
	t.Id = tileId
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

func processZoneData(zData []byte, tiles []gosoh.TileInfo) gosoh.ZoneInfo {
	z := gosoh.ZoneInfo{}

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

	zt := int(zData[14])
	switch zt {
	case 1:
		// TODO: pick the Teleporter maps out of here
		z.Type = "Plain"
		z.IsOverworld = true
	case 2:
		z.Type = "GateToNorth"
		z.IsOverworld = true
	case 3:
		z.Type = "GateToSouth"
		z.IsOverworld = true
	case 4:
		z.Type = "GateToEast"
		z.IsOverworld = true
	case 5:
		z.Type = "GateToWest"
		z.IsOverworld = true
	case 6:
		z.Type = "PortalEnter"
		z.IsOverworld = true
	case 7:
		z.Type = "PortalExit"
		z.IsOverworld = true
	case 8:
		z.Type = "Interior"
		z.IsOverworld = false
	case 9:
		z.Type = "OpeningSplash"
		z.IsOverworld = false
	case 10:
		z.Type = "FinalDestination"
		z.IsOverworld = true
	case 11:
		z.Type = "HomeBase"
		z.IsOverworld = true
	case 13:
		z.Type = "WinSplash"
		z.IsOverworld = false
	case 14:
		z.Type = "LoseSplash"
		z.IsOverworld = false
	case 15:
		z.Type = "ItemForTool"
		z.IsOverworld = true
	case 16:
		z.Type = "ItemForItem"
		z.IsOverworld = true
	case 17:
		z.Type = "ItemForTask"
		z.IsOverworld = true
	case 18:
		z.Type = "FindTheForce"
		z.IsOverworld = true
	}

	p := int(zData[20])
	switch p {
	case 1:
		z.Biome = "desert"
	case 2:
		z.Biome = "snow"
	case 3:
		z.Biome = "forest"
	case 5:
		z.Biome = "swamp"
	default:
		z.Biome = "UNKNOWN"
	}

	// Grab tiles starting at byte 22: each one has 3x two-byte ints, for 3 tiles / cell
	z.TileMaps.Terrain = make([]int, z.Width*z.Height)
	z.TileMaps.Objects = make([]int, z.Width*z.Height)
	z.TileMaps.Overlay = make([]int, z.Width*z.Height)
	for j := 0; j < (z.Width * z.Height); j++ {
		z.TileMaps.Terrain[j] = int(binary.LittleEndian.Uint16(zData[6*j+22:]))
		z.TileMaps.Objects[j] = int(binary.LittleEndian.Uint16(zData[6*j+24:]))
		z.TileMaps.Overlay[j] = int(binary.LittleEndian.Uint16(zData[6*j+26:]))
	}

	offset := (6 * z.Width * z.Height) + 22

	// Parse entries for tile-based triggers
	triggerTypes := []string{
		"trigger_location",
		"spawn_location",
		"force_location",
		"vehicle_to_secondary_map",
		"vehicle_to_primary_map",
		"object_gives_locator",
		"item_related_trigger",
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

	numTriggers := int(binary.LittleEndian.Uint16(zData[offset:]))
	if numTriggers > 0 {
		z.TileTriggers = make([]gosoh.TileTrigger, numTriggers)
		for k := 0; k < numTriggers; k++ {
			z.TileTriggers[k].Type = triggerTypes[int(binary.LittleEndian.Uint16(zData[offset+2:]))]
			z.TileTriggers[k].X = int(binary.LittleEndian.Uint16(zData[offset+6:]))
			z.TileTriggers[k].Y = int(binary.LittleEndian.Uint16(zData[offset+8:]))
			z.TileTriggers[k].Arg = int(binary.LittleEndian.Uint16(zData[offset+12:]))
			offset += 12
		}
	}

	// IZAX: Zone Actors (e.g. enemy creatures wandering around the map when it loads)
	// 4B header, 2B section length
	//   2B unused, 2B Unknown, and 2B to count X 44B commands afterward
	//   X * 44B Actors
	//     2B Creature ID
	//     4B X and Y coord on the map where it spawns
	//     6B Args
	//     ...and the rest is usually just FF? What are the rest of these bytes for?
	// 4B unused
	offset += 6
	sectionLength := int(binary.LittleEndian.Uint16(zData[offset:]))
	offset += 2
	izax := zData[offset : offset+sectionLength-6]
	numItems := int(binary.LittleEndian.Uint16(zData[offset+4:]))
	z.ZoneActors = make([]gosoh.ZoneActor, numItems)
	if numItems > 0 {
		for i := 0; i < numItems; i++ {
			zax := gosoh.ZoneActor{
				CreatureId: int(binary.LittleEndian.Uint16(izax[6+(44*i):])),
				ZoneX:      int(binary.LittleEndian.Uint16(izax[8+(44*i):])),
				ZoneY:      int(binary.LittleEndian.Uint16(izax[10+(44*i):])),
				Args:       izax[12+(44*i) : 18+(44*i)],
			}
			chk := 0
			for _, x := range izax[18+(44*i) : 50+(44*i)] {
				chk += int(x)
			}

			if chk != 8160 { // 32 * 0xFF is all 'empties'
				zax.Unknown = (izax[18+(44*i) : 50+(44*i)])
			}

			z.ZoneActors[i] = zax
		}
	}

	// IZX2: Item rewards
	// 6B Header + section length
	// 2B Number of items
	//   2B Item ID
	offset += sectionLength - 2
	sectionLength = int(binary.LittleEndian.Uint16(zData[offset:]))
	offset += 2
	// How many reward items?
	numItems = int(binary.LittleEndian.Uint16(zData[offset+2:]))
	z.RewardItems = make([]int, numItems)
	for i := 0; i < numItems; i++ {
		z.RewardItems[i] = int(binary.LittleEndian.Uint16(zData[offset+4+(2*i):]))
	}

	// IZX3: Quest-related NPCs
	// 6B Header + section length
	// 2B Number of items
	//   2B Item ID
	offset += sectionLength - 2
	sectionLength = int(binary.LittleEndian.Uint16(zData[offset:]))
	offset += 2
	// How many reward items?
	numItems = int(binary.LittleEndian.Uint16(zData[offset+2:]))
	z.QuestNPCs = make([]int, numItems)
	for i := 0; i < numItems; i++ {
		z.QuestNPCs[i] = int(binary.LittleEndian.Uint16(zData[offset+4+(2*i):]))
	}

	// Separate out the relevant parts
	offset += sectionLength - 2
	z.Izx4a = int(zData[offset+4])
	zFlags := reverse(fmt.Sprintf("%08b", zData[offset+6]))
	rep := strings.NewReplacer("1", "Y", "0", ".")
	z.Izx4b = rep.Replace(zFlags)
	offset += 8
	trigz := make([][]byte, 0)

	// Parse actions, if there are any
	for len(zData) >= offset+4 {
		// IACT + sectionLength = 8 bytes
		sectionLength = int(binary.LittleEndian.Uint32(zData[offset+4:]))
		offset += 8
		trigz = append(trigz, zData[offset:offset+sectionLength])
		offset += sectionLength
	}

	z.ActionTriggers = make([]gosoh.ActionTrigger, len(trigz))
	for i, act := range trigz {
		// Different offset for the trigger data
		tOffset := 0
		numTriggers = int(binary.LittleEndian.Uint16(act[0:]))
		fmt.Printf("   Found %d trigger(s)", numTriggers)
		trg := gosoh.ActionTrigger{}
		trg.Conditions = make([]gosoh.TriggerCondition, numTriggers)
		tOffset += 2
		for x := 0; x < numTriggers; x++ { // Each condition is 14B
			con := gosoh.TriggerCondition{
				Condition: gosoh.TriggerConditionType(act[tOffset]),
			}
			con.Args = make([]int, 6)
			for y := 0; y < 6; y++ {
				con.Args[y] = int(binary.LittleEndian.Uint16(act[2+tOffset+(2*y):]))
			}
			trg.Conditions[x] = con
			tOffset += 14
		}
		numActions := int(binary.LittleEndian.Uint16(act[tOffset:]))
		fmt.Printf(" and %d action(s)...\n", numActions)
		trg.Actions = make([]gosoh.TriggerAction, numActions)
		tOffset += 2
		for x := 0; x < numActions; x++ {
			actn := gosoh.TriggerAction{
				Action: gosoh.TriggerActionType(act[tOffset]),
			}
			actn.Args = make([]int, 5)
			for y := 0; y < 5; y++ {
				actn.Args[y] = int(binary.LittleEndian.Uint16(act[2+tOffset+(2*y):]))
			}
			strLen := int(binary.LittleEndian.Uint16(act[tOffset+12:]))
			if strLen > 0 {
				actn.Text = string(act[tOffset+14 : tOffset+strLen+14])
				tOffset += strLen
			}
			trg.Actions[x] = actn

			tOffset += 14
		}

		z.ActionTriggers[i] = trg
	}

	return z
}

func processPuzzleData(pData []byte, numTiles int) (ret []gosoh.PuzzleInfo) {
	ret = make([]gosoh.PuzzleInfo, 0)
	offset := 0
	for len(pData) > (offset) {
		// 2 bytes of puzzle ID, plus 4 for the IPUZ header
		p := gosoh.PuzzleInfo{}
		p.Id = int(binary.LittleEndian.Uint16(pData[offset:]))
		if p.Id == 65535 { // End of puzzle section: we're out!
			return
		}
		offset += 6

		// (X - 2) bytes to hold the puzzle text
		puzzleLength := int(binary.LittleEndian.Uint16(pData[offset:]))
		puzBytes := pData[offset+2 : offset+puzzleLength]
		// TODO: interpret 0x20 as a "newline" for dialogs?

		puzTypeId := puzBytes[2]
		switch puzTypeId {
		case 0x00:
			p.Type = "ItemForItem"
		case 0x01:
			p.Type = "ItemForTask"
		case 0x02:
			p.Type = "ItemForTask2"
		case 0x03:
			p.Type = "MainQuest"
		}

		puzItemTypeId := puzBytes[6]
		switch puzItemTypeId {
		case 0x00:
			p.ItemType = "Keycard"
		case 0x01:
			p.ItemType = "Tool"
		case 0x02:
			p.ItemType = "Part"
		case 0x04:
			p.ItemType = "PlotItem"
		default:
			p.ItemType = "UNKNOWN"
			log.Fatal(fmt.Sprintf("Found unknown puzzle type: %d", puzBytes[6]))
		}

		p.NeedText, p.DoneText, p.HaveText = slurpPuzzleText(puzBytes)

		offset += puzzleLength

		// 2 bytes for the puzzle Item: either this is required to complete the thing,
		// or it's given as a reward for a different thing?
		// Might rename these, later
		p.LockItemId = int(binary.LittleEndian.Uint16(pData[offset:]))
		reward := int(binary.LittleEndian.Uint16(pData[offset+2:]))
		if reward > 0 && reward < numTiles {
			p.RewardItemId = reward
		} else { // if it's not referencing a tile, then it's probably bitflags...?
			p.RewardItemId = 0
			p.RewardFlags = reverse(fmt.Sprintf("%016b", reward))
		}

		offset += 4

		ret = append(ret, p)
	}

	return
}

func slurpPuzzleText(pb []byte) (need, done, have string) {
	textLength := 0
	textStart := 16
	need = ""
	done = ""
	have = ""

	out := ""
	ret := make([]string, 0)
	// Creep forward to the first non-zero value; that's the start of our text
	for ok := true; ok; ok = pb[textStart] == 0x00 {
		if pb[textStart] == 0x00 {
			textStart = textStart + 2
		}
	}

	for ok := true; ok; ok = textStart < (len(pb) - 4) {
		textLength = int(binary.LittleEndian.Uint16(pb[textStart:]))
		out = string(pb[textStart+2 : textStart+textLength+2])
		textStart = textStart + textLength + 2
		ret = append(ret, out)
	}
	if len(ret) == 3 {
		need = ret[0]
		done = ret[1]
		have = ret[2]
	} else if len(ret) == 2 {
		done = ret[0]
		have = ret[1]
	} else if len(ret) == 1 {
		have = ret[0]
	} else {
		log.Fatal(spew.Sprint(pb))
	}

	return need, done, have
}

func processItemList(iData []byte) (ret []gosoh.ItemInfo) {
	ret = make([]gosoh.ItemInfo, 0)
	// Each item entry is 26 bytes long
	for i := 0; i < len(iData)-26; i += 26 {
		iInfo := gosoh.ItemInfo{}
		if iInfo.Id == 65535 { // End of items section: we're out!
			return
		}
		iInfo.Id = int(binary.LittleEndian.Uint16(iData[i:]))
		// Trim the zeros from the end of the "line"
		nameLength := 0
		for j := 25; j > 1; j-- {
			if iData[i+j] != 0 {
				nameLength = j
				break
			}
		}
		iInfo.Name = string(iData[i+2 : i+1+nameLength])

		ret = append(ret, iInfo)
	}
	return
}

func processCharList(cData []byte) (ret []gosoh.CreatureInfo) {
	ret = make([]gosoh.CreatureInfo, 0)
	// Each creature entry is 84 bytes long
	for i := 0; i < len(cData)-84; i += 84 {
		cInfo := gosoh.CreatureInfo{}
		cInfo.Id = int(binary.LittleEndian.Uint16(cData[i:]))

		// Name starts at 10 and ends at the first 0
		cName := ""
		offset := 10
		for cData[offset+i] != 0x00 {
			cName += string(cData[offset+i])
			offset += 1
		}
		cInfo.Name = cName

		// These all appear to be in the same spots
		img := make(map[gosoh.CardinalDirection]int)
		img[gosoh.UpLeft] = int(binary.LittleEndian.Uint16(cData[i+36:]))
		img[gosoh.DownRight] = int(binary.LittleEndian.Uint16(cData[i+38:]))
		img[gosoh.Up] = int(binary.LittleEndian.Uint16(cData[i+40:]))
		img[gosoh.Left] = int(binary.LittleEndian.Uint16(cData[i+42:]))
		img[gosoh.DownLeft] = int(binary.LittleEndian.Uint16(cData[i+44:]))
		img[gosoh.UpRight] = int(binary.LittleEndian.Uint16(cData[i+46:]))
		img[gosoh.Right] = int(binary.LittleEndian.Uint16(cData[i+48:]))
		img[gosoh.Down] = int(binary.LittleEndian.Uint16(cData[i+50:]))
		cInfo.Images = img

		ret = append(ret, cInfo)
	}
	return
}

func reverse(str string) (result string) {
	// Given a string, return it in reverse order
	for _, v := range str {
		result = string(v) + result
	}
	return
}
