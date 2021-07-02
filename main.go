package main

import (
	"os"
)

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	err = cfg.GenerateGormCode()
	if err != nil {
		os.Exit(1)
	}

}
