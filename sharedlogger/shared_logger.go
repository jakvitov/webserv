package sharedlogger

import (
	"cz/jakvitov/webserv/config"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const ERROR_PREFIX string = "[ERROR]"
const INFO_PREFIX string = "[INFO]"
const WARN_PREFIX string = "[WARN]"
const FATAL_PREFIX string = "[FATAL]"

/*
This implementation build on top of golang logger is designer to be shared as one struct instance
and allow logging trough simple methods
*/
type SharedLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	level       string
}

// Return shared logger used in tests, that just prints to the standard output
func TestingSharedLoggerInit() *SharedLogger {
	return &SharedLogger{
		infoLogger:  log.New(io.MultiWriter(os.Stdout), "", 0),
		errorLogger: log.New(io.MultiWriter(os.Stderr), "", 0),
		level:       "INFO",
	}
}

// The output file streams are to print out to multiple files as well as stdio
func SharedLoggerInit(outputFileStream *os.File, cnf *config.Config) *SharedLogger {
	res := &SharedLogger{}
	if outputFileStream == nil {
		res.infoLogger = log.New(io.MultiWriter(os.Stdout), "", 0)
		res.errorLogger = log.New(io.MultiWriter(os.Stderr), "", 0)
	} else {
		res.infoLogger = log.New(io.MultiWriter(os.Stdout, outputFileStream), "", 0)
		res.errorLogger = log.New(io.MultiWriter(os.Stderr, outputFileStream), "", 0)
	}
	res.level = cnf.Logger.Level
	return res
}

func getDateTimeLogPrefix() string {
	//todo customize with config
	return fmt.Sprintf("[%s]", time.Now().Format(time.DateTime))
}

// Simple print to stdout and other streams
func (s *SharedLogger) Info(message string) {
	if s.level == config.INFO {
		s.infoLogger.Printf("%s;%s;%s\n", INFO_PREFIX, getDateTimeLogPrefix(), message)
	}
}

func (s *SharedLogger) Finfo(input string, args ...any) {
	if s.level == config.INFO {
		s.infoLogger.Printf("%s;%s;%s\n", INFO_PREFIX, getDateTimeLogPrefix(), fmt.Sprintf(input, args...))
	}
}

func (s *SharedLogger) Warn(message string) {
	if s.level == config.INFO || s.level == config.WARN {
		s.infoLogger.Printf("%s;%s;%s\n", WARN_PREFIX, getDateTimeLogPrefix(), message)
	}
}

// Not followed by any other action than log
func (s *SharedLogger) Error(message string) {
	if s.level == config.INFO || s.level == config.WARN || s.level == config.ERROR {
		s.errorLogger.Printf("%s;%s;%s\n", ERROR_PREFIX, getDateTimeLogPrefix(), message)
	}
}

// Followed by os.Exit(1) syscall call
func (s *SharedLogger) Fatal(message string) {
	s.errorLogger.Printf("%s;%s;%s\n", FATAL_PREFIX, getDateTimeLogPrefix(), message)
	os.Exit(1)
}
