package tiled

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Layer struct {
	Id         int
	Name       string
	Width      uint
	Height     uint
	Data       []uint32 // use a tile struct?
	Properties CustomProperties
}

type xmlLayer struct {
	Id         int             `xml:"id,attr"`
	Name       string          `xml:"name,attr"`
	Width      int             `xml:"width,attr"`
	Height     int             `xml:"height,attr"`
	Data       xmlLayerData    `xml:"data"`
	Properties xmlPropertyList `xml:"properties>property"`
}

type xmlLayerData struct {
	//XMLName    string       `xml:"data"`
	Data        []byte `xml:",innerxml"`
	Encoding    string `xml:"encoding,attr"`
	Compression string `xml:"compression,attr"`
}

func (a Layer) Merge(b Layer) (Layer, error) {
	if a.Width != b.Width || a.Height != b.Height {
		return Layer{}, fmt.Errorf("Layer dimension mismatch: %dx%d vs %dx%d", a.Width, a.Height, b.Width, b.Height)
	}

	n := Layer{
		Name:       fmt.Sprintf("%s + %s", a, b),
		Width:      a.Width,
		Height:     a.Height,
		Data:       make([]uint32, len(a.Data)),
		Properties: CustomProperties{}, // merge this too?
	}

	copy(n.Data, a.Data)
	for i, d := range b.Data {
		if d != 0 {
			n.Data[i] = d
		}
	}

	return n, nil
}

func decodeLayerData(encoding, compression string, data []byte) ([]uint32, error) {
	var uncompressed []byte
	switch compression {
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, reader); err != nil {
			return nil, err
		}
		uncompressed = buf.Bytes()
		break

	//case "zlib":
	//	break
	//case "zstd":
	//	break

	case "":
		uncompressed = data
		break

	default:
		return nil, fmt.Errorf("Unsupported compression format: %q", compression)
	}

	d := []uint32{}
	switch encoding {
	case "csv":
		split := strings.Split(strings.ReplaceAll(strings.ReplaceAll(string(uncompressed), "\n", ""), "\r", ""), ",")
		for _, s := range split {
			u64, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("Error parsing data: %v", err)
			}

			d = append(d, uint32(u64))
		}
		break
	case "base64":
		buf := bytes.NewBuffer(uncompressed)
		dec := base64.NewDecoder(base64.StdEncoding, buf)
		byteData, err := io.ReadAll(dec)
		if err != nil {
			return nil, fmt.Errorf("Error decoding base64 data: %v", err)
		}

		if len(byteData)%4 != 0 {
			return nil, fmt.Errorf("Invalid base64 data length: %d", len(byteData))
		}

		for i := 0; i < len(byteData); i += 4 {
			duint := uint32(byteData[i+0]) | uint32(byteData[i+1])<<8 | uint32(byteData[i+2])<<16 | uint32(byteData[i+3])<<24
			d = append(d, duint)
		}
	default:
		return nil, fmt.Errorf("Unsupported encoding: %q", encoding)
	}

	return d, nil
}

func decodeLayers(input []xmlLayer) ([]Layer, error) {
	lst := []Layer{}
	for _, l := range input {
		data, err := decodeLayerData(l.Data.Encoding, l.Data.Compression, l.Data.Data)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode layer %q: %v", l.Name, err)
		}

		cp, err := l.Properties.CustomProps()
		if err != nil {
			return nil, fmt.Errorf("Error parsing layer properties: %v", err)
		}

		layer := Layer{
			Id:         l.Id,
			Name:       l.Name,
			Data:       data,
			Width:      uint(l.Width),
			Height:     uint(l.Height),
			Properties: cp,
		}

		lst = append(lst, layer)
	}

	return lst, nil
}
