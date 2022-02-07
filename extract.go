package main

// Functions to extract and process the data from .DTA resources
import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"strings"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/davecgh/go-spew/spew"
	"github.com/ghostiam/binstruct"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Palette data extracted from the de-compiled Yoda Stories binary
var PaletteData = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x8B, 0x00, 0xC3, 0xCF, 0x4B, 0x00,
	0x8B, 0xA3, 0x1B, 0x00, 0x57, 0x77, 0x00, 0x00, 0x8B, 0xA3, 0x1B, 0x00, 0xC3, 0xCF, 0x4B, 0x00,
	0xFB, 0xFB, 0xFB, 0x00, 0xEB, 0xE7, 0xE7, 0x00, 0xDB, 0xD3, 0xD3, 0x00, 0xCB, 0xC3, 0xC3, 0x00,
	0xBB, 0xB3, 0xB3, 0x00, 0xAB, 0xA3, 0xA3, 0x00, 0x9B, 0x8F, 0x8F, 0x00, 0x8B, 0x7F, 0x7F, 0x00,
	0x7B, 0x6F, 0x6F, 0x00, 0x67, 0x5B, 0x5B, 0x00, 0x57, 0x4B, 0x4B, 0x00, 0x47, 0x3B, 0x3B, 0x00,
	0x33, 0x2B, 0x2B, 0x00, 0x23, 0x1B, 0x1B, 0x00, 0x13, 0x0F, 0x0F, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0xC7, 0x43, 0x00, 0x00, 0xB7, 0x43, 0x00, 0x00, 0xAB, 0x3F, 0x00, 0x00, 0x9F, 0x3F, 0x00,
	0x00, 0x93, 0x3F, 0x00, 0x00, 0x87, 0x3B, 0x00, 0x00, 0x7B, 0x37, 0x00, 0x00, 0x6F, 0x33, 0x00,
	0x00, 0x63, 0x33, 0x00, 0x00, 0x53, 0x2B, 0x00, 0x00, 0x47, 0x27, 0x00, 0x00, 0x3B, 0x23, 0x00,
	0x00, 0x2F, 0x1B, 0x00, 0x00, 0x23, 0x13, 0x00, 0x00, 0x17, 0x0F, 0x00, 0x00, 0x0B, 0x07, 0x00,
	0x4B, 0x7B, 0xBB, 0x00, 0x43, 0x73, 0xB3, 0x00, 0x43, 0x6B, 0xAB, 0x00, 0x3B, 0x63, 0xA3, 0x00,
	0x3B, 0x63, 0x9B, 0x00, 0x33, 0x5B, 0x93, 0x00, 0x33, 0x5B, 0x8B, 0x00, 0x2B, 0x53, 0x83, 0x00,
	0x2B, 0x4B, 0x73, 0x00, 0x23, 0x4B, 0x6B, 0x00, 0x23, 0x43, 0x5F, 0x00, 0x1B, 0x3B, 0x53, 0x00,
	0x1B, 0x37, 0x47, 0x00, 0x1B, 0x33, 0x43, 0x00, 0x13, 0x2B, 0x3B, 0x00, 0x0B, 0x23, 0x2B, 0x00,
	0xD7, 0xFF, 0xFF, 0x00, 0xBB, 0xEF, 0xEF, 0x00, 0xA3, 0xDF, 0xDF, 0x00, 0x8B, 0xCF, 0xCF, 0x00,
	0x77, 0xC3, 0xC3, 0x00, 0x63, 0xB3, 0xB3, 0x00, 0x53, 0xA3, 0xA3, 0x00, 0x43, 0x93, 0x93, 0x00,
	0x33, 0x87, 0x87, 0x00, 0x27, 0x77, 0x77, 0x00, 0x1B, 0x67, 0x67, 0x00, 0x13, 0x5B, 0x5B, 0x00,
	0x0B, 0x4B, 0x4B, 0x00, 0x07, 0x3B, 0x3B, 0x00, 0x00, 0x2B, 0x2B, 0x00, 0x00, 0x1F, 0x1F, 0x00,
	0xDB, 0xEB, 0xFB, 0x00, 0xD3, 0xE3, 0xFB, 0x00, 0xC3, 0xDB, 0xFB, 0x00, 0xBB, 0xD3, 0xFB, 0x00,
	0xB3, 0xCB, 0xFB, 0x00, 0xA3, 0xC3, 0xFB, 0x00, 0x9B, 0xBB, 0xFB, 0x00, 0x8F, 0xB7, 0xFB, 0x00,
	0x83, 0xB3, 0xF7, 0x00, 0x73, 0xA7, 0xFB, 0x00, 0x63, 0x9B, 0xFB, 0x00, 0x5B, 0x93, 0xF3, 0x00,
	0x5B, 0x8B, 0xEB, 0x00, 0x53, 0x8B, 0xDB, 0x00, 0x53, 0x83, 0xD3, 0x00, 0x4B, 0x7B, 0xCB, 0x00,
	0x9B, 0xC7, 0xFF, 0x00, 0x8F, 0xB7, 0xF7, 0x00, 0x87, 0xB3, 0xEF, 0x00, 0x7F, 0xA7, 0xF3, 0x00,
	0x73, 0x9F, 0xEF, 0x00, 0x53, 0x83, 0xCF, 0x00, 0x3B, 0x6B, 0xB3, 0x00, 0x2F, 0x5B, 0xA3, 0x00,
	0x23, 0x4F, 0x93, 0x00, 0x1B, 0x43, 0x83, 0x00, 0x13, 0x3B, 0x77, 0x00, 0x0B, 0x2F, 0x67, 0x00,
	0x07, 0x27, 0x57, 0x00, 0x00, 0x1B, 0x47, 0x00, 0x00, 0x13, 0x37, 0x00, 0x00, 0x0F, 0x2B, 0x00,
	0xFB, 0xFB, 0xE7, 0x00, 0xF3, 0xF3, 0xD3, 0x00, 0xEB, 0xE7, 0xC7, 0x00, 0xE3, 0xDF, 0xB7, 0x00,
	0xDB, 0xD7, 0xA7, 0x00, 0xD3, 0xCF, 0x97, 0x00, 0xCB, 0xC7, 0x8B, 0x00, 0xC3, 0xBB, 0x7F, 0x00,
	0xBB, 0xB3, 0x73, 0x00, 0xAF, 0xA7, 0x63, 0x00, 0x9B, 0x93, 0x47, 0x00, 0x87, 0x7B, 0x33, 0x00,
	0x6F, 0x67, 0x1F, 0x00, 0x5B, 0x53, 0x0F, 0x00, 0x47, 0x43, 0x00, 0x00, 0x37, 0x33, 0x00, 0x00,
	0xFF, 0xF7, 0xF7, 0x00, 0xEF, 0xDF, 0xDF, 0x00, 0xDF, 0xC7, 0xC7, 0x00, 0xCF, 0xB3, 0xB3, 0x00,
	0xBF, 0x9F, 0x9F, 0x00, 0xB3, 0x8B, 0x8B, 0x00, 0xA3, 0x7B, 0x7B, 0x00, 0x93, 0x6B, 0x6B, 0x00,
	0x83, 0x57, 0x57, 0x00, 0x73, 0x4B, 0x4B, 0x00, 0x67, 0x3B, 0x3B, 0x00, 0x57, 0x2F, 0x2F, 0x00,
	0x47, 0x27, 0x27, 0x00, 0x37, 0x1B, 0x1B, 0x00, 0x27, 0x13, 0x13, 0x00, 0x1B, 0x0B, 0x0B, 0x00,
	0xF7, 0xB3, 0x37, 0x00, 0xE7, 0x93, 0x07, 0x00, 0xFB, 0x53, 0x0B, 0x00, 0xFB, 0x00, 0x00, 0x00,
	0xCB, 0x00, 0x00, 0x00, 0x9F, 0x00, 0x00, 0x00, 0x6F, 0x00, 0x00, 0x00, 0x43, 0x00, 0x00, 0x00,
	0xBF, 0xBB, 0xFB, 0x00, 0x8F, 0x8B, 0xFB, 0x00, 0x5F, 0x5B, 0xFB, 0x00, 0x93, 0xBB, 0xFF, 0x00,
	0x5F, 0x97, 0xF7, 0x00, 0x3B, 0x7B, 0xEF, 0x00, 0x23, 0x63, 0xC3, 0x00, 0x13, 0x53, 0xB3, 0x00,
	0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0xEF, 0x00, 0x00, 0x00, 0xE3, 0x00, 0x00, 0x00, 0xD3, 0x00,
	0x00, 0x00, 0xC3, 0x00, 0x00, 0x00, 0xB7, 0x00, 0x00, 0x00, 0xA7, 0x00, 0x00, 0x00, 0x9B, 0x00,
	0x00, 0x00, 0x8B, 0x00, 0x00, 0x00, 0x7F, 0x00, 0x00, 0x00, 0x6F, 0x00, 0x00, 0x00, 0x63, 0x00,
	0x00, 0x00, 0x53, 0x00, 0x00, 0x00, 0x47, 0x00, 0x00, 0x00, 0x37, 0x00, 0x00, 0x00, 0x2B, 0x00,
	0x00, 0xFF, 0xFF, 0x00, 0x00, 0xE3, 0xF7, 0x00, 0x00, 0xCF, 0xF3, 0x00, 0x00, 0xB7, 0xEF, 0x00,
	0x00, 0xA3, 0xEB, 0x00, 0x00, 0x8B, 0xE7, 0x00, 0x00, 0x77, 0xDF, 0x00, 0x00, 0x63, 0xDB, 0x00,
	0x00, 0x4F, 0xD7, 0x00, 0x00, 0x3F, 0xD3, 0x00, 0x00, 0x2F, 0xCF, 0x00, 0x97, 0xFF, 0xFF, 0x00,
	0x83, 0xDF, 0xEF, 0x00, 0x73, 0xC3, 0xDF, 0x00, 0x5F, 0xA7, 0xCF, 0x00, 0x53, 0x8B, 0xC3, 0x00,
	0x2B, 0x2B, 0x00, 0x00, 0x23, 0x23, 0x00, 0x00, 0x1B, 0x1B, 0x00, 0x00, 0x13, 0x13, 0x00, 0x00,
	0xFF, 0x0B, 0x00, 0x00, 0xFF, 0x00, 0x4B, 0x00, 0xFF, 0x00, 0xA3, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	0x00, 0xFF, 0x00, 0x00, 0x00, 0x4B, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0x33, 0x2F, 0x00,
	0x00, 0x00, 0xFF, 0x00, 0x00, 0x1F, 0x97, 0x00, 0xDF, 0x00, 0xFF, 0x00, 0x73, 0x00, 0x77, 0x00,
	0x6B, 0x7B, 0xC3, 0x00, 0x57, 0x57, 0xAB, 0x00, 0x57, 0x47, 0x93, 0x00, 0x53, 0x37, 0x7F, 0x00,
	0x4F, 0x27, 0x67, 0x00, 0x47, 0x1B, 0x4F, 0x00, 0x3B, 0x13, 0x3B, 0x00, 0x27, 0x77, 0x77, 0x00,
	0x23, 0x73, 0x73, 0x00, 0x1F, 0x6F, 0x6F, 0x00, 0x1B, 0x6B, 0x6B, 0x00, 0x1B, 0x67, 0x67, 0x00,
	0x1B, 0x6B, 0x6B, 0x00, 0x1F, 0x6F, 0x6F, 0x00, 0x23, 0x73, 0x73, 0x00, 0x27, 0x77, 0x77, 0x00,
	0xFF, 0xFF, 0xEF, 0x00, 0xF7, 0xF7, 0xDB, 0x00, 0xF3, 0xEF, 0xCB, 0x00, 0xEF, 0xEB, 0xBB, 0x00,
	0xF3, 0xEF, 0xCB, 0x00, 0xE7, 0x93, 0x07, 0x00, 0xE7, 0x97, 0x0F, 0x00, 0xEB, 0x9F, 0x17, 0x00,
	0xEF, 0xA3, 0x23, 0x00, 0xF3, 0xAB, 0x2B, 0x00, 0xF7, 0xB3, 0x37, 0x00, 0xEF, 0xA7, 0x27, 0x00,
	0xEB, 0x9F, 0x1B, 0x00, 0xE7, 0x97, 0x0F, 0x00, 0x0B, 0xCB, 0xFB, 0x00, 0x0B, 0xA3, 0xFB, 0x00,
	0x0B, 0x73, 0xFB, 0x00, 0x0B, 0x4B, 0xFB, 0x00, 0x0B, 0x23, 0xFB, 0x00, 0x0B, 0x73, 0xFB, 0x00,
	0x00, 0x13, 0x93, 0x00, 0x00, 0x0B, 0xD3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00,
}

