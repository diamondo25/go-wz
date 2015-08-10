package wz

type WZUOL struct {
	*WZImageObject

	Reference string
}

func NewWZUOL(name string, parent *WZSimpleNode) *WZUOL {
	node := new(WZUOL)
	node.WZImageObject = NewWZImageObject(name, parent)
	return node
}

func (m *WZUOL) Parse(file *WZFileBlob, offset int64) {
	if file.Debug {
		m.debug(file, "> WZUOL::Parse")
		defer func() { m.debug(file, "< WZUOL::Parse") }()
	}

	file.skip(1) // Version number?
	m.Reference = file.readWZObjectUOL(m.GetPath(), offset)
}
