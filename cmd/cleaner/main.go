package main

import (
	"flag"
	"whale-cleaner/cleaner"

	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "config path")
	dry := flag.Bool("dry", false, "dry run")
	flag.Parse()

	config, err := cleaner.LoadConfig(*configPath)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load config.")
	}

	if err := config.Validate(); err != nil {
		logrus.WithError(err).Fatal("Configuration is invalid.")
	}

	if *dry {
		logrus.Info("Dry run. Images will not be deleted.")
	}

	for _, repository := range config.Repositories {
		logrus.Infof("Cleaning repository: %s", repository.Name)
		err := cleaner.Clean(repository, *dry)
		if err != nil {
			logrus.WithError(err).Error("Error during repository clean")
		}
	}
}
