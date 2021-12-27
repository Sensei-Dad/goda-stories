# Go-da Stories

Port of Yoda Stories to Go, based on [an article](https://www.gamedeveloper.com/programming/reverse-engineering-the-binary-data-format-for-star-wars-yoda-stories) about the datafile format

To get started, simply copy the YODESK.DTA file from your Yoda Stories installation into `data/`.


## Section notes

### TILE
On run, will output tile data and export pngs to `assets/tiles`.

> TODO: Create the tiles dir if it doesn't exist

### ZONE

Section layout:
```
0:2 	2 B: ID
2:6 	4 B: "IZON"
6:10	4 B: (unknown)
10:12	2 B: map width (W)
12:14	2 B: map height (H)
14		1 B: map flags (TODO)
15:20	5 B: unused (same values for every map)
20		1 B: planet
			0x01 = desert
			0x02 = snow
			0x03 = forest
			0x05 = swamp
21		1 B: unused (same values for every map)
22:X	(W * H) * 6 B: map data
			Each "cell" in the map has 3 tiles, each denoted by a Uint16
			Tile values correspond to the tile number
		2 B: object info entry count (X)
		(X * 12) B: object info data
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

