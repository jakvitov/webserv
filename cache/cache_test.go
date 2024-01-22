package cache

import (
	"cz/jakvitov/webserv/sharedlogger"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

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

func testCacheBenchmark(b *testing.B) {

}
