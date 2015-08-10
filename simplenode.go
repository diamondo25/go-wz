package wz

import "fmt"

type WZSimpleNode struct {
	*WithParent
	Name string

	cachedPath    string
	hasCachedPath bool
}

func NewWZSimpleNode(name string, parent *WZSimpleNode) *WZSimpleNode {
	node := new(WZSimpleNode)
	node.WithParent = &WithParent{parent}
	node.Name = name
	return node
}

func (m *WZSimpleNode) GetPath() string {
	if !m.hasCachedPath {
		x := m
		buffer := ""
		for x != nil {
			buffer = "/" + x.Name + buffer
			x = x.Parent
		}

		m.cachedPath = buffer[1:]
		m.hasCachedPath = true
	}

	return m.cachedPath
}

func (m *WZSimpleNode) debug(file *WZFileBlob, args ...interface{}) {
	if !file.Debug {
		return
	}
	tmp := make([]interface{}, 1)
	tmp[0] = fmt.Sprint("[ ", file.pos(), " : ", m.GetPath(), "]")
	tmp = append(tmp, args...)
	file.debug(tmp...)
}
