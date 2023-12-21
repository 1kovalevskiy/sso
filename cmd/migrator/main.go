package main

import (
	"flag"
	"log"

	"github.com/1kovalevskiy/sso/config"
	"github.com/1kovalevskiy/sso/pkg/migrator"
)

func main() {
	var configPath, migrationsPath string

	flag.StringVar(&configPath, "config-path", "", "path to config")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if configPath == "" {
		log.Fatalf("Set config-path")
	}
	if migrationsPath == "" {
		log.Fatalf("Set migrations-path")

	}

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	storagePath := cfg.SQL.URL
	migrationsTable := "migrations"


	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	migrator.MigrateUp(migrationsPath, storagePath, migrationsTable)
}
