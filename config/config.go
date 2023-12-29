package config

import (
	"cz/jakvitov/webserv/err"
	"encoding/json"
	"os"
)

// This file includes the configuration for this HTTP server, that is given as input
type WebserverConfig struct {
	Ports   []int  `json:"ports"`
	RootDir string `json:"content_root"`
	LogDest string `json:"log_dest"`
}

// Check for mandatory fields - non empty
func checkMandatoryFields(wc *WebserverConfig) *err.ConfigParseError {
	mf := make([]string, 0)
	if len(wc.Ports) == 0 {
		mf = append(mf, "ports")
	}
	if len(wc.LogDest) == 0 {
		mf = append(mf, "log_dest")
	}
	if len(wc.RootDir) == 0 {
		mf = append(mf, "content_root")
	}
	if len(mf) == 0 {
		return nil
	}
	return err.ConfigParseErrorInit(mf)
}

// Given a path, return a ptr to a read and parsed Webserver config
// Return err if not found
func ReadConfig(filePath string) (*WebserverConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	res := &WebserverConfig{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if err := checkMandatoryFields(res); err != nil {
		return nil, err
	}
	return res, nil
}
