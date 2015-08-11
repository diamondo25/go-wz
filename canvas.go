package wz

type WZCanvas struct {
	*WZImageObject

	Width    int32
	Height   int32
	Format   int32
	MagLevel uint8

	data []byte

	Properties WZProperty
}

func NewWZCanvas(name string, parent *WZSimpleNode) *WZCanvas {
	node := new(WZCanvas)
	node.WZImageObject = NewWZImageObject(name, parent)
	return node
}

func (m *WZCanvas) Parse(file *WZFileBlob, offset int64) {
	if file.Debug {
		m.debug(file, "> WZCanvas::Parse")
		defer func() { m.debug(file, "< WZCanvas::Parse") }()
	}
	file.skip(1)

	if file.readByte() == 1 {
		m.Properties = ParseProperty(m.WZSimpleNode, file, offset)
	}

	m.Width = file.readWZInt()
	m.Height = file.readWZInt()

	if m.Width >= 0x10000 || m.Height >= 0x10000 {
		panic("File corrupt? Width and/or Height is too big.")
	}

	m.Format = file.readWZInt()
	m.MagLevel = file.readByte()

	if file.readInt32() != 0 {
		panic("4 bytes must equal zero.")
	}

	len := file.readInt32()

	m.debug(file, "Canvas len: ", len)
	len -= 1
	// skip first byte
	file.skip(1)

	m.data = file.readBytes(len)
}
