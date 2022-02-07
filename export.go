package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/MasterShizzle/goda-stories/gosoh"
	"github.com/davecgh/go-spew/spew"
)

const xmlIndent string = " "
const tilesetFileName string = "yodatiles.tsx"
const tilesetImagePath string = "assets/yodatiles.png"

func saveXMLToFile(filepath string, foo interface{}) error {
	// Create the output file
	outFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	xmlstring, err := xml.MarshalIndent(foo, "", xmlIndent)
	if err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		spew.Fprintf(outFile, "%s\n", xmlstring)
	} else {
		log.Fatal(err)
	}
	fmt.Printf("[saveXMLToFile] Saved to file: %s\n", filepath)
	return err
}

type TiledXMLMap struct {
	XMLName         xml.Name
	Version         string                `xml:"version,attr"`
	TiledVersion    string                `xml:"tiledversion,attr"`
	TileOrientation string                `xml:"orientation,attr"`
	RenderOrder     string                `xml:"renderorder,attr"`
	Width           int                   `xml:"width,attr"`
	Height          int                   `xml:"height,attr"`
	TileWidth       int                   `xml:"tilewidth,attr"`
	TileHeight      int                   `xml:"tileheight,attr"`
	Infinite        int                   `xml:"infinite,attr"`
	BackgroundColor string                `xml:"backgroundcolor,attr"`
	NextLayerId     int                   `xml:"nextlayerid,attr"`
	NextObjectId    int                   `xml:"nextobjectid,attr"`
	Properties      []TiledXMLProperty    `xml:"properties>property"`
	Objects         []TiledXMLObjectGroup `xml:"objectgroup"`
	Tileset         TiledTilesetRef       `xml:"tileset"`
	Layers          []TiledXMLMapLayer    `xml:"layer"`
}

type TiledXMLProperty struct {
	XMLName xml.Name
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
}

type TiledXMLObjectGroup struct {
	XMLName xml.Name
	Id      int              `xml:"id,attr"`
	Name    string           `xml:"name,attr"`
	Objects []TiledXMLObject `xml:"object"`
}

type TiledXMLObject struct {
	XMLName  xml.Name
	Id       int           `xml:"id,attr"`
	Name     string        `xml:"name,attr"`
	Type     int           `xml:"type,attr"`
	X        int           `xml:"x,attr"`
	Y        int           `xml:"y,attr"`
	Width    int           `xml:"width,attr"`
	Height   int           `xml:"height,attr"`
	Rotation int           `xml:"rotation,attr"`
	TileGid  int           `xml:"gid,attr"`
	Mark     TiledXMLShape `xml:"ellipse"`
}

type TiledXMLShape struct {
	XMLName xml.Name
}

