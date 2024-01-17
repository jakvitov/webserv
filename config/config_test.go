package config_test

import (
	"cz/jakvitov/webserv/config"
	"testing"

	"gotest.tools/v3/assert"
)

const inputCorrectConfigPath string = "../test/config/correct_simple_config.yaml"
const inputErrorConfigPath string = "../test/config/error_config.yaml"
const minimalConfigPath string = "../test/config/minimal_config.yaml"

func TestReadConfigCoorect(t *testing.T) {
	conf, err := config.ReadConfig(inputCorrectConfigPath)
	assert.NilError(t, err)
	assert.Check(t, conf.Logger.Level == "INFO")
	assert.Check(t, conf.Ports.HttpPort == 8080)
}

func TestDefaultConfigCreation(t *testing.T) {
	conf, err := config.ReadConfig(minimalConfigPath)
	assert.NilError(t, err)
	assert.Check(t, conf.Logger.Level == "INFO")
	assert.Check(t, conf.Ports.HttpPort == 8080)
	assert.Check(t, conf.Logger.AppendOutput == false)
}

func TestReadWrongConfig(t *testing.T) {
	conf, err := config.ReadConfig(inputErrorConfigPath)
	assert.Check(t, err != nil)
	assert.Check(t, conf == nil)
}
