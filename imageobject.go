package wz

type WZImageObject struct {
	*WZSimpleNode

	Parse func(file *WZFileBlob, offset int64)
}

func NewWZImageObject(name string, parent *WZSimpleNode) *WZImageObject {
	node := new(WZImageObject)
	node.WZSimpleNode = NewWZSimpleNode(name, parent)
	return node
}
