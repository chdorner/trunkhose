package main

import (
	"fmt"
	"os"
	"path"

	flag "github.com/spf13/pflag"
)

func main() {
	var configPath string
	var historyPath string
	flag.StringVarP(&configPath, "config", "c", "", "path to the config file")
	flag.StringVar(&historyPath, "history", "", "path to the history file")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintln(os.Stderr, "missing path to config file")
		os.Exit(1)
	}

	if historyPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		historyPath = path.Join(wd, "trunkhose-history.json")
	}

	config, err := ParseConfig(configPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	err = process(config, historyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
