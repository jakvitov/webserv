package cache

import (
	"cz/jakvitov/webserv/sharedlogger"
	"os"
	"sync"
)

// In memmory cache to hold data, so they do not have to be
// loaded from the files
type Cache struct {
	files     map[string]*CachedFile
	maxBytes  int64
	totalSize int64
	//Lock for rebalance
	mutex   sync.Mutex
	locked  bool
	logger  *sharedlogger.SharedLogger
	enabled bool
}

// Create new empty cache with max size
func CacheInit(max int64, enabled bool, logger *sharedlogger.SharedLogger) *Cache {
	return &Cache{
		files:     make(map[string]*CachedFile),
		maxBytes:  max,
		totalSize: 0,
		logger:    logger,
		locked:    false,
		enabled:   enabled,
	}
}

func (c *Cache) encache(cf *CachedFile) {
	c.logger.Finfo("Cached file [%s]\n", cf.GetPath())
	c.totalSize += cf.GetSize()
	c.files[cf.GetPath()] = cf
}

func (c *Cache) addOrRebalance(cf *CachedFile) {
	c.mutex.Lock()
	c.locked = true
	defer func() { c.locked = false }()
	defer c.mutex.Unlock()

	//We have free space to accomoddate this file
	if cf.GetSize() <= c.maxBytes-c.totalSize {
		c.encache(cf)
		return
	}

	c.logger.Finfo("Rebalancing cache to add [%s]\n", cf.path)
	//We do not have space to accomodate this file and need to remove the less used ones
	//Min heap based on the least used files (by our coefficient)
	fh := &CachedFileHeap{}
	for _, val := range c.files {
		fh.Push(val)
	}

	//We clear the files from the least used until we make room for the new one
	for c.maxBytes-c.totalSize >= cf.GetSize() {
		toRem := fh.Pop().(*CachedFile)
		delete(c.files, toRem.path)
		c.totalSize -= toRem.GetSize()
	}

	//Now we have made the space required for the new file
	c.encache(cf)
}

// Add a file to cache if he fits or free space for him in a separate thread
func (c *Cache) Get(path string) ([]byte, error) {
	//We are rebalancing in another goroutine or we do not use cache
	if c.locked || !c.enabled {
		res, err := os.ReadFile(path)
		return res, err
	}

	file, found := c.files[path]
	//We cannot cache same file twice
	if found {
		return file.GetData(), nil
	}

	cf, err := CachedFileInit(path)
	if err != nil {
		c.logger.Warn("Error while opening file: " + err.Error())
		return nil, err
	}
	//The file is too large to be cached, we just return it and do not cache
	if cf.GetSize() > c.maxBytes {
		return cf.GetData(), nil
	}
	go c.addOrRebalance(cf)
	return cf.GetData(), nil
}
