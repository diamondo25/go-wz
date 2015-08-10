package wz

import (
	//"bytes"
	"strings"
)

type Encryption struct {
	encryptedStrings []string
	aesIV            []byte
	aesKey           []byte
	xorKey           []byte
}

const VariantGMS = byte(1)
const VariantSEA = byte(2)

func NewEncryption(variant byte) *Encryption {
	m := new(Encryption)
	m.setWZVariant(variant)
	return m
}

func (m *Encryption) IsEncrypted(uol string) bool {
	for _, str := range m.encryptedStrings {
		if strings.Index(uol, str) == 0 {
			return true
		}
	}

	return false
}

func (m *Encryption) TransformBuffer(buffer []byte) {
	m.tryExpandXorKey(len(buffer))
	for i := 0; i < len(buffer); i++ {
		buffer[i] ^= m.xorKey[i]
	}
}

func (m *Encryption) setWZVariant(variant byte) {
	switch variant {
	case VariantGMS:
		m.aesIV = GMS_WZ_IV
		m.aesKey = WZ_AES_KEY
		break
	case VariantSEA:
		m.aesIV = SEA_WZ_IV
		m.aesKey = WZ_AES_KEY
		break
	default:
		// When the WZ key is set to this, do not expect good results
		// There has been no version that used this key yet.
		m.aesIV = DEFAULT_WZ_IV
		m.aesKey = WZ_AES_KEY
	}

	m.xorKey = []byte{}
	m.tryExpandXorKey(400) // Pre-built a WZ key for 400 characters. Should be enough for most of the simple data/strings.
}

func (m *Encryption) tryExpandXorKey(length int) {
	// Check if we already have enough data
	if len(m.xorKey) < length {
		return
	}

	m.aesIV, m.xorKey = expandXorKey(m.aesIV, m.aesKey, m.xorKey, length)
}
