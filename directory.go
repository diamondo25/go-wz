package wz

import (
	"fmt"
)

type WZDirectoryLoader struct {
	Directory *WZDirectory
	FileBlob  *WZFileBlob
	Offset    int64
}

func (m *WZDirectoryLoader) DoWork(workRoutine int) {
	m.Directory.Parse(m.FileBlob, m.Offset)
}

type WZImageLoader struct {
	Image    *WZImage
	FileBlob *WZFileBlob
	Offset   int64
}

func (m *WZImageLoader) DoWork(workRoutine int) {
	m.Image.Parse(m.FileBlob, m.Offset)
}

type WZDirectory struct {
	*WZSimpleNode

	Directories map[string]*WZDirectory
	Images      map[string]*WZImage
}

func NewWZDirectory(name string, parent *WZSimpleNode) *WZDirectory {
	node := new(WZDirectory)
	node.WZSimpleNode = NewWZSimpleNode(name, parent)
	node.Directories = make(map[string]*WZDirectory)
	node.Images = make(map[string]*WZImage)

	return node
}

func (m *WZDirectory) Parse(file *WZFileBlob, offset int64) {
	file.seek(offset)

	entries := file.readWZInt()
	var i int32 = 0

	for ; i < entries; i++ {
		elementType := file.readByte()
		var name string = ""

		switch elementType {
		case 1:
			someData := file.readBytes(10)
			m.debug(file, "WZDirectory::Parse  found type 1: ", someData)

			continue // What does these 10 bytes contain?
		case 2: // UOL basically
			subOffset := int64(file.readInt32() + file.contentsStart)
			file.peekFor(func() {
				file.seek(subOffset)
				elementType = file.readByte()
				name = file.readWZString(m.GetPath())
			})
			break
		case 3, 4:
			name = file.readWZString(m.GetPath())
			break
		default:
			panic(fmt.Sprint("Unknown type in directory? ", elementType))
			return
		}

		/*size := */ file.readWZInt() // Blob size
		file.readWZInt()              // Checksum?
		dataOffset := int64(file.readWZOffset())
		curpos := file.pos()

		if elementType == 3 {

			newDir := NewWZDirectory(name, m.WZSimpleNode)
			m.Directories[name] = newDir
			if true {
				work := new(WZDirectoryLoader)
				work.Directory = newDir
				work.FileBlob = file.Copy()
				work.Offset = dataOffset

				if err := file.workPool.PostWork("directory loader", work); err != nil {
					fmt.Println("ERROR ", err)
				}
			} else {
				newDir.Parse(file, dataOffset)
			}

		} else {
			img := NewWZImage(name, m.WZSimpleNode)
			m.Images[name] = img
			if !file.file.LazyLoading {
				if false {
					// Goroutine spamming
					subfile := file.CopySliced(int(dataOffset))
					img.Parse(subfile, 0)
				} else if true {
					// Workpool
					work := new(WZImageLoader)
					work.Image = img
					work.FileBlob = file.Copy()
					work.Offset = dataOffset

					if err := file.workPool.PostWork("image loader", work); err != nil {
						fmt.Println("ERROR ", err)
					}
				} else {
					// Sync loading
					img.Parse(file, dataOffset)
				}
			} else {
				img.parseFuncInfo = func() {
					img.Parse(file, dataOffset)
				}
			}
		}
		file.seek(curpos)
	}
}
