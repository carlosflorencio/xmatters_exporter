package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Config Configuration
)

type Configuration struct {
	Debug bool   `arg:"env:DEBUG"`
	Url   string `arg:"env:XMATTERS_URL"`
	Token string `arg:"env:XMATTERS_TOKEN"`
}

func init() {
	Config = Configuration{
		Debug: false,
		Url:   "https://company.xmatters.com/api/xm/1",
	}
	arg.MustParse(&Config)

	if Config.Token == "" {
		log.Fatal().Msg("Environment variable XMATTERS_TOKEN is required.")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func main() {
	fmt.Println(Config)
}


