package logger

import (
	"log"
)

type PrintfLogger struct {
	logger *log.Logger
}

func NewPrintfLogger() *PrintfLogger {
	return &PrintfLogger{logger: log.Default()}
}

func (p *PrintfLogger) Debug(msg string, args ...any) { p.logger.Printf("DEBUG: "+msg, args...) }
func (p *PrintfLogger) Info(msg string, args ...any)  { p.logger.Printf("INFO: "+msg, args...) }
func (p *PrintfLogger) Warn(msg string, args ...any)  { p.logger.Printf("WARN: "+msg, args...) }
func (p *PrintfLogger) Error(msg string, args ...any) { p.logger.Printf("ERROR: "+msg, args...) }
func (p *PrintfLogger) Fatal(msg string, args ...any) { p.logger.Fatalf("FATAL: "+msg, args...) }
