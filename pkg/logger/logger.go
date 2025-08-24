package logger

import (
	"log"
	"os"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type logger struct {
	stdLog *log.Logger
}

func NewLogger() *logger {
	return &logger{
		stdLog: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Infof logs an informational message to the standard output
func (logger *logger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// Fatalf logs a fatal error message to the standard output
func (logger *logger) Fatalf(format string, args ...interface{}) {
	log.Printf("[FATAL] "+format, args...)
}

// Errorf logs an error message to the standard output
func (logger *logger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// Debugf logs a debug message to the standard output
func (logger *logger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}
