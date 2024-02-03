//go:build !excludeTest

package server_test

import (
	"crypto/tls"
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
)

// Serves only index simple webpage
const LOCALHOST_URL string = "http://localhost:8080/"
const HTTPS_LOCALHOST string = "https://localhost:8443"

const ONLY_INDEX_FILE string = "../test/web_content/only_index_webpage/index.html"

const CORRECT_CONFIG string = "../test/config/correct_simple_config.yaml"

const HTTPS_CONFIG string = "../test/config/https_server.yaml"

func TestServerInit(t *testing.T) {
	term := new(sync.WaitGroup)
	cnf, err := config.ReadConfig(CORRECT_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf, term)
	wg := srv.StartListening()

	wg.Wait()

	//We give the server 200 ms to initialize
	t.Logf("Sending a get request.\n")
	res, err := http.Get(LOCALHOST_URL)
	assert.NilError(t, err)
	defer res.Body.Close()
	assert.Check(t, res != nil)
	srv.Shutdown()
	term.Wait()
}

// Test serving of the only index webpage
func TestServerOnlyIndexWebpage(t *testing.T) {
	term := new(sync.WaitGroup)
	cnf, err := config.ReadConfig(CORRECT_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf, term)
	wg := srv.StartListening()
	//Let the server load properly
	wg.Wait()

	res, err := http.Get(LOCALHOST_URL)
	assert.NilError(t, err)
	expected, err := os.ReadFile(ONLY_INDEX_FILE)
	assert.NilError(t, err)
	body, err := io.ReadAll(res.Body)
	//defer res.Body.Close()
	assert.NilError(t, err)
	assert.DeepEqual(t, body, expected)
	srv.Shutdown()
	term.Wait()
}

func BenchmarkCacheOneFile(b *testing.B) {
	term := new(sync.WaitGroup)
	cnf, err := config.ReadAndVerify("../test/config/minimal_config.yaml")
	assert.NilError(b, err)
	index, err := os.ReadFile("../test/web_content/only_index_webpage/index.html")
	assert.NilError(b, err)
	srv := server.ServerInit(cnf, term)
	wg := srv.StartListening()
	wg.Wait()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		res, err := http.Get(LOCALHOST_URL)
		b.StopTimer()
		assert.NilError(b, err)
		resData, err := io.ReadAll(res.Body)
		assert.NilError(b, err)
		assert.DeepEqual(b, resData, index)
	}
	srv.Shutdown()
	term.Wait()
}

// Create HTTP client, that accepts unknown Certificate authorities
func CreateTestApliClient() *http.Client {
	c := &http.Client{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c = &http.Client{Transport: tr}
	return c
}

// Test for simple https connection
func TestHttpsServer(t *testing.T) {
	client := CreateTestApliClient()
	term := new(sync.WaitGroup)
	cnf, err := config.ReadConfig(HTTPS_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf, term)
	wg := srv.StartListening()

	//Let the server load properly
	wg.Wait()

	t.Logf("Sending a http GET request.\n")
	res, err := client.Get(LOCALHOST_URL)
	assert.NilError(t, err)
	assert.Check(t, res != nil)

	t.Logf("Sending a https GET request.\n")
	res, err = client.Get(HTTPS_LOCALHOST)
	assert.NilError(t, err)
	assert.Check(t, res != nil)
	assert.Equal(t, res.StatusCode, 200)

	srv.Shutdown()
	term.Wait()
	client.CloseIdleConnections()
}
