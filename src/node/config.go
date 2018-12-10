package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	NetworkInterface string
	SharedDirectory  string
	ChunkSize        string
	ChunkSizeInt     int
}

// Create a new Configuration object with default values
func initConfig() *Configuration {
	config := Configuration{}
	config.NetworkInterface = "eth0"
	config.SharedDirectory = "/tmp"
	config.ChunkSize = "2048" // 2KB
	config.ChunkSizeInt = 2048
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

	// Update ChunkSizeInt
	tempInt64, _ := strconv.ParseInt(config.ChunkSize, 10, 64)
	config.ChunkSizeInt = int(tempInt64)
	return config
}
