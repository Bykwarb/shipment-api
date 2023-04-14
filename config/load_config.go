package config

import (
	"log"
	"os"

	"task/iternal/model"

	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) *model.Config {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Failed to read YAML file: %v", err)
	}
	var config model.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Panicf("Failed to unmarshal YAML: %v", err)
	}
	return &config
}
