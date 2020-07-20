package tiled

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Tileset struct {
	FirstGid uint
	Source   string
	Name     string

	TileWidth  uint
	TileHeight uint
	TileCount  uint
	Columns    uint

	Tiles []Tile
}

type Tile struct {
	Id         uint
	Width      uint
	Height     uint
	Source     string
	Properties CustomProperties
	Image      string // source string
}

type xmlTileset struct {
	Version      string `xml:"version,attr"`
	TiledVersion string `xml:"tiledversion,attr"`
	Name         string `xml:"name,attr"`
	TileWidth    uint   `xml:"tilewidth,attr"`
	TileHeight   uint   `xml:"tileheight,attr"`
	TileCount    uint   `xml"tilecount,attr"`
	Columns      uint   `xml:"columns:attr"`

	Grid  *xmlTsGrid
	Tiles []xmlTsTile `xml:"tile"`
	Image *xmlTsImage `xml:"image"`
}

type xmlTsImage struct {
	Source string `xml:"source,attr"`
	Trans  string `xml:"trans,attr"`
	Width  uint   `xml:"width,attr"`
	Height uint   `xml:"height,attr"`
}

type xmlTsGrid struct {
	Orientation string `xml:"orientation,attr"`
	Width       uint   `xml:"width,attr"`
	Height      uint   `xml:"height,attr"`
}

type xmlTsTile struct {
	Id         uint            `xml:"id,attr"`
	Properties xmlPropertyList `xml:"properties>property"`
	Image      *xmlTsImage     `xml:"image"`
}

func LoadTileset(filename string) (*Tileset, error) {
	rawXml, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading XML file: %v", err)
	}

	return LoadTilesetRaw(rawXml)
}

func LoadTilesetRaw(rawXml []byte) (*Tileset, error) {
	var td xmlTileset
	err := xml.Unmarshal(rawXml, &td)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal XML: %v", err)
	}

	ts := &Tileset{
		TileWidth:  td.TileWidth,
		TileHeight: td.TileHeight,
		TileCount:  td.TileCount,
		Columns:    td.Columns,

		Tiles: []Tile{},
	}

	for _, tile := range td.Tiles {
		//fmt.Printf("[%d] %v\n", i, tile)
		t := Tile{
			Id: tile.Id,
			//Properties: CustomProperties{},
		}

		props, err := tile.Properties.CustomProps()
		if err != nil {
			return nil, fmt.Errorf("Unable to parse properties: %v", err)
		}
		t.Properties = props

		if tile.Image != nil {
			t.Width = tile.Image.Width
			t.Height = tile.Image.Height
			t.Image = tile.Image.Source
		}

		ts.Tiles = append(ts.Tiles, t)
	}

	return ts, nil
}
