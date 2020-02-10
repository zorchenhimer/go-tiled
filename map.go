package tiled

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Map struct {
	Properties MapProperties
	CustomProperties CustomProperties
	Tilesets []Tileset
	Layers []Layer

	version string
	tiledVersion string
}

func (m Map) Version() string {
	return m.version
}

func (m Map) TiledVersion() string {
	return m.tiledVersion
}

type xmlMap struct {
	XMLName    string       `xml:"map"`
	Layers     []xmlLayer   `xml:"layer"`
	Tilesets   []xmlTileset `xml:"tileset"`
	SourceFile string       `xml:"-"`

	Version string `xml:"version,attr"`
	TiledVersion string `xml:"tiledversion,attr"`
}

func LoadMap(filename string) (*Map, error) {
	rawXml, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading XML file: %v", err)
	}

	return LoadMapRaw(rawXml)
}

func LoadMapRaw(rawXml []byte) (*Map, error) {
	var md xmlMap
	err := xml.Unmarshal(rawXml, &md)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling XML: %v", err)
	}

	layers, err := decodeLayers(md.Layers)
	if err != nil {
		return nil, err
	}

	m := &Map{
		version: md.Version,
		tiledVersion: md.TiledVersion,
		Layers: layers,
	}

	return m, nil
}

func (m Map) GetLayerByName(name string) []Layer {
	ret := []Layer{}
	for _, layer := range m.Layers {
		if layer.Name == name {
			ret = append(ret, layer)
		}
	}
	return ret
}

func (m Map) GetLayer(id int) (Layer, error) {
	for _, l := range m.Layers {
		if l.Id == id {
			return l, nil
		}
	}
	return Layer{}, fmt.Errorf("No such layer")
}

