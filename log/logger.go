package log

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	gostLogger *log.Logger
	gostFile   *os.File
)

var (
	//ErrLoggerNotInitialized is thrown when uninitialized logger instance is requested
	ErrLoggerNotInitialized = errors.New("LoggerNotInitialized")
)

// GetLoggerInstance returns singleton instance of the logger
func GetLoggerInstance() (*log.Logger, error) {
	if gostLogger == nil {
		return log.New(), ErrLoggerNotInitialized
	}

	return gostLogger, nil
}

// InitializeLogger with various properties
func InitializeLogger(file *os.File, logFileName string, format log.Formatter, verboseFlag bool) (*log.Logger, error) {
	var err error
	gostLogger = log.New()
	gostFile = file

	if logFileName != "" {
		gostFile, err = os.OpenFile(logFileName+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Log file cannot be opened")
			gostFile.Close()
		}
	}

	if logFileName == "" || err != nil {
		gostFile = os.Stdout
	}

	gostLogger.Out = gostFile
	if format != nil {
		gostLogger.Formatter = format
	}

	if verboseFlag {
		gostLogger.Level = log.DebugLevel
	} else {
		gostLogger.Level = log.InfoLevel
	}

	return gostLogger, err
}

// DebugWithElapsedTime writes a new debug line, including a field with elapsed time
// call with defer at the start of a function: defer DebugWithElapsedTime(logger, time.Now(), "test")
func DebugWithElapsedTime(entry *log.Entry, start time.Time, args ...interface{}) {
	elapsed := time.Since(start)
	l := entry.WithFields(log.Fields{"elapsed": elapsed})
	l.Debug(args...)
}

// DebugfWithElapsedTime writes a new debug format line, including a field with elapsed time
// call with defer at the start of a function: defer DebugWithElapsedTime(logger, time.Now(), "test %v", "test")
func DebugfWithElapsedTime(entry *log.Entry, start time.Time, format string, args ...interface{}) {
	elapsed := time.Since(start)
	l := entry.WithFields(log.Fields{"elapsed": elapsed})
	l.Debugf(format, args...)
}

// CleanUp after the logger is closed
func CleanUp() {
	if gostFile != nil {
		gostFile.Close()
	}
}
