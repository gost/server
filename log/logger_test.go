package log

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

var (
	testFile        *os.File
	testLogger      *log.Logger
	err             error
	testLogFileName = "testLog"
	testVerboseFlag = true
)

func TestGetInstanceWithoutInitialization(t *testing.T) {
	testLogger, err := GetLoggerInstance()

	assert.Error(t, err, "Function must throw error")
	assert.EqualError(t, err, ErrLoggerNotInitialized.Error(), "Error should be logger not initialized")

	assert.NotNil(t, testLogger, "Logger should not be nil")
}

func TestLoggerLifecycle(t *testing.T) {
	//Setup
	testLogger, err = InitializeLogger(testFile, testLogFileName, new(log.TextFormatter), testVerboseFlag)
	assert.NoError(t, err, "Initialization error should be nil")
	assert.NotNil(t, testLogger, "Logger must have been initialized")

	//GetInstance
	loggerInstance, err := GetLoggerInstance()
	assert.NoError(t, err, "GetInstance should not return error")
	assert.NotNil(t, loggerInstance, "Logger instance should not be nil")
	loggerInstance.Debug("Debug")
	loggerInstance.Info("Info")
	loggerInstance.Warn("Warn")
	loggerInstance.Error("Error")
	_, err = os.Stat("./" + testLogFileName + ".log")
	assert.False(t, os.IsNotExist(err), "Test log file should exist")

	//Cleanup
	CleanUp()
	assert.Nil(t, testFile, "Test file should be Nil")
}

func TestCleanUpWithoutInitialization(t *testing.T) {
	assert.NotPanics(t, func() { CleanUp() }, "Must not panic")
}

func TestDebugWithElapsedTime(t *testing.T) {
	testLogger, err = InitializeLogger(testFile, "", new(log.TextFormatter), true)
	entry := testLogger.WithFields(log.Fields{"package": "gost.server.log"})

	f := func() { DebugWithElapsedTime(entry, time.Now(), "test") }
	message := captureStdout(f)

	//assert
	assert.Contains(t, message, "elapsed")
	assert.Contains(t, message, "test")
}

func TestDebugfWithElapsedTime(t *testing.T) {
	testLogger, err = InitializeLogger(testFile, "", new(log.TextFormatter), true)
	entry := testLogger.WithFields(log.Fields{"package": "gost.server.log"})

	f := func() { DebugfWithElapsedTime(entry, time.Now(), "test %s", "1") }
	message := captureStdout(f)

	//assert
	assert.Contains(t, message, "elapsed")
	assert.Contains(t, message, "test 1")
}

func captureStdout(f func()) string {
	old := testLogger.Out
	r, w, _ := os.Pipe()
	testLogger.Out = w

	f()

	w.Close()
	testLogger.Out = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
