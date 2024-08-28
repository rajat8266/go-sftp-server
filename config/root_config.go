package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// Config holds the configuration for the SFTP daemon.
var RootConfig Config

// LoadRootConfig parses the config file and returns the configuration.
func LoadRootConfig() *Config {
	configPath := flag.String("config-path", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Please specify the config-path parameter")
		flag.Usage()
		os.Exit(1)
	}

	parseRootConfig(*configPath)
	return &RootConfig
}

// parseRootConfig reads the JSON config file and populates the RootConfig variable.
func parseRootConfig(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading config file at %s: %v", path, err)
	}

	if err := json.Unmarshal(content, &RootConfig); err != nil {
		log.Fatalf("failed to unmarshal config data: %v", err)
	}
}
