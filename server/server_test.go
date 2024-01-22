package server_test

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

// Serves only index simple webpage
const LOCALHOST_URL string = "http://localhost:8080/"

const ONLY_INDEX_FILE string = "../test/web_content/only_index_webpage/index.html"

const CORRECT_CONFIG string = "../test/config/correct_simple_config.yaml"

func TestServerInit(t *testing.T) {
	wg := new(sync.WaitGroup)
	cnf, err := config.ReadConfig(CORRECT_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf)
	srv.StartListening(wg)

	time.Sleep(50 * time.Microsecond)

	//We give the server 200 ms to initialize
	t.Logf("Sending a get request.\n")
	res, err := http.Get(LOCALHOST_URL)
	assert.NilError(t, err)
	defer res.Body.Close()
	assert.Check(t, res != nil)
	srv.Shutdown()
	wg.Wait()
}

// Test serving of the only index webpage
func TestServerOnlyIndexWebpage(t *testing.T) {
	wg := new(sync.WaitGroup)
	cnf, err := config.ReadConfig(CORRECT_CONFIG)
	assert.NilError(t, err)
	srv := server.ServerInit(cnf)
	srv.StartListening(wg)
	//Let the server load properly
	time.Sleep(50 * time.Microsecond)

	res, err := http.Get(LOCALHOST_URL)
	assert.NilError(t, err)
	expected, err := os.ReadFile(ONLY_INDEX_FILE)
	assert.NilError(t, err)
	body, err := io.ReadAll(res.Body)
	//defer res.Body.Close()
	assert.NilError(t, err)
	assert.DeepEqual(t, body, expected)
	srv.Shutdown()
	wg.Wait()
}

func BenchmarkCacheOneFile(b *testing.B) {
	wg := new(sync.WaitGroup)
	cnf, err := config.ReadAndVerify("../test/config/minimal_config.yaml")
	assert.NilError(b, err)
	index, err := os.ReadFile("../test/web_content/only_index_webpage/index.html")
	assert.NilError(b, err)
	srv := server.ServerInit(cnf)
	srv.StartListening(wg)
	time.Sleep(50 * time.Microsecond)
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
	wg.Wait()
}
