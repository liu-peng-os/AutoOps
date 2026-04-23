package main

import (
	"dodevops-api/common/config"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/log"
	"fmt"
)

func main() {
	if err := config.LoadConfig(""); err != nil {
		panic("failed to load config: " + err.Error())
	}

	log.Setup()

	if err := db.SetupDBLink(); err != nil {
		fmt.Printf("database connection failed: %v\n", err)
		panic(err)
	}

	fmt.Printf("Database initialized and migrations applied from %s.\n", config.Config.Db.MigrationPath)
}
