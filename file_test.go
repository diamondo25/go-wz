package wz

import (
	"fmt"
	"testing"
)

/*
func TestWZEtcWZ(t *testing.T) {
	return
	fmt.Println("test")
	base, err := NewFile("Etc.wz")
	if err != nil {
		panic(err)
	}
	base.Debug = true
	base.SetWZVariant(VariantGMS)

	base.Parse()

}
*/

func parseDir(dir *WZDirectory) {
	fmt.Println("Found dir: ", dir.GetPath())

	for _, subnode := range dir.Directories {
		parseDir(subnode)
	}

	for name, node := range dir.Images {
		fmt.Println("Found image: ", name, ":", node.GetPath())
		//subnode.Parse()
	}
}

func TestWZDataWZ(t *testing.T) {
	fmt.Println("test")
	base, err := NewFile("Data.wz")
	if err != nil {
		panic(err)
	}
	//base.Debug = true

	base.Parse()
	base.WaitUntilLoaded()
	/*
		for _, node := range base.Root.Directories {
			parseDir(node)
		}

		for name, node := range base.Root.Images {
			fmt.Println("Found image: ", name, ":", node.GetPath())
			//node.Parse()
		}
	*/
}
