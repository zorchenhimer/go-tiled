package tiled

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type Orientation int

const (
	Orient_Orthogonal Orientation = iota
	Orient_Isometric
	Orient_IsometricStaggered
	Orient_HexagonalStaggered
)

type StaggerAxis int

const (
	Stagger_X StaggerAxis = iota
	Stagger_Y
)

type StaggerIndex int

const (
	Stagger_Odd StaggerIndex = iota
	Stagger_Even
)

type LayerFormat int

const (
	LayerFormat_Base64Uncompressed LayerFormat = iota
	LayerFormat_Base64Gzip
	LayerFormat_Base64Zlib
	LayerFormat_Base64Zstandard
	LayerFormat_CSV
)

type RenderOrder int

const (
	RightDown RenderOrder = iota
	RightUp
	LeftDown
	LeftUp
)

type MapProperties struct {
	Orientation Orientation

	Width  uint
	Height uint

	TileWidth  uint
	TileHeight uint

	Infinite       bool
	TileSideLength int

	StaggerAxis  StaggerAxis
	StaggerIndex StaggerIndex

	TileLayerFormat LayerFormat

	OutputChunkWidth  uint
	OutputChunkHeight uint

	TileRenderOrder  RenderOrder
	CompressionLevel uint

	BackgroundColor color.RGBA
}

type PropertyType string

const (
	TypeBool   PropertyType = "bool"
	TypeColor  PropertyType = "color"
	TypeFloat  PropertyType = "float"
	TypeFile   PropertyType = "file"
	TypeInt    PropertyType = "int"
	TypeString PropertyType = "string"
)

type CustomProperty struct {
	Name string
	Type PropertyType

	value interface{}
}

func (cp CustomProperty) ValueBool() (bool, error) {
	if cp.Type != TypeBool {
		return false, fmt.Errorf("Property %q is type %s, not bool", cp.Name, cp.Type)
	}

	return cp.value.(bool), nil
}

type CustomProperties map[string]CustomProperty

func (cp CustomProperties) GetBoolProperty(key string, fallback bool) bool {
	m := map[string]CustomProperty(cp)
	if prop, ok := m[key]; ok {
		if prop.Type == TypeBool {
			return prop.value.(bool)
		}
	}
	return fallback
}

func (cp CustomProperties) GetIntProperty(key string, fallback int) int {
	m := map[string]CustomProperty(cp)
	if prop, ok := m[key]; ok {
		if prop.Type == TypeInt {
			return prop.value.(int)
		}
	}
	return fallback
}

func (cp CustomProperties) GetStringProperty(key string, fallback string) string {
	m := map[string]CustomProperty(cp)
	if prop, ok := m[key]; ok {
		if prop.Type == TypeString {
			return prop.value.(string)
		}
	}
	return fallback
}

type xmlPropertyList []xmlProperty

type xmlProperty struct {
	XMLName string `xml:"property"`
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Value   string `xml:"value,attr"`
}

func (xp xmlProperty) String() string {
	return fmt.Sprintf("%s:%q", xp.Name, xp.Value)
}

func (pl xmlPropertyList) String() string {
	p := []string{}
	for _, prop := range pl {
		p = append(p, prop.String())
	}
	return strings.Join(p, " ")
}

func (pl xmlPropertyList) GetProperty(name string) string {
	for _, p := range pl {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

func (pl xmlPropertyList) CustomProps() (CustomProperties, error) {
	cp := map[string]CustomProperty{}
	for _, p := range pl {
		prop := CustomProperty{
			Name: p.Name,
			Type: PropertyType(p.Type),
		}

		//TypeBool PropertyType = "bool"
		//TypeColor PropertyType = "color"
		//TypeFloat PropertyType = "float"
		//TypeFile PropertyType = "file"
		//TypeInt PropertyType = "int"
		//TypeString PropertyType = "string"
		switch prop.Type {
		case TypeString:
			prop.value = p.Value
			break
		case TypeBool:
			if p.Value == "true" {
				prop.value = true
			} else {
				prop.value = false
			}
			break
		case TypeInt:
			i64, err := strconv.ParseInt(p.Value, 10, 32)
			if err != nil {
				return nil, err
			}
			prop.value = int(i64)
			break
		default:
			// Default to string instead of breaking.  Doesn't look like
			// the type attribute is written from Tiled if the type is
			// string.
			//return nil, fmt.Errorf("Property type %s not implemented yet", prop.Type)
			prop.value = p.Value
			prop.Type = TypeString
		}

		//cp = append(cp, prop)
		cp[p.Name] = prop
	}

	return cp, nil
}
