package server_test

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	"net/http"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const BASIC_SERVER_PATH string = "../test/config/minimal_config.yaml"
const REVERSE_PROXY_CONFIG string = "../test/config/reverse_proxy_config.yaml"
const URL string = "http://localhost:3000/"
const FORWARDED string = "X-Forwarded-For"

// Runs the server from the given config
func startServer(t *testing.T, path string) (*server.Server, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	cnf, err := config.ReadConfig(path)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf)
	srv.StartListening(wg)
	//Let the server load properly
	time.Sleep(50 * time.Microsecond)
	return srv, wg
}

// Test the reverse proxy functionality
func TestProxy(t *testing.T) {
	//Start server
	basicSrv, basicWg := startServer(t, BASIC_SERVER_PATH)
	proxySrv, proxyWg := startServer(t, REVERSE_PROXY_CONFIG)
	defer func() {
		basicSrv.Shutdown()
		proxySrv.Shutdown()
		basicWg.Wait()
		proxyWg.Wait()
	}()

	res, err := http.Get(LOCALHOST_URL)
	assert.NilError(t, err)
	forwarded, found := res.Header[FORWARDED]
	assert.Equal(t, found, true)
	assert.Assert(t, len(forwarded) > 0)

}
