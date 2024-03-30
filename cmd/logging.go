package commands

import (
	"log"
	"os"
)

func init() {
	Log = *NewLogging()
}

var Log Logging

type Logging struct {
	Stdout     *log.Logger
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	DebugLog   *log.Logger
	SessionLog *log.Logger
}

func openLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

func NewLogging() *Logging {
	stdout := log.New(os.Stdout, "", 0)
	infoLog := log.New(openLogFile("info.log"), "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog := log.New(openLogFile("debug.log"), "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	sessionLog := log.New(openLogFile("session.log"), "\t", 0)

	app := &Logging{
		Stdout:     stdout,
		ErrorLog:   errorLog,
		InfoLog:    infoLog,
		DebugLog:   debugLog,
		SessionLog: sessionLog,
	}

	return app
}
