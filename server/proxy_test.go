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
	cnf, err := config.ReadConfig(path)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf)
	wg := srv.StartListening()
	//Let the server load properly
	time.Sleep(100 * time.Millisecond)
	return srv, wg
}

// Test the reverse proxy functionality
func TestProxy(t *testing.T) {
	//Channels for concurrency orchestration
	sdnormal := make(chan bool, 1)
	sdproxy := make(chan bool, 1)
	waitchan := make(chan int, 2)
	//Start server
	go func(sdchan chan bool, waitchan chan int) {
		basicSrv, basicWg := startServer(t, BASIC_SERVER_PATH)
		basicWg.Wait()
		waitchan <- 1
		shutdown := <-sdchan
		if shutdown {
			basicSrv.Shutdown()
			waitchan <- 1
		}
	}(sdnormal, waitchan)

	go func(sdchan chan bool, waitchan chan int) {
		proxySrv, proxyWg := startServer(t, REVERSE_PROXY_CONFIG)
		proxyWg.Wait()
		waitchan <- 1
		shutdown := <-sdchan
		if shutdown {
			proxySrv.Shutdown()
			waitchan <- 1
		}
	}(sdproxy, waitchan)

	//Wait for both servers to start
	ready := 0
	for ready != 2 {
		sig := <-waitchan
		ready += sig
	}

	//Test the response
	res, err := http.Get(URL)
	assert.NilError(t, err)
	forwarded, found := res.Header[FORWARDED]
	assert.Equal(t, found, true)
	assert.Assert(t, len(forwarded) > 0)

	//Shutdown the
	sdnormal <- true
	sdproxy <- true
	ready = 0
	for ready != 2 {
		sig := <-waitchan
		ready += sig
	}

}
