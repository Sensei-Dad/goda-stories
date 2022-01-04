package main

import (
	"fmt"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const tileWidth, tileHeight = 32, 32 // Tile width and height, in pixels
// const uiPadding = 5                  // Padding between UI elements, in pixels
const ViewportWidth, ViewportHeight = 12, 12

type Game struct {
	CurrentScreen int
	World         GameWorld
	ECSManager    *ecs.Manager
	ECSTags       map[string]ecs.Tag
	tick          int64
}

var blankTile = ebiten.NewImage(tileWidth, tileHeight)
var tiles = []*ebiten.Image{}

func NewGame(tInfo []TileInfo, zones []ZoneInfo) *Game {
	// TODO: Distinguish between "init game" and "new game"
	tiles = LoadAllTiles(tInfo)
	g := &Game{}
	g.World = NewWorld()

	// ECS!
	mgr, tags := g.InitializeWorld()
	g.ECSManager = mgr
	g.ECSTags = tags

	g.tick = 0

	return g
}

func (g *Game) Update() error {
	// g.tick++
	g.ProcessInput()
	// TODO: Handle AI, randomly move critters around, etc.
	// ProcessCreatures(g)
	g.ProcessMovement()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each layer, starting at the bottom:
	//     = Overlay =
	//  ** Renderables **
	//  ==   Objects   ==
	// ===   Terrain   ===
	layers := g.World.GetLayers()
	layers.DrawLayer(Terrain, screen)
	layers.DrawLayer(Objects, screen)
	ProcessRenderables(g, layers, screen)
	layers.DrawLayer(Overlay, screen)

	// Show FPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(w, h int) (int, int) {
	// for now, return the map with nothing else around it
	return (ViewportWidth * tileWidth), (ViewportHeight * tileHeight)
}
