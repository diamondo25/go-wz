package wz

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/goinggo/workpool"
	"strconv"
)

type WZFileBlob struct {
	reader        *bytes.Reader
	encryption    *Encryption
	file          *WZFile
	contentsStart int32
	Debug         bool
	Name          string
	debug         func(...interface{})
	data          []byte

	workPool *workpool.WorkPool
}

func NewWZFileBlob(data []byte, encryption *Encryption, file *WZFile) *WZFileBlob {
	m := new(WZFileBlob)
	m.contentsStart = 0
	m.encryption = encryption
	m.Debug = file.Debug
	m.Name = file.Filename
	m.debug = file.debug
	m.file = file
	m.data = data
	m.reader = bytes.NewReader(m.data)
	m.workPool = file.workPool

	return m
}

func (m *WZFileBlob) Copy() *WZFileBlob {
	obj := NewWZFileBlob(m.data, m.encryption, m.file)
	obj.contentsStart = m.contentsStart
	return obj
}

func (m *WZFileBlob) CopySliced(start int) *WZFileBlob {
	obj := NewWZFileBlob(m.data[start:], m.encryption, m.file)
	obj.contentsStart = m.contentsStart
	return obj
}

func (m *WZFileBlob) readByte() (out uint8) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readSByte() (out int8) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readInt16() (out int16) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readInt32() (out int32) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readInt64() (out int64) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readUInt16() (out uint16) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readUInt32() (out uint32) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readUInt64() (out uint64) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readFloat32() (out float32) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readFloat64() (out float64) {
	if err := binary.Read(m.reader, binary.LittleEndian, &out); err != nil {
		panic(err)
	}
	return
}

func (m *WZFileBlob) readBytes(size int32) []byte {
	var out []byte = make([]byte, size)

	amount, err := m.reader.Read(out)
	if err != nil {
		panic(err)
	}

	if int32(amount) != size {
		panic(errors.New(fmt.Sprintln("Expected ", size, " bytes, got ", amount)))
	}

	return out
}

// readASCIIZString reads strings until the null terminator
func (m *WZFileBlob) readASCIIZString() string {
	ret := make([]byte, 0)
	for {
		b, err := m.reader.ReadByte()
		if err != nil {
			panic(err)
		}

		if b == 0 {
			break
		}

		ret = append(ret, b)
	}

	return string(ret)
}

func (m *WZFileBlob) readASCIIString(length int32) string {
	return string(m.readBytes(length))
}

// Helper functions
func (m *WZFileBlob) seek(offset int64) {
	if m.Debug {
		m.debug("Seeking ", offset, " bytes (@ ", m.pos(), ")")
	}
	newOffset, err := m.reader.Seek(offset, 0)
	if err != nil {
		panic(err)
	}
	if m.Debug {
		m.debug("New offset ", newOffset, " (@ ", m.pos(), ")")
	}
}

func (m *WZFileBlob) skip(offset int64) {
	_, err := m.reader.Seek(offset, 1)
	if err != nil {
		panic(err)
	}
}

func (m *WZFileBlob) pos() int64 {
	offset, err := m.reader.Seek(0, 1)
	if err != nil {
		panic(err)
	}
	return offset
}

func (m *WZFileBlob) len() int64 {
	curpos := m.pos()
	len, err := m.reader.Seek(0, 2)
	if err != nil {
		panic(err)
	}
	m.seek(curpos)
	return len
}

type AnonFunc func()

func (m *WZFileBlob) peekFor(f AnonFunc) {
	offset, err := m.reader.Seek(0, 1)
	defer func() {
		// Seek back to where we were
		m.seek(offset)
	}()
	if err != nil {
		panic(err)
	}
	f()
}

// WZ data related functions

func (m *WZFileBlob) readDeDuplicatedWZString(uol string, possibleNeededOffset int64, addOffset bool) (result string) {
	key := m.readByte()
	str := ""
	switch key {
	case 0, 0x73:
		str = m.readWZString(uol)
	case 1, 0x1B:
		m.peekFor(func() {
			tmp := possibleNeededOffset
			if addOffset {
				tmp += int64(m.readUInt32())
			} else {
				tmp -= int64(m.readUInt32())
			}
			if m.Debug {
				m.debug("Reading dedup string at ", tmp)
			}
			m.seek(tmp)
			str = m.readWZString(uol)
		})
	default:
		panic("Unknown deduplicated wz string type: " + strconv.Itoa(int(key)) + " at " + uol + " AT " + strconv.Itoa(int(m.pos())))
	}

	if m.Debug {
		m.debug("Dedupped string for ", uol, ": ", str)
	}
	return str
}

func (m *WZFileBlob) readWZObjectUOL(uol string, possibleNeededOffset int64) (result string) {
	key := m.readByte()
	str := ""
	switch key {
	case 0, 0x73:
		str = m.readWZString(uol)
		break
	case 1, 0x1B:
		possibleNeededOffset += int64(m.readUInt32())
		m.peekFor(func() {
			tmp := possibleNeededOffset
			if m.Debug {
				m.debug("Reading wz object uol at ", tmp)
			}
			m.seek(tmp)
			str = m.readWZString(uol)
		})
		break
	default:
		panic("Unknown deduplicated wz string type: " + strconv.Itoa(int(key)) + " at " + uol + " AT " + strconv.Itoa(int(m.pos())))
	}

	if m.Debug {
		m.debug("Dedupped string for ", uol, ": ", str)
	}

	return str
}

func (m *WZFileBlob) readWZString(uol string) (result string) {
	var size int32 = int32(m.readSByte())

	ascii := false

	// Unicode strings are positive, ASCII negative
	if size > 0 {
		if size == 127 {
			size = m.readInt32()
		}
		size *= 2
	} else {
		if size == -128 {
			size = m.readInt32()
		} else {
			size *= -1
		}

		ascii = true
	}

	if size == 0 {
		return ""
	}

	characters := m.readBytes(size)

	var i int32 = 0
	if !ascii {
		var mask uint16 = 0xAAAA

		for ; i < size; i += 2 {
			var char uint16 = uint16(characters[i] | characters[i+1]<<8)
			char ^= mask
			characters[i] = byte(char)
			characters[i+1] = byte(char >> 8)

			mask++
		}

	} else {
		var mask uint8 = 0xAA

		for ; i < size; i++ {
			x := characters[i]
			x ^= mask

			characters[i] = x
			mask++
		}
	}

	if m.encryption != nil && m.encryption.IsEncrypted(uol) {
		m.encryption.TransformBuffer(characters)
	}

	return string(characters)
}

func (m *WZFileBlob) readWZInt() int32 {
	if possibleSize := m.readSByte(); possibleSize == -128 {
		return m.readInt32()
	} else {
		return int32(possibleSize)
	}
}

func (m *WZFileBlob) readWZLong() int64 {
	if possibleSize := m.readSByte(); possibleSize == -128 {
		return m.readInt64()
	} else {
		return int64(possibleSize)
	}
}

func (m *WZFileBlob) readWZOffset() uint32 {
	offset := uint32(m.pos())
	offset = (offset - uint32(m.contentsStart)) ^ 0xFFFFFFFF
	offset *= m.file.versionHash
	offset -= uint32(0x581C3F6D) // Who doesn't like magic values?
	offset = rotl(offset, byte(offset&0x1F))

	encryptedOffset := m.readUInt32()
	offset ^= encryptedOffset
	offset += (uint32(m.contentsStart) * 2)

	return offset
}