type TiledTilesetRef struct {
	XMLName  xml.Name
	FirstGid int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

type TiledXMLMapLayer struct {
	XMLName  xml.Name
	Id       int                 `xml:"id,attr"`
	Name     string              `xml:"name,attr"`
	Width    int                 `xml:"width,attr"`
	Height   int                 `xml:"height,attr"`
	TileData TiledCSVTileMapData `xml:"data"`
}

type TiledCSVTileMapData struct {
	XMLName  xml.Name
	Encoding string `xml:"encoding,attr"`
	Body     string `xml:",chardata"`
}

func saveTiledMaps(g *Game) {
	// Save Tileset and Zones to Tiled-compatible files
	fmt.Println("[saveTiledMaps] Stitching Zones to Tiled maps...")
	for zId, zData := range gosoh.Zones {
		mapNum := fmt.Sprintf("%03d", zId)
		mapFilePath := "assets/maps/zone_" + mapNum + ".tmx"
		saveZoneToTiledMap(mapFilePath, zData)
	}
	fmt.Printf("    %d zones extracted.\n", len(gosoh.Zones))
}

func saveZoneToTiledMap(filepath string, zData gosoh.ZoneInfo) {
	zoneXML := TiledXMLMap{
		XMLName:         xml.Name{Local: "map"},
		Version:         "1.5",
		TiledVersion:    "1.7.2",
		TileOrientation: "orthogonal",
		RenderOrder:     "right-down",
		Width:           zData.Width,
		Height:          zData.Height,
		TileWidth:       gosoh.TileWidth,
		TileHeight:      gosoh.TileHeight,
		Infinite:        0,
		BackgroundColor: "#000000",
		NextLayerId:     1,
		NextObjectId:    1,
	}
	zoneXML.Tileset = TiledTilesetRef{
		XMLName:  xml.Name{Local: "tileset"},
		FirstGid: 1,
		Source:   "../" + tilesetFileName,
	}

	// Type
	zoneXML.Properties = make([]TiledXMLProperty, 0)
	zProp := TiledXMLProperty{
		XMLName: xml.Name{Local: "property"},
		Name:    "Biome",
		Value:   zData.Biome,
	}
	zoneXML.Properties = append(zoneXML.Properties, zProp)
	zProp = TiledXMLProperty{
		XMLName: xml.Name{Local: "property"},
		Name:    "Type",
		Value:   zData.Type,
	}
	zoneXML.Properties = append(zoneXML.Properties, zProp)
	zProp = TiledXMLProperty{
		XMLName: xml.Name{Local: "property"},
		Name:    "IsOverworld",
		Value:   fmt.Sprintf("%t", zData.IsOverworld),
	}
	zoneXML.Properties = append(zoneXML.Properties, zProp)

	totalObjs := 0
	zoneXML.Objects = make([]TiledXMLObjectGroup, 0)

	// Hotspots
	hotspots := TiledXMLObjectGroup{
		Id:      4,
		Name:    "Hotspots",
		Objects: make([]TiledXMLObject, 0),
	}
	npcs := TiledXMLObjectGroup{
		Id:      5,
		Name:    "QuestNPCs",
		Objects: make([]TiledXMLObject, 0),
	}
	actors := TiledXMLObjectGroup{
		Id:      6,
		Name:    "ZoneActors",
		Objects: make([]TiledXMLObject, 0),
	}
	triggers := TiledXMLObjectGroup{
		Id:      7,
		Name:    "ActionTriggers",
		Objects: make([]TiledXMLObject, 0),
	}

	for _, hs := range zData.Hotspots {
		totalObjs++
		spot := TiledXMLObject{
			XMLName:  xml.Name{Local: "object"},
			Id:       totalObjs,
			Name:     hs.ToString(),
			Type:     int(hs.Type),
			X:        (hs.X * gosoh.TileWidth),
			Y:        (hs.Y * gosoh.TileHeight),
			Width:    gosoh.TileWidth - 2,
			Height:   gosoh.TileHeight - 2,
			Rotation: 0,
		}
		hotspots.Objects = append(hotspots.Objects, spot)
	}

	zoneXML.Objects = append(zoneXML.Objects, hotspots)
	zoneXML.Objects = append(zoneXML.Objects, npcs)
	zoneXML.Objects = append(zoneXML.Objects, actors)
	zoneXML.Objects = append(zoneXML.Objects, triggers)

	// TODO:
	// - Actors
	// - Action Triggers
	// - Quest NPCs

	zoneXML.Layers = make([]TiledXMLMapLayer, 0)

	// Add Terrain tiles
	terrainLayer := TiledXMLMapLayer{
		XMLName: xml.Name{Local: "layer"},
		Id:      1,
		Name:    "Terrain",
		Width:   zData.Width,
		Height:  zData.Height,
	}
	terrainLayer.TileData = TiledCSVTileMapData{
		XMLName:  xml.Name{Local: "data"},
		Encoding: "csv",
		Body:     getLayerCSV(zData.TileMaps.Terrain),
	}

	// Walls
	wallsLayer := TiledXMLMapLayer{
		XMLName: xml.Name{Local: "layer"},
		Id:      2,
		Name:    "Walls",
		Width:   zData.Width,
		Height:  zData.Height,
	}
	wallsLayer.TileData = TiledCSVTileMapData{
		XMLName:  xml.Name{Local: "data"},
		Encoding: "csv",
		Body:     getLayerCSV(zData.TileMaps.Walls),
	}

	// Overlay
	overlayLayer := TiledXMLMapLayer{
		XMLName: xml.Name{Local: "layer"},
		Id:      3,
		Name:    "Overlay",
		Width:   zData.Width,
		Height:  zData.Height,
	}
	overlayLayer.TileData = TiledCSVTileMapData{
		XMLName:  xml.Name{Local: "data"},
		Encoding: "csv",
		Body:     getLayerCSV(zData.TileMaps.Overlay),
	}

	zoneXML.Layers = append(zoneXML.Layers, terrainLayer)
	zoneXML.Layers = append(zoneXML.Layers, wallsLayer)
	zoneXML.Layers = append(zoneXML.Layers, overlayLayer)

	err := saveXMLToFile(filepath, zoneXML)
	if err != nil {
		log.Fatal(err)
	}
}

func getLayerCSV(layerData []int) (ret string) {
	ret = ""

	for _, tNum := range layerData {
		if tNum == 65535 {
			ret += fmt.Sprintf("%d,", 0)
		} else {
			ret += fmt.Sprintf("%d,", tNum+1)
		}
	}

	return
}
