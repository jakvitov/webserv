package config

import (
	"cz/jakvitov/webserv/err"
	"encoding/json"
	"os"
)

// Given a path, return a ptr to a read and parsed Webserver config
// Return err if not found
func ReadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	res := &Config{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func verifyPorts(prts *Ports, e *err.ConfigParseError) {
	if prts.HttpPort == 0 && prts.HttpsPort == 0 {
		e.AddMissingField("Http/https port")
	}
}

// Verify the logger configuration
func verifyLogger(lg *Logger, e *err.ConfigParseError) {
	if lg.Level != WARN && lg.Level != INFO && lg.Level != ERROR && lg.Level != FATAL {
		e.AppendOrCreate("Logger level")
	}
}

// Verify the given config and return possible errors
// Else return nil
func verifyConfig(cnf *Config) *err.ConfigParseError {
	var e *err.ConfigParseError = nil
	verifyPorts(&cnf.Ports, e)
	verifyLogger(&cnf.Logger, e)
	return e
}

func ReadAndVerify(path string) (*Config, error) {
	cnf, err := ReadConfig(path)
	if err != nil {
		return nil, err
	}
	err = verifyConfig(cnf)
	if err != nil {
		return nil, err
	}
	return cnf, nil
}
