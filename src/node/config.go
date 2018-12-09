package main

import (
    "log"
    "os"
    "encoding/json"
)

type Configuration struct {
    NetworkInterface string
}

// Create a new Configuration object with default values
func initConfig() (*Configuration) {
    config := Configuration{}
    config.NetworkInterface = "eth0"

    return &config
}

func ReadConfigFile(filepath string) (*Configuration) {
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
