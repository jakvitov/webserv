//go:build !excludeTest

package server_test

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	echoserver "cz/jakvitov/webserv/test/echo_server"
	"net/http"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
)

const BASIC_SERVER_PATH string = "../test/config/minimal_config.yaml"
const REVERSE_PROXY_CONFIG string = "../test/config/reverse_proxy_config.yaml"
const REVERSE_PROXY_REGEX_CONFIG string = "../test/config/reverse_proxy_regex_config.yaml"

const URL string = "http://localhost:3000/"
const FORWARDED string = "X-Forwarded-For"

// Runs the server from the given config
func startServer(t *testing.T, path string, term *sync.WaitGroup) (*server.Server, *sync.WaitGroup) {
	cnf, err := config.ReadConfig(path)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf, term)
	wg := srv.StartListening()
	return srv, wg
}

func TestProxyForwarded(t *testing.T) {
	sdResourceSrv := make(chan bool, 1)
	sdProxySrv := make(chan bool, 1)
	startedSrvs := make(chan int, 2)

	terminatedWg := new(sync.WaitGroup)
	go func(term *sync.WaitGroup) {
		srv, startupWg := startServer(t, BASIC_SERVER_PATH, term)
		startupWg.Wait()
		startedSrvs <- 1
		<-sdResourceSrv
		srv.Shutdown()
	}(terminatedWg)

	go func(term *sync.WaitGroup) {
		srv, startupWg := startServer(t, REVERSE_PROXY_CONFIG, term)
		startupWg.Wait()
		startedSrvs <- 1
		<-sdProxySrv
		srv.Shutdown()
	}(terminatedWg)

	//Wait for both servers to start
	ready := 0
	for ready != 2 {
		sig := <-startedSrvs
		ready += sig
	}

	//Test the response
	res, err := http.Get(URL)
	assert.NilError(t, err)
	forwarded, found := res.Header[FORWARDED]
	assert.Equal(t, found, true)
	assert.Assert(t, len(forwarded) > 0)

	//Shutdown both servers
	sdResourceSrv <- true
	sdProxySrv <- true
	//Await both servers termination
	terminatedWg.Wait()
}

func TestProxyRegex(t *testing.T) {
	echoStartup := new(sync.WaitGroup)
	echoShutdown := make(chan bool, 1)
	echoShutdownCompleted := make(chan bool, 1)

	proxyShutdownWg := new(sync.WaitGroup)
	proxyShutdownSignal := make(chan bool, 1)
	proxyStartup := make(chan bool, 1)

	go func() {
		go echoserver.RunEchoServer(8080, echoStartup, echoShutdown, echoShutdownCompleted)
	}()

	go func() {
		srv, startupWg := startServer(t, REVERSE_PROXY_REGEX_CONFIG, proxyShutdownWg)
		startupWg.Wait()
		proxyStartup <- true
		<-proxyShutdownSignal
		srv.Shutdown()
	}()

	<-proxyStartup
	echoStartup.Wait()

	//Test the response that is not proxied
	res, err := http.Get("http://localhost:3000/djflkajd")
	assert.NilError(t, err)
	assert.Equal(t, res.StatusCode, 404)

	//This should proxy to echo server
	res, err = http.Get("http://localhost:3000/test/djdjdjdj")
	assert.NilError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	echoShutdown <- true
	<-echoShutdownCompleted
	proxyShutdownSignal <- true
	proxyShutdownWg.Wait()
}
