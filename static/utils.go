package static

import (
	"bytes"
	"crypto/tls"
	"cz/jakvitov/webserv/sharedlogger"
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

// Loads a TLS cerfificate from a given cert and key files
func LoadTlsCerfificate(cert, key string, lg *sharedlogger.SharedLogger) (*tls.Certificate, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		lg.Ferror("Error loading TLS certificate [%s]\n", err.Error())
		return nil, err
	}
	return &certificate, nil
}
