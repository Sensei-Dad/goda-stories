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
	Tileset         TiledTilesetRef       `xml:"tileset"`
	Properties      []TiledXMLProperty    `xml:"properties>property"`
	Objects         []TiledXMLObjectGroup `xml:"objectgroup"`
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
	XMLName    xml.Name
	Id         int                `xml:"id,attr"`
	Name       string             `xml:"name,attr"`
	Type       int                `xml:"type,attr"`
	X          int                `xml:"x,attr"`
	Y          int                `xml:"y,attr"`
	Width      int                `xml:"width,attr"`
	Height     int                `xml:"height,attr"`
	Rotation   int                `xml:"rotation,attr"`
	Properties []TiledXMLProperty `xml:"properties>property"`
	TileGid    int                `xml:"gid,attr"`
	Mark       TiledXMLShape      `xml:"ellipse"`
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

	zoneXML.Objects = make([]TiledXMLObjectGroup, 0)

	// Hotspots
	hotspots := TiledXMLObjectGroup{
		Id:      4,
		Name:    "Hotspots",
		Objects: make([]TiledXMLObject, 0),
	}
	zoneXML.NextLayerId++
	npcs := TiledXMLObjectGroup{
		Id:      5,
		Name:    "QuestNPCs",
		Objects: make([]TiledXMLObject, 0),
	}
	zoneXML.NextLayerId++
	actors := TiledXMLObjectGroup{
		Id:      6,
		Name:    "ZoneActors",
		Objects: make([]TiledXMLObject, 0),
	}
	zoneXML.NextLayerId++
	rewards := TiledXMLObjectGroup{
		Id:      7,
		Name:    "Rewards",
		Objects: make([]TiledXMLObject, 0),
	}
	zoneXML.NextLayerId++
	triggers := TiledXMLObjectGroup{
		Id:      8,
		Name:    "ActionTriggers",
		Objects: make([]TiledXMLObject, 0),
	}
	zoneXML.NextLayerId++

	for _, hs := range zData.Hotspots {
		spot := TiledXMLObject{
			XMLName:  xml.Name{Local: "object"},
			Id:       zoneXML.NextObjectId,
			Name:     hs.ToString(),
			Type:     int(hs.Type),
			X:        (hs.X * gosoh.TileWidth) + 2,
			Y:        (hs.Y * gosoh.TileHeight) + 2,
			Width:    gosoh.TileWidth - 4,
			Height:   gosoh.TileHeight - 4,
			Rotation: 0,
		}
		hotspots.Objects = append(hotspots.Objects, spot)
		zoneXML.NextObjectId++
	}
	for _, act := range zData.ZoneActors {
		crtr := gosoh.Creatures[act.CreatureId]
		actor := TiledXMLObject{
			XMLName:  xml.Name{Local: "object"},
			Id:       zoneXML.NextObjectId,
			Name:     crtr.Name,
			Type:     act.CreatureId,
			X:        (act.ZoneX * gosoh.TileWidth),
			Y:        (act.ZoneY * gosoh.TileHeight),
			Width:    gosoh.TileWidth,
			Height:   gosoh.TileHeight,
			Rotation: 0,
			TileGid:  gosoh.GetCreatureTNum(act.CreatureId) + 1,
		}
		actors.Objects = append(actors.Objects, actor)
		zoneXML.NextObjectId++
	}
	for i, npcId := range zData.QuestNPCs {
		npc := TiledXMLObject{
			XMLName:  xml.Name{Local: "object"},
			Id:       zoneXML.NextObjectId,
			Name:     fmt.Sprintf("NPC %d", i),
			Type:     npcId,
			X:        (i * gosoh.TileWidth),
			Y:        (-1 * gosoh.TileHeight),
			Width:    gosoh.TileWidth,
			Height:   gosoh.TileHeight,
			Rotation: 0,
			TileGid:  npcId + 1,
		}
		npcs.Objects = append(npcs.Objects, npc)
		zoneXML.NextObjectId++
	}
	for i, rewardId := range zData.RewardItems {
		rew := TiledXMLObject{
			XMLName:  xml.Name{Local: "object"},
			Id:       zoneXML.NextObjectId,
			Name:     gosoh.GetItemName(rewardId),
			Type:     rewardId,
			X:        (i * gosoh.TileWidth),
			Y:        (-3 * gosoh.TileHeight),
			Width:    gosoh.TileWidth,
			Height:   gosoh.TileHeight,
			Rotation: 0,
			TileGid:  rewardId + 1,
		}
		rewards.Objects = append(rewards.Objects, rew)
		zoneXML.NextObjectId++
	}
	for i, actTrg := range zData.ActionTriggers {
		trgr := TiledXMLObject{
			XMLName:    xml.Name{Local: "object"},
			Id:         zoneXML.NextObjectId,
			Name:       fmt.Sprintf("Trg_%d", i),
			Type:       i,
			X:          (i * gosoh.TileWidth),
			Y:          ((zData.Height + 1) * gosoh.TileHeight),
			Width:      gosoh.TileWidth,
			Height:     gosoh.TileHeight,
			Rotation:   0,
			Properties: make([]TiledXMLProperty, 0),
		}

		// Express conditions and actions as properties
		for j, c := range actTrg.Conditions {
			prop := TiledXMLProperty{
				XMLName: xml.Name{Local: "property"},
				Name:    fmt.Sprintf("IF_%03d", j),
				Value:   c.ToString(),
			}
			trgr.Properties = append(trgr.Properties, prop)
		}
		for j, a := range actTrg.Actions {
			prop := TiledXMLProperty{
				XMLName: xml.Name{Local: "property"},
				Name:    fmt.Sprintf("THEN_%03d", j),
				Value:   a.ToString(),
			}
			trgr.Properties = append(trgr.Properties, prop)
		}
		triggers.Objects = append(triggers.Objects, trgr)
		zoneXML.NextObjectId++
	}

	zoneXML.Objects = append(zoneXML.Objects, hotspots)
	zoneXML.Objects = append(zoneXML.Objects, npcs)
	zoneXML.Objects = append(zoneXML.Objects, actors)
	zoneXML.Objects = append(zoneXML.Objects, rewards)
	zoneXML.Objects = append(zoneXML.Objects, triggers)

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
	zoneXML.NextLayerId += 3

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
