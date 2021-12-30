# Go-da Stories

Port of Yoda Stories to Go, based on [an article](https://www.gamedeveloper.com/programming/reverse-engineering-the-binary-data-format-for-star-wars-yoda-stories) about the datafile format

To get started, simply copy the YODESK.DTA file from your Yoda Stories installation into `data/`.


## Section notes

### TILE
On run, will output tile data and export pngs to `assets/tiles`.

> TODO: Create the tiles dir if it doesn't exist

#### Demystifying the Tile flag bits
Based on Zach's incredible work, I could make some easy eliminations right off the bat when designing the `TileInfo` struct:

* Without any real work, it was easy to verify that bits 9-15 are unused (*i.e.* the same value for all tiles), so we get to ignore those "for free".
* He refers to the first group as "Type" bits. I checked and permuted through the first nine bits of all the tiles, finding only 10 unique types (Huzzahs for all!). Checking examples of each (and once again cheating off of Zach's answers), we can make some assumptions and try to categorize all 2123 tiles in the game:
  * `010000000` - (441 tiles) Terrain, walkable. Drawn on the bottom layer, can be walked-upon by the player / mobs.
  * `101000000` - (416 tiles) Terrain, non-walkable. Closed doors, chests, furniture, etc. Most all of these tiles have transparent bits when I checked samples, so it's safe to assume this means they should be drawn on the middle layer. Uses none of the later bits.
  * `001000000` - (472 tiles) Terrain, non-walkable. Unlike the above non-walkable category, none of these tiles have transparent pixels to them, so we can assume they're also usually drawn on the *bottom* layer... though if the player can't walk over them, it doesn't make much difference.
  * `101100000` - (18 tiles) Pushable object, pretty easily verified since there are so few examples. These are the things you can move around by holding Shift.
  * `100010000` - (333 tiles) Terrain, walkable. Drawn on the top layer, after the player and other objects (arches, the top bits of tall buildings the player can walk behind, etc.).
  * `000010000` - (5 tiles) Terrain, only these 5 appear to be "full" bottom-layer tiles for some reason. These look like water / jungle tiles, mostly, so possibly these five show the player half-submerged or something...?
    * TODO: Once we figure out the scripting, look up the maps where these 5 tiles are actually used
  * `100000001` - (246 tiles) These are the various characters, monsters, etc. in the game (we'll say Creatures, to borrow some D&D terminology).
  * `100000010` - (167 tiles) Game objects that can be picked up, and get shown in the inventory (alongside weapons)
  * `100000100` - (10 tiles) Weapon object, or The Force. Stuff that Luke can equip in his "weapon" slot and use with the left mouse button.
  * `000001000` - (15 tiles) These are used for the Locator mini-map, but oddly it doesn't include ALL the tiles for the minimap. Luckily we can find the odd ones grouped together in-between the ones with this bitmask (starting at tile 817), so it shouldn't be too difficult to program around.

...and there we go. Scrunch a couple of these together (weapons and items are both "show on map => pick up" kinds of tiles, etc.), we can start verifying what the last groups of bits are used for.

#### Final tile sort
From the above, we can start deciding on tile categories that could be useful for programming, and decipher all the remaining bits:
* Terrain Tiles, always drawn on the bottom
  * 142 of these use bit 16. Zach refers to it as a "door" bit in terms of terrain, but on inspection this flag also includes tiles that depict switches, staircases, plain squares of terrain, etc... Since all these are walkable, I'm thinking that this bit indicates any tile that triggers `<some generic event>` when stepped-upon, and door functionality (*i.e.* a "move the player elsewhere" event) happens to be included in those functions: A handy thing to keep in mind while we're looking at the action scripts.
* Object Tiles, always drawn on the middle layer
* Block Tiles, which can be pushed / pulled
* Overlay Tiles, drawn on top of the rest
  * It's worth noting (*i.e.* it saves me some work) that none of the Object, Block, or Overlay tiles use any of the other bits.
* Creature Tiles, which move around and do stuff.
* Item Tiles, including keycards, key items, and weapons. We're doing this inside an ECS, so we can just make components corresponding to Weapons, Keycards, ConsumableItemsThatHealThePlayer, etc.
* Minimap Tiles, notable for being the only ones which are tiled on a different "grid" than the others: the game worlds / Locator maps are 10 tiles across with a border, when the UI usually displays 9 tiles per axis.

One nice thing is that (for now) we don't *really* need to worry overmuch about which order things get drawn in, since the layer(s) for each tile are also embedded in the ZONE data. Perhaps this will change if we start wanting to create our own maps.

### ZONE

Section layout:
| index 		| len 			| desc |
|---------------|---------------|------|
| 0:2   		| 2 B 			| ID |
| 2:6 			| 4 B			| "IZON" |
| 6:10			| 4 B			| (unknown) |
| 10:12			| 2 B			| map width (W) |
| 12:14			| 2 B			| map height (H) |
| 14			| 1 B			| map flags (TODO) |
| 15:20			| 5 B			| unused (same values for every map) |
| 20			| 1 B			| planet type (desert / snow / forest / swamp) |
| 21			| 1 B 			| unused (same values for every map) |
| 22:X			| (W * H) * 6 B	| map data: Each "cell" in the map has 3 tiles, each denoted by a Uint16; tile values correspond to the tile number |
| ??			| 2 B			| object info entry count (X) |
| ?? 			| (X * 12) B	| object triggers |

TODO
```
		4 B: "IZAX"
		2 B: length (X)
		(X - 6) B: IZAX data
		4 B: "IZX2"
		2 B: length (X)
		(X - 6) B: IZX2data
		4 B: "IZX3"
		2 B: length (X)
		(X - 6) B: IZX3 data
		4 B: "IZX4"
		8 B: IZX4 data
		4 B: "IACT"
		4 B: length (X)
		X B: action data
		...
		4 B: "IACT"
		4 B: length (X)
		X B: action data
```

