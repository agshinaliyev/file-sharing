package logger

import (
	"file-sharing/config"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// TODO Configure these values via external configuration
const (
	LOG_FILE = "tracker.log"
)

func Init() (*os.File, error) {
	logLevel, err := log.ParseLevel(config.Conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(logLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := OpenFile(LOG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))

	log.Info("Logging started...")
	log.Info("Log level is set to: ", logLevel)
	return file, nil
}

func OpenFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
