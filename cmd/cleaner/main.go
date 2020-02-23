package main

import (
	"flag"
	"fmt"
	"os"
	"whale-cleaner/cleaner"

	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("config path is required")
		os.Exit(1)
	}

	dryPtr := flag.Bool("dry", false, "dry run")
	flag.Parse()

	configPath := os.Args[1]
	config, err := cleaner.LoadConfig(configPath)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load config.")
	}

	if err := config.Validate(); err != nil {
		logrus.WithError(err).Fatal("Configuration is invalid.")
	}

	for _, repository := range config.Repositories {
		logrus.Infof("Cleaning repository: %s", repository.Name)
		err := cleaner.Clean(repository, *dryPtr)
		logrus.WithError(err).Error("Error during repository clean")
	}
}
