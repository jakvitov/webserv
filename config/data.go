package config

import (
	"encoding/json"
	"os"
)

// Constants for logger levels
const (
	INFO  string = "INFO"
	WARN  string = "WARN"
	ERROR string = "ERROR"
	FATAL string = "FATAL"
)

/*This file includes the data structures used in the config file and to parse them*/

type Ports struct {
	HttpPort  int `json:"http_port"`
	HttpsPort int `json:"https_port"`
}

type Logger struct {
	Level        string `json:"level"`
	OutputStd    bool   `json:"output_std"`
	OutputFile   string `json:"output_file"`
	AppendOutput bool   `json:"append_output"`
}

type Handler struct {
	ContentRoot    string `json:"content_root"`
	ReadTimeout    int    `json:"read_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
	MaxHeaderBytes int    `json:"max_header_bytes"`
}

type Route struct {
	From string `json:"from"`
	To   int    `json:"to"`
}

type ReverseProxy struct {
	Routes []Route `json:"routes"`
}

type Security struct {
	CertPath   string `json:"cert_path"`
	SpamFilter bool   `json:"spam_filter"`
}

type Config struct {
	Ports        Ports        `json:"ports"`
	Logger       Logger       `json:"logger"`
	Handler      Handler      `json:"handler"`
	ReverseProxy ReverseProxy `json:"reverse_proxy"`
	Security     Security     `json:"security"`
}

// Function used to parse input config into a data structure
func DecodeConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &Config{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
