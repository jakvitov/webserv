package cache

import (
	"os"
	"sync"
	"time"
)

// One particular file loaded in a memmory
type CachedFile struct {
	data    []byte
	created time.Time
	read    int64
	path    string
	mutex   sync.RWMutex
}

// Read a cached file from a memory
func CachedFileInit(path string) (*CachedFile, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &CachedFile{
		data:    bytes,
		created: time.Now(),
		read:    int64(0),
		path:    path,
	}, nil
}

// Get coefficient of usage rate of the file
// Lower the coefficient, the less used the file is
func (c *CachedFile) GetCoefficient() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return (c.created.Unix()) / (c.read * c.read)
}

// Get file data
func (c *CachedFile) GetData() []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.read += 1
	return c.data
}

func (c *CachedFile) GetSize() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return int64(len(c.data))
}

func (c *CachedFile) GetPath() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.path
}
