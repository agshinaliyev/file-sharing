package main

import (
	"file-sharing/config"
	"file-sharing/db"
	"file-sharing/logger"
	"file-sharing/server"
	log "github.com/sirupsen/logrus"
)

func main() {

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}
	file, err := logger.Init()

	if err != nil {
		log.Fatal(err)
	}
	if err := db.Init(); err != nil {
		log.Fatal("Error initializing database", err)
	}

	if err := server.Init(); err != nil {
		log.Fatal("Error initializing server", err)

	}
	defer func() {
		log.Info("Server stopped...")
		file.Close()
	}()

}
