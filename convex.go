package wz

import "strconv"

func ParseConvex(parent *WZSimpleNode, file *WZFileBlob, offset int64) []interface{} {
	if file.Debug {
		parent.debug(file, "> WZConvex::Parse")
		defer func() { parent.debug(file, "< WZConvex::Parse") }()
	}

	propcount := int(file.readWZInt())
	if file.Debug {
		parent.debug(file, "Object count: ", propcount)
	}
	objects := make([]interface{}, propcount)

	for i := 0; i < propcount; i++ {
		typename := file.readWZObjectUOL(parent.GetPath(), offset)

		if file.Debug {
			parent.debug(file, "Prop ", i, " has typename ", typename)
		}

		objects[i] = ParseObject(strconv.Itoa(i), typename, parent, file, offset)
	}

	return objects
}
