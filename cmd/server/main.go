package main

import (
	"log"

	"github.com/AntonPashechko/gophkeeper/internal/logger"
	"github.com/AntonPashechko/gophkeeper/internal/server/app"
	"github.com/AntonPashechko/gophkeeper/internal/server/config"
	"github.com/AntonPashechko/gophkeeper/internal/storage"
)

func main() {
	//Инициализируем синглтон логера
	if err := logger.Initialize("info"); err != nil {
		log.Fatalf("cannot initialize logger: %s\n", err)
	}

	cfg, err := config.Create()
	if err != nil {
		log.Fatalf("cannot load config: %s\n", err)
	}

	storage, err := storage.NewKeeperStorage(cfg.DataBaseDNS)
	if err != nil {
		log.Fatalf("cannot create db store: %s\n", err)
	}
	defer storage.Close()

	app, err := app.Create(cfg, storage)
	if err != nil {
		logger.Error("cannot create app: %s", err)
		return
	}

	go app.Run()

	logger.Info("Running server: address %s", cfg.Endpoint)

	<-app.ServerDone()

	if err := app.Shutdown(); err != nil {
		logger.Error("Server shutdown failed: %s", err)
	}

	logger.Info("Server has been shutdown")
}
