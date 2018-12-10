package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	NetworkInterface string
    SharedDirectory  string
}

// Create a new Configuration object with default values
func initConfig() *Configuration {
	config := Configuration{}
	config.NetworkInterface = "eth0"
    config.SharedDirectory = "/tmp"
	return &config
}

func ReadConfigFile(filepath string) *Configuration {
	config := initConfig()
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("[Configuration] Unable to open the config file. Setting default values")
		return config
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		log.Printf("[Configuration] Unable to open the config file. Setting default values")
		return config
	}

	return config
}
