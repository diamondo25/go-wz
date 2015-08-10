package wz

import (
	"errors"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"github.com/goinggo/workpool"
	"os"
	"path/filepath"
	"strconv"
	"runtime"
	"time"
)

type WZFile struct {
	filemap     mmap.MMap
	versionHash uint32
	mainBlob    *WZFileBlob

	workPool *workpool.WorkPool

	FileDescription string
	Debug           bool
	Filename        string
	Root            *WZDirectory
}

func NewFile(filename string) (*WZFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	filemap, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	wz := new(WZFile)
	wz.filemap = filemap
	wz.Debug = false
	wz.Filename = filename
	wz.workPool = workpool.New(runtime.NumCPU() * 2, 7000)
	wz.mainBlob = NewWZFileBlob(wz.filemap, nil, wz)

	return wz, nil
}

func (m *WZFile) debug(args ...interface{}) {
	if m.Debug {
		fmt.Println(fmt.Sprint("[WZFile: ", m.Filename, "] ", fmt.Sprint(args...)))
	}
}

func (m *WZFile) Close() {
	m.filemap.Unmap()
}

func (m *WZFile) Parse() {
  runtime.GOMAXPROCS(runtime.NumCPU())

	m.debug("Starting parsing...")
	m.mainBlob.seek(0)

	header := m.mainBlob.readASCIIString(4)

	m.debug("Header: ", header)
	if header != "PKG1" {
		panic(errors.New("Not a PKG1/WZ file"))
	}

	m.mainBlob.skip(8) // Filesize

	m.mainBlob.contentsStart = m.mainBlob.readInt32()
	m.debug("Content starts at ", m.mainBlob.contentsStart)

	m.FileDescription = m.mainBlob.readASCIIZString()
	m.debug("File description: ", m.FileDescription)

	m.determineVersion()
}

// determineVersion is a bruteforcer on the hash stored inside the
// wz file.
func (m *WZFile) determineVersion() {
	m.mainBlob.seek(int64(m.mainBlob.contentsStart))

	encryptedVersion := m.mainBlob.readUInt16()
	var realVersion uint16 = 0

	calculateHash := func(versionNumber uint16) (uint16, uint32) {

		versionAsString := strconv.Itoa(int(versionNumber))
		b := []byte(versionAsString)
		// Should return "31 33" on ver 13

		var x uint16 = 0xFF
		var y uint32 = 0
		for _, val := range b {
			x ^= uint16(val + 1) // Lolwat.
			y = (y << 8) | uint32(val+1)
		}

		return x, y
	}

	for {
		realVersion++
		calcVersion, calcHash := calculateHash(realVersion)
		if calcVersion != encryptedVersion {
			m.debug("It cannot be version ", realVersion)
		} else {
			m.debug("It is probably version ", realVersion, "! (hash ", calcHash, ")")
			m.versionHash = calcHash
			// Now, see if we can actually do something with this version
			if dir := m.isParsableWithVersion(); dir != nil {
				m.debug("Yes, this is usable!")

				m.Root = dir

				return
			} else {
				m.debug("Nope, not the correct version")
				continue
			}

		}
	}
}

func (m *WZFile) isParsableWithVersion() (result *WZDirectory) {
	defer func() {
		if r := recover(); r != nil {
			m.debug("Its not this version, reason: ", r)
			panic(r)
			result = nil
		}
	}()

	dir := NewWZDirectory(filepath.Base(m.Filename), nil)
	dir.Parse(m.mainBlob, m.mainBlob.pos())

	return dir
}

func (m *WZFile) WaitUntilLoaded() {
  for m.workPool.QueuedWork() != 0 {
    time.Sleep(100 * time.Millisecond)
  }
}
