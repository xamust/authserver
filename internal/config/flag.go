package config

import (
	"flag"
	"os"
)

func Flag() string {
	var confFile string
	flag.StringVar(&confFile, "config", "", "config file path")
	flag.Parse()
	if confFile == "" {
		if _, err := os.Stat(confFile); os.IsNotExist(err) {
			panic("config file does not exist")
		}
		panic("config file is not set")
	}
	return confFile
}
