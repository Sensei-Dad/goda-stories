# Extracting the Yoda Stories Datafile

## Section notes

### TILE
On run, will output tile data and export pngs to `assets/tiles`.

> TODO:
>   * Create the tiles dir if it doesn't exist
>   * Embed these things, instead of creating the actual PNGs

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
Each IZON entry in the ZONE section describes a map, or "screen", or whatever. In the original game, it describes all the stuff that happens before any kind of screen or room transitions.

Between the hotspots of the zone, all the possible enemies on it, and all the action triggers for cutscenes, quest progress, etc... there's a LOT of information embedded into each IZON entry.

* Tilemaps for the Terrain layer, Walls, and Overlay
* Tile-based triggers (Hotspots)
* IZAX handles Zone Actors
  * 4B IZAX header, 2B section length
  * 2B unused, 2B unknown, then 2B to count the number of 44-byte creature entries
* IZX2 is a list of possible items that this Zone can "produce", e.g. what's allowed to drop here
  * 2 bytes unused
  * 2 bytes: Count X of how many item entries to follow
  * X * 2 bytes (Uint16): Tile ID / Item ID of the item in question
* IZX3 is the list of NPC character tiles that can be used on this quest
  * Very similar to the previous: A Uint16 to count by, then a list of tile IDs.
* IZX4 is some kind of eldritch nonsense which I don't understand. I'm putting it in the "stuff I might fathom if I ever learned C" pile, where it can be ignored properly.
* IACT is a lot of entries for action triggers, and moving spawned objects around the map; these deserve their own doc.

### PUZ2
These are entries that assign each item a bit of text that makes sense for it in context. Each text block has one or more of the following text strings, which take the general form of HaveText, NeedText, and DoneText. By looking up the various strings for two different items, you can cobble together an NPC's entire dialogue.

* 2 bytes - Puzzle ID
* 4 bytes - IPUZ header
* 2 bytes - length of the puzzle text
* X bytes of puzzle data:
    * 2 bytes unused
    * 2 bytes Puzzle Type, which determines the "shape" of the text block:
        * NeedText - "Hey, I need a ___" said when speaking to an NPC who needs this item
        * DoneText - "Thanks for the ___" said after the item-related puzzle has been finished
        * HaveText - "If you help, you can have this ___" said by an NPC who has this item and needs something else
        * There are 15 "Main" missions: you always get these from Yoda, and finishing these will win the game
            * The NeedText is unused, but that also means a lot of funny blurbs from the devs live here
            * DoneText is a Quest synopsis / mission title
            * HaveText is Yoda's conversation you get at the start of a new game, after finding him
    * 2 bytes unused
    * 2 bytes for the type of Item it is (for what shows on the Locator screens, etc.): Keycard / Tool / Part / Valuable
    * 2 bytes unused
    * X bytes of Puzzle text strings, in 2-byte length-denoted blocks (0x20 is a "newline")
* 2 bytes for a tile ID
* 2 more bytes for another tile ID? Though not always? Maybe these are flags.

### SNDS and TNAM
These are simple lists of strings, which refer to sound files and tile names respectively.

## TODO
```
CHWP
CAUX
```
