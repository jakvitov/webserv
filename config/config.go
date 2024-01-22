package config

import (
	"cz/jakvitov/webserv/err"
	"os"

	"gopkg.in/yaml.v2"
)

// Defaults for missing values in the config
const READ_TIMEOUT_MS_DEFAULT int = 1000
const WRITE_TIMEOUT_MS_DEFAULT int = 1000
const MAX_HEADER_BYTES_DEFAULT int = 1 << 20

const DEFAULT_HTTP_PORT int = 8080
const DEFAULT_LOG_LEVEL string = "INFO"
const OUTPUT_TO_FILE_DEFAULT bool = true
const OUTPUT_FILE_DEFAULT string = "./websersv_log.log"
const APPEND_OUTPUT_DEFAULT bool = false

// Return pointer to a preset config with default values
func getDefaultConfig() *Config {
	res := &Config{
		Ports: Ports{
			HttpPort: DEFAULT_HTTP_PORT,
		},
		Logger: Logger{
			Level:        DEFAULT_LOG_LEVEL,
			OutputToFile: OUTPUT_TO_FILE_DEFAULT,
			OutputFile:   OUTPUT_FILE_DEFAULT,
			AppendOutput: APPEND_OUTPUT_DEFAULT,
		},
	}
	return res
}

// Given a path, return a ptr to a read and parsed Webserver config
// Return err if not found
func ReadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	//Default config values
	res := getDefaultConfig()
	if err := yaml.Unmarshal(data, &res); err != nil {
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
func verifyConfig(cnf *Config) error {
	var e *err.ConfigParseError = nil
	verifyPorts(&cnf.Ports, e)
	verifyLogger(&cnf.Logger, e)
	verifyHanlder(&cnf.Handler)
	if e == nil {
		return nil
	}
	return e
}

func verifyHanlder(hd *Handler) {
	if hd.MaxHeaderBytes == 0 {
		hd.MaxHeaderBytes = MAX_HEADER_BYTES_DEFAULT
	}
	if hd.ReadTimeout == 0 {
		hd.ReadTimeout = READ_TIMEOUT_MS_DEFAULT
	}
	if hd.WriteTimeout == 0 {
		hd.WriteTimeout = WRITE_TIMEOUT_MS_DEFAULT
	}
	if hd.CacheEnabled == true && hd.MaxCacheSize == 0 {
		hd.MaxCacheSize = 100000 * 1000
	}
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
