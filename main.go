package main

import (
	"log"

	"geekible.todolist/src/config"
	"geekible.todolist/src/routers"
)

func main() {
	cfg := config.InitServiceConfig()
	serviceConfig, err := cfg.BuildConfig()
	if err != nil {
		log.Fatalf("error getting service configuration: %v", err)
	}

	migration := config.InitMigration(serviceConfig)
	migration.DoMigration()

	serviceMux := cfg.BuilderMux()

	routers.InitToDoRoutes(serviceMux, cfg).RegisterRoutes()
}
