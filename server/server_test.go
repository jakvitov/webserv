package server_test

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const CORRECT_CONFIG string = "../test/config/config_correct.json"

func TestServerInit(t *testing.T) {
	cnf, err := config.ReadConfig(CORRECT_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf)
	srv.StartListening(new(sync.WaitGroup))

	time.Sleep(50 * time.Microsecond)

	//We give the server 200 ms to initialize
	t.Logf("Sending a get request.\n")
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", cnf.Ports[0]))
	assert.NilError(t, err)
	assert.Check(t, res != nil)

}
