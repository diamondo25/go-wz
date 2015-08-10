package wz

type WZVector struct {
	*WZImageObject

	X int32
	Y int32
}

func NewWZVector(name string, parent *WZSimpleNode) *WZVector {
	node := new(WZVector)
	node.WZImageObject = NewWZImageObject(name, parent)
	return node
}

func (m *WZVector) Parse(file *WZFileBlob, offset int64) {
	if file.Debug {
		m.debug(file, "> WZVector2D::Parse")
		defer func() { m.debug(file, "< WZVector2D::Parse") }()
	}

	m.X = file.readWZInt()
	m.Y = file.readWZInt()
}
