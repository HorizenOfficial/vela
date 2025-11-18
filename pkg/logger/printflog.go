package logger

import (
	"io"
	"log"
	"os"
)

type PrintfLogger struct {
	logger *log.Logger
}

func NewPrintfLogger(cfg *Config) *PrintfLogger {
	writers := []io.Writer{}

	if cfg.Console {
		writers = append(writers, os.Stderr)
	}

	if cfg.FileName != "" {
		logFile, err := os.OpenFile(cfg.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		writers = append(writers, logFile)
	}

	var writer io.Writer
	if len(writers) > 0 {
		writer = io.MultiWriter(writers...)
	} else {
		// default to stderr if no output is specified
		writer = os.Stderr
	}

	return &PrintfLogger{logger: log.New(writer, "", log.LstdFlags)}
}

func (p *PrintfLogger) Trace(msg string, args ...any) { p.logger.Printf("TRC: "+msg, args...) }
func (p *PrintfLogger) Debug(msg string, args ...any) { p.logger.Printf("DBG: "+msg, args...) }
func (p *PrintfLogger) Info(msg string, args ...any)  { p.logger.Printf("INF: "+msg, args...) }
func (p *PrintfLogger) Warn(msg string, args ...any)  { p.logger.Printf("WRN: "+msg, args...) }
func (p *PrintfLogger) Error(msg string, args ...any) { p.logger.Printf("ERR: "+msg, args...) }
func (p *PrintfLogger) Fatal(msg string, args ...any) { p.logger.Fatalf("FTL: "+msg, args...) }
func (p *PrintfLogger) Panic(msg string, args ...any) { p.logger.Panicf("PNC: "+msg, args...) }
