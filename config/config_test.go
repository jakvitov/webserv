package config_test

import (
	"cz/jakvitov/webserv/config"
	"testing"

	"gotest.tools/v3/assert"
)

const inputCorrectConfigPath string = "../test/config/config_correct.json"
const inputMissingConfigPath string = "../test/config/config_incorrect.json"

func TestReadConfigCoorect(t *testing.T) {
	inputConfigPath := inputCorrectConfigPath
	conf, err := config.ReadConfig(inputConfigPath)
	assert.NilError(t, err)
	assert.Check(t, len(conf.LogDest) != 0)
	assert.Check(t, len(conf.Ports) != 0)
	assert.Check(t, len(conf.RootDir) != 0)
}

func TestReadConfigIncorrect(t *testing.T) {
	_, err := config.ReadConfig(inputMissingConfigPath)
	assert.Check(t, err != nil)
}
