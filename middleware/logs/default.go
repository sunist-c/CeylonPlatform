package logs

import (
	"fmt"
	"io"
	"log"
)

func Default(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		logger: log.New(out, prefix, flag),
	}
}

type Logger struct {
	logger *log.Logger
}

func (d Logger) Info(logText string) {
	d.logger.Println(fmt.Sprintf("(info): %v", logText))
}

func (d Logger) Warning(logText string) {
	d.logger.Println(fmt.Sprintf("(warning): %v", logText))
}

func (d Logger) Error(logText string) {
	d.logger.Println(fmt.Sprintf("(error): %v", logText))
}

func (d Logger) Fatal(logText string) {
	d.logger.Fatalln(fmt.Sprintf("(fatal): %v", logText))
}
