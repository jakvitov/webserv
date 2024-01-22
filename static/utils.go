package static

import (
	"bytes"
	"encoding/gob"
)

// Deep object copy of object
func DeepObjectCopy(original interface{}) interface{} {
	var result interface{}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	enc.Encode(original)
	buf.Reset()
	dec.Decode(&result)

	return result
}
