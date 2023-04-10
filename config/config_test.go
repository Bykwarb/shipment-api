package config

import (
	"testing"
)

func TestLoadConfigFromYML(t *testing.T) {

	config := LoadConfig("../config.yml")
	expectedConfig := Config{}
	expectedConfig.Server.Port = "8080"
	expectedConfig.Server.Host = "localhost"

	expectedConfig.Database.DBName = "shipment_api"
	expectedConfig.Database.Port = "5432"
	expectedConfig.Database.Host = "localhost"
	expectedConfig.Database.User = "postgres"
	expectedConfig.Database.Password = "Xxagrorog123"

	if *config != expectedConfig {
		t.Errorf("Excepted %+v, got %+v", expectedConfig, config)
	}

}
