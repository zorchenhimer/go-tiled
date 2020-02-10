package tiled

import (
	//"encoding/xml"
)

type Tileset struct {
	FirstGid uint
	Source string
	Name string

	TileWidth uint
	TileHeight uint
	TileCount uint
	Columns uint
}

type xmlTileset struct {
	FirstId int    `xml:"firstgid,attr"`
	Source  string `xml:"source,attr"`
}
