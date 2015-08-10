package wz

type WZProperty map[string]*WZVariant

func ParseProperty(parent *WZSimpleNode, file *WZFileBlob, offset int64) WZProperty {
	if file.Debug {
		parent.debug(file, "> WZProperty::Parse")
		defer func() { parent.debug(file, "< WZProperty::Parse") }()
	}

	file.skip(2) // Unk
	propcount := int(file.readWZInt())

	if file.Debug {
  	parent.debug(file, "Properties of ", parent.GetPath(), ": ", propcount)
  }

	variants := make(WZProperty)

	for i := 0; i < propcount; i++ {
		name := file.readWZObjectUOL(parent.GetPath(), offset)
		if file.Debug {
		  parent.debug(file, "Prop ", i, " has name ", name)
		}
		variant := NewWZVariant(name, parent)
		variant.Parse(file, offset)
		variants[name] = variant
	}

	return variants
}
