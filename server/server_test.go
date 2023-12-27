package server_test

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	"testing"

	"gotest.tools/v3/assert"
)

const CORRECT_CONFIG string = "../test/config/config_correct.json"

func TestServerInit(t *testing.T) {
	cnf, err := config.ReadConfig(CORRECT_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf)
	srv.StartListening()
}
