package service

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type LogRorate interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type logRorate struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewLogrorate() LogRorate {
	file, err := os.OpenFile("./log/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &logRorate{
		infoLogger:    InfoLogger,
		warningLogger: WarningLogger,
		errorLogger:   ErrorLogger,
	}
}

// Error implements LogRorate
func (l *logRorate) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

// Info implements LogRorate
func (l *logRorate) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

// Warn implements LogRorate
func (l *logRorate) Warn(v ...interface{}) {
	l.warningLogger.Println(v...)
}
