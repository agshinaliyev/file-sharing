package config

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type args struct {
	//application config
	AppName  string `arg:"env:APP_NAME"`
	LogLevel string `arg:"env:LOG_LEVEL"`
	//database config
	DBHost string `arg:"env:DB_HOST"`
	DBPort string `arg:"env:DB_PORT"`
	DBUser string `arg:"env:DB_USER"`
	DBPass string `arg:"env:DB_PASS"`
	DBName string `arg:"env:DB_NAME"`
	UserId string `arg:"env:USER_TRANSPORT_KEY"`
	//server config
	Port string `arg:"env:APP_PORT"`
}

var Opts struct {
	Profile string `short:"p" long:"profile" default:"default" description:"Application run profile"`
}

func IsDefaultProfile() bool {
	return Opts.Profile == "default"
}

// Conf - struct of application configuration data
var Conf args

func Init() error {
	_, err := flags.Parse(&Opts)
	if err != nil {
		return err
	}

	fileName := "default.env"
	if err := godotenv.Load(fileName); err != nil {
		return errors.New(fmt.Sprintf("Error in loading environment variables from %s: %s", fileName, err))
	} else {
		log.Info("Environment variables loaded from: default.env")
	}

	if !IsDefaultProfile() {
		profileFileName := Opts.Profile + ".env"
		if err := godotenv.Overload(profileFileName); err != nil {
			return errors.New(fmt.Sprintf("Error in loading environment variables from %s: %s", profileFileName, err))
		} else {
			log.Info("Environment variables overloaded from: ", profileFileName)
		}
	}

	_ = arg.Parse(&Conf)

	log.Info("Application is configured from profile: ", Opts.Profile)
	return nil
}
