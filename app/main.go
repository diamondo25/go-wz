package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/diamondo25/go-wz"
	"os"
)

func main() {
	base, err := wz.NewFile("Data.wz")
	if err != nil {
		panic(err)
	}
	// base.Debug = true
	flag.BoolVar(&base.Debug, "debug", false, "Toggles debugging mode")
	flag.BoolVar(&base.LazyLoading, "lazyloading", true, "If disabled, all data will be loaded in memory")
	flag.Parse()

	fmt.Println("Loading ", base.Filename)
	base.Parse()

	base.WaitUntilLoaded()

	fmt.Println("Loaded!")

	fmt.Println("Smap: ", base.GetFromPath("smap.img"))
	fmt.Println("Falling tomb thing z value: ", base.GetFromPath("Effect/Tomb.img/fall/0/z"))
	fmt.Println("Falling tomb thing origin X: ", base.GetFromPath("Effect/Tomb.img/fall/0/origin/X"))

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	fmt.Println("Shutting Down")
}
