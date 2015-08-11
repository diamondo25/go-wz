package wz

import (
	"encoding/hex"
	"fmt"
)

type WZVariant struct {
	*WZImageObject

	Type  uint8
	Value interface{}
}

func NewWZVariant(name string, parent *WZSimpleNode) *WZVariant {
	node := new(WZVariant)
	node.WZImageObject = NewWZImageObject(name, parent)
	return node
}

func (m *WZVariant) Parse(file *WZFileBlob, offset int64) {
	if file.Debug {
		m.debug(file, "> WZVariant::Parse")
		defer func() { m.debug(file, "< WZVariant::Parse") }()
	}

	m.Type = file.readByte()

	if file.Debug {
		m.debug(file, "Type: ", m.Type)
	}

	switch m.Type {
	// no data
	case 0:
		m.Value = nil
		break // Nothing

	// int16
	case 2, 11:
		m.Value = file.readInt16()
		break

	// int32
	case 3, 19:
		m.Value = int32(file.readWZInt())
		break

	// int64
	case 20:
		m.Value = file.readWZLong()
		break

	// float32
	case 4:
		if file.readByte() == 0x80 {
			m.Value = file.readFloat32()
		} else {
			m.Value = float32(0.0)
		}
		break

	// float64
	case 5:
		m.Value = file.readFloat64()
		break

	// String
	case 8:
		m.Value = file.readWZObjectUOL(m.GetPath(), offset)
		break

		// Sub object
	case 9:
		size := int64(file.readInt32())
		x := file.pos()

		if file.Debug {
			m.debug(file, "Size: ", size, " - x: ", x)
			m.debug(file, "Offset: ", offset)
		}
		typename := file.readWZObjectUOL(m.GetPath(), offset)

		if file.Debug {
			m.debug(file, "typename: ", typename)
		}

		m.Value = ParseObject(m.Name, typename, m.WZSimpleNode, file, offset)

		if x+size != file.pos() {
			x += size
			x -= file.pos()
			m.debug(file, "NOT ENOUGH PARSED: ", x, " bytes left???")
			if x > 0 {

				m.debug(file, "Bytes: \n", hex.Dump(file.readBytes(int32(x))))
			} else {
				file.skip(x)
			}
		}
		break

	default:
		panic(fmt.Sprint("Unknown wz prop type ", m.Type, " at ", m.GetPath(), " AT ", file.pos()))
	}
}
