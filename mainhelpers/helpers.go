package mainhelpers

import (
	"flag"
	// "fmt"
	"github.com/xtforgame/azgoapi/config"
	"github.com/xtforgame/azgoapi/mainservice"
)

func ParseConfig() (*config.Config, error) {
	var (
		conf = flag.String("config", "", "config filename")
	)
	flag.Parse()

	cfg, err := config.ParseConfig(*conf)
	if err != nil {
		panic(err)
	}

	return cfg, err
}

func NewSbMainServiceForDev() *mainservice.SbMainService {
	cfg, _ := ParseConfig()
	return mainservice.NewSbMainService(cfg, mainservice.SbMainServiceOptions{})
}

func NewSbMainServiceForProd() *mainservice.SbMainService {
	cfg, _ := ParseConfig()
	return mainservice.NewSbMainService(cfg, mainservice.SbMainServiceOptions{})
}
