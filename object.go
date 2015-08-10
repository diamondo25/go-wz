package wz

func ParseObject(name string, typename string, parent *WZSimpleNode, file *WZFileBlob, offset int64) interface{} {
	if file.Debug {
		parent.debug(file, "> WZObject::Parse")
		parent.debug(file, typename)
		defer func() { parent.debug(file, "< WZObject::Parse") }()
	}

	switch typename {
	case "Property":
		return ParseProperty(parent, file, offset)

	case "Canvas":
		canvas := NewWZCanvas(name, parent)
    canvas.Parse(file, offset)
		return canvas

	case "Shape2D#Convex2D":
		return ParseConvex(parent, file, offset)

	case "Shape2D#Vector2D":
		vector2d := NewWZVector(name, parent)
		vector2d.Parse(file, offset)
		return vector2d

	case "UOL":
		uol := NewWZUOL(name, parent)
		uol.Parse(file, offset)
		return uol

	case "Sound_DX8":
		sound := NewWZSoundDX8(name, parent)
		sound.Parse(file, offset)
		return sound

	default:
		panic("Unknown typename: " + typename)
	}

}
