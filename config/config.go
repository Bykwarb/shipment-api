package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Filter struct {
		Enabled                  bool
		ExpectedNumElements      int     `yaml:"expected_num_elements"`
		FalsePositiveProbability float64 `yaml:"false_positive_probability"`
	}
}

func LoadConfig(path string) *Config {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Failed to read YAML file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Panicf("Failed to unmarshal YAML: %v", err)
	}
	return &config
}
