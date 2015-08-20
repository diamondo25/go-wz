package wz

import (
	"crypto/aes"
	"strconv"
)

// If a WZ key is configured, it's copied 4 times...
var GMS_WZ_IV = []byte{0x4D, 0x23, 0xC7, 0x2B, 0x4D, 0x23, 0xC7, 0x2B, 0x4D, 0x23, 0xC7, 0x2B, 0x4D, 0x23, 0xC7, 0x2B}
var SEA_WZ_IV = []byte{0xB9, 0x7D, 0x63, 0xE9, 0xB9, 0x7D, 0x63, 0xE9, 0xB9, 0x7D, 0x63, 0xE9, 0xB9, 0x7D, 0x63, 0xE9}

// The default WZ key (if none is configured)
var DEFAULT_WZ_IV = []byte{0x53, 0xF2, 0xA8, 0x42, 0x9D, 0x7F, 0x77, 0x09, 0x1D, 0x26, 0x42, 0x53, 0x88, 0x7C, 0x73, 0x3A}

// WZ_AES_KEY is a cut down version of the defined AES key inside the client!
var WZ_AES_KEY = []byte{
	0x13, 0x00, 0x00, 0x00,
	0x08, 0x00, 0x00, 0x00,
	0x06, 0x00, 0x00, 0x00,
	0xB4, 0x00, 0x00, 0x00,
	0x1B, 0x00, 0x00, 0x00,
	0x0F, 0x00, 0x00, 0x00,
	0x33, 0x00, 0x00, 0x00,
	0x52, 0x00, 0x00, 0x00}

func expandXorKey(currentIV, aesKey, currentXorKey []byte, neededLength int) (newIV []byte, newXorKey []byte) {
	finalLength := int(neededLength / 16)
	if neededLength%16 != 0 {
		finalLength += 1
	}
	finalLength *= 16
	finalLength -= len(currentXorKey)

	nextBlock := make([]byte, finalLength-len(currentXorKey))

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}

	// For each 16 bytes, encrypt the IV (result will be put into nextBlock), then use the result to do the next one...
	for i := 0; i < len(nextBlock); i += 16 {
		curblock := nextBlock[i : i+16]
		block.Encrypt(curblock, currentIV)
		currentIV = curblock
	}

	currentXorKey = append(currentXorKey, nextBlock...)

	return currentIV, currentXorKey
}

// rotl is a Bitshift to the left, that'll put the bits pushed off at the right
func rotl(value uint32, times uint8) uint32 {
	return uint32((value << times) | (value >> (32 - times)))
}

// rotr is a Bitshift to the right, that'll put the bits pushed off at the left
func rotr(value uint32, times uint8) uint32 {
	return uint32((value >> times) | (value << (32 - times)))
}

func calculateHash(versionNumber uint16) (uint16, uint32) {
	versionAsString := strconv.Itoa(int(versionNumber))
	b := []byte(versionAsString)
	// Should return "31 33" on ver 13

	var y uint32 = 0
	for _, val := range b {
		y = (y << 5)
		y += uint32(val + 1)
	}

	var x uint16 = 0xFF
	x ^= uint16((y >> 24) & 0xFF)
	x ^= uint16((y >> 16) & 0xFF)
	x ^= uint16((y >> 8) & 0xFF)
	x ^= uint16((y >> 0) & 0xFF)

	return x, y
}
