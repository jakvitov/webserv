package config

import (
	"os"

	"gopkg.in/yaml.v2"
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
	HttpPort  int `yaml:"http_port"`
	HttpsPort int `yaml:"https_port"`
}

type Logger struct {
	Level        string `yaml:"level"`
	OutputToFile bool   `yaml:"output_to_file"`
	OutputFile   string `yaml:"output_file"`
	AppendOutput bool   `yaml:"append_output"`
}

type Handler struct {
	ContentRoot    string `yaml:"content_root"`
	ReadTimeout    int    `yaml:"read_timeout"`
	WriteTimeout   int    `yaml:"write_timeout"`
	MaxHeaderBytes int    `yaml:"max_header_bytes"`
}

type Route struct {
	From string `yaml:"from"`
	To   int    `yaml:"to"`
}

type ReverseProxy struct {
	Routes []Route `yaml:"routes"`
}

type Security struct {
	CertPath   string `yaml:"cert_path"`
	SpamFilter bool   `yaml:"spam_filter"`
}

type Config struct {
	Ports        Ports        `yaml:"ports"`
	Logger       Logger       `yaml:"logger"`
	Handler      Handler      `yaml:"handler"`
	ReverseProxy ReverseProxy `yaml:"reverse_proxy"`
	Security     Security     `yaml:"security"`
}

// Function used to parse input config into a data structure
func DecodeConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &Config{}
	if err := yaml.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
