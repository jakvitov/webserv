package config_test

import (
	"cz/jakvitov/webserv/config"
	"testing"

	"gotest.tools/v3/assert"
)

const inputCorrectConfigPath string = "../test/config/correct_simple_config.yaml"

func TestReadConfigCoorect(t *testing.T) {
	inputConfigPath := inputCorrectConfigPath
	conf, err := config.ReadConfig(inputConfigPath)
	assert.NilError(t, err)
	assert.Check(t, conf.Logger.Level == "INFO")
	assert.Check(t, conf.Ports.HttpPort == 80)
}
