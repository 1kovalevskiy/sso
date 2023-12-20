package main

import (
	"flag"
	"log"

	"github.com/1kovalevskiy/sso/config"
	"github.com/1kovalevskiy/sso/internal/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "", "path to config")
	flag.Parse()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)

}
