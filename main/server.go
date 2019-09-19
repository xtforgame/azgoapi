package main

import (
	"github.com/xtforgame/azgoapi/agapiserver"
	"github.com/xtforgame/azgoapi/utils"
)

func main() {
	defer utils.FinalReport()
	// os.Exit(0)

	hs := agapiserver.NewHttpServer()
	hs.Init()
	hs.Start()
}
