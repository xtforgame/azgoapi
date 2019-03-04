package main

import (
	"github.com/xtforgame/gbol/tests/basicws"
	"github.com/xtforgame/gbol/utils"
)

func main() {
	defer utils.FinalReport()
	// os.Exit(0)

	hs := basicws.NewHttpServer()
	hs.Init()
	hs.Start()
}
