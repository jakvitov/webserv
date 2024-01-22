package cache

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/server"
	"cz/jakvitov/webserv/sharedlogger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const LOCALHOST_URL string = "http://localhost:8080/"

func TestCacheEnabled(t *testing.T) {
	//Read the whole test package and serve it in limit of 2kB of cache -> causes often rebalancing
	cache := CacheInit(2000, true, sharedlogger.TestingSharedLoggerInit())
	err := filepath.WalkDir("../test/", func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {
			data, err := cache.Get(path)
			assert.NilError(t, err)
			comp, err := os.ReadFile(path)
			assert.NilError(t, err)
			assert.DeepEqual(t, comp, data)
		}
		return nil
	})
	assert.NilError(t, err)
}

func TestCacheDisabled(t *testing.T) {
	//Read the whole test package and serve it in limit of 2kB of cache -> causes often rebalancing
	cache := CacheInit(2000, false, sharedlogger.TestingSharedLoggerInit())
	err := filepath.WalkDir("../test/", func(path string, info os.DirEntry, err error) error {
		if !info.IsDir() {
			data, err := cache.Get(path)
			assert.NilError(t, err)
			comp, err := os.ReadFile(path)
			assert.NilError(t, err)
			assert.DeepEqual(t, comp, data)
		}
		return nil
	})
	assert.NilError(t, err)
}

func BenchmarkCacheOneFile(b *testing.B) {
	wg := new(sync.WaitGroup)
	cnf, err := config.ReadConfig("../test/config/minimal_config.yaml")
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
