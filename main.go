package main

import (
	"fmt"
	"log"
	"net/http"

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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serviceConfig.Port), serviceMux); err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}