func processYodaFile(fileName string) (*ebiten.Image, []gosoh.TileInfo, []gosoh.ZoneInfo, []gosoh.ItemInfo, []gosoh.PuzzleInfo, []gosoh.CreatureInfo, []string) {
	yodaFilePath := "data/" + fileName
	tileImageBytes := make([][]byte, 0)
	outTiles := make([]gosoh.TileInfo, 0)
	outZones := make([]gosoh.ZoneInfo, 0)
	outItems := make([]gosoh.ItemInfo, 0)
	outPuzzles := make([]gosoh.PuzzleInfo, 0)
	outCreatures := make([]gosoh.CreatureInfo, 0)
	outSounds := make([]string, 0)

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
		case "STUP", "CHWP", "CAUX":
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

			// Extract tile bits into images
			for i := 0; i < numTiles; i++ {
				flags, _ := reader.ReadUint32()
				tData := processTileData(i, flags)
				outTiles = append(outTiles, tData)

				_, tileBytes, _ := reader.ReadBytes(0x400)

				tileImageBytes = append(tileImageBytes, tileBytes)
			}
			fmt.Printf("    Extracted %d tile images\n", numTiles)
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
		case "SNDS":
			sectionLength, _ := reader.ReadUint32()
			_, soundData, err := reader.ReadBytes(int(sectionLength))
			outSounds = processSoundList(soundData)
			if err != nil {
				fmt.Printf("Error reading section %s\n", s)
				log.Fatal(err)
			}
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

	// Draw tiles to a tileset image, and save
	tileRows := int(numTiles/gosoh.TilesetColumns) + 1
	tImg := image.NewNRGBA(image.Rect(0, 0, gosoh.TilesetColumns*gosoh.TileWidth, tileRows*gosoh.TileHeight))
	for tNum, t := range tileImageBytes {
		tileX, tileY := gosoh.GetTileCoords(tNum)
		for j := 0; j < len(t); j++ {
			pixel := int(t[j])
			if pixel != 0 {
				rVal := PaletteData[pixel*4+2]
				gVal := PaletteData[pixel*4+1]
				bVal := PaletteData[pixel*4+0]
				clr := color.NRGBA{R: rVal, G: gVal, B: bVal, A: 255}
				tImg.Set((j%gosoh.TileWidth)+tileX, (j/gosoh.TileHeight)+tileY, clr)
			} else {
				tImg.Set((j%gosoh.TileWidth)+tileX, (j/gosoh.TileHeight)+tileY, color.Transparent)
			}
		}
	}
	f, _ := os.Create(tilesetImagePath)
	png.Encode(f, tImg)
	fmt.Printf("[%s] Saved tileset image: %s\n", fileName, tilesetImagePath)

	tilesetImage, _, err := ebitenutil.NewImageFromFile(tilesetImagePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[%s] Processed data file.\n", yodaFile)

	return tilesetImage, outTiles, outZones, outItems, outPuzzles, outCreatures, outSounds
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
		t.Type = "Wall"
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
	// Locator minimap tiles: #817-837
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
	z.TileMaps.Walls = make([]int, z.Width*z.Height)
	z.TileMaps.Overlay = make([]int, z.Width*z.Height)
	for j := 0; j < (z.Width * z.Height); j++ {
		z.TileMaps.Terrain[j] = int(binary.LittleEndian.Uint16(zData[6*j+22:]))
		z.TileMaps.Walls[j] = int(binary.LittleEndian.Uint16(zData[6*j+24:]))
		z.TileMaps.Overlay[j] = int(binary.LittleEndian.Uint16(zData[6*j+26:]))
	}

	offset := (6 * z.Width * z.Height) + 22

	// Parse entries for hotspots
	numTriggers := int(binary.LittleEndian.Uint16(zData[offset:]))
	if numTriggers > 0 {
		z.Hotspots = make([]gosoh.ZoneHotspot, numTriggers)
		for k := 0; k < numTriggers; k++ {
			z.Hotspots[k].Type = gosoh.TriggerHotspotType(int(binary.LittleEndian.Uint16(zData[offset+2:])))
			z.Hotspots[k].Id = k
			z.Hotspots[k].X = int(binary.LittleEndian.Uint16(zData[offset+6:]))
			z.Hotspots[k].Y = int(binary.LittleEndian.Uint16(zData[offset+8:]))
			z.Hotspots[k].Arg = int(binary.LittleEndian.Uint16(zData[offset+12:]))
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
				Index:      i,
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

func processSoundList(sData []byte) (ret []string) {
	ret = make([]string, 0)
	offset := 2
	for offset < len(sData) {
		strLen := int(binary.LittleEndian.Uint16(sData[offset:]))
		ret = append(ret, string(sData[offset+2:offset+strLen+1]))
		offset = offset + strLen + 2
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
