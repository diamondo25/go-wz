package main

import (
"bufio"
"fmt"
  "github.com/diamondo25/go-wz"
  "os"
)

func main() {
	base, err := wz.NewFile("Data.wz")
	if err != nil {
		panic(err)
	}
	//base.Debug = true

  fmt.Println("Loading ", base.Filename)
	base.Parse()

  base.WaitUntilLoaded()

  fmt.Println("Loaded!")

  base.GetFromPath("smap.img")
  base.GetFromPath("Effect/Tomb.img/fall/0/z")

  reader := bufio.NewReader(os.Stdin)
  reader.ReadString('\n')

  fmt.Println("Shutting Down")
}
