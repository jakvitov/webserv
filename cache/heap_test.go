package cache

import (
	"testing"

	"gotest.tools/v3/assert"
)

// Test heap push pops and rebalancing with loging cached files
func TestMinHeap(t *testing.T) {
	cf1, err := CachedFileInit("../main.go")
	cf2, err2 := CachedFileInit("./heap_test.go")
	cf3, err3 := CachedFileInit("./heap.go")

	assert.NilError(t, err)
	assert.NilError(t, err2)
	assert.NilError(t, err3)
	cf1.read = 10
	cf2.read = 5
	cf3.read = 1

	fh := &CachedFileHeap{}
	//heap.Init(fh)
	fh.Push(cf1)
	fh.Push(cf2)
	fh.Push(cf3)

	res := fh.Pop().(*CachedFile)
	//We remove the least coefficients first - the least used files in the cache
	assert.Equal(t, res.read, int64(1))
}
