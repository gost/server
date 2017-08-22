package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestGetInstanceWithoutInitialization(t *testing.T) {
	testLogger, err := GetInstance()

	assert.Error(t, err, "Function must throw error")
	assert.EqualError(t, err, ErrLoggerNotInitialized.Error(), "Error should be logger not initialized")

	assert.NotNil(t, testLogger, "Logger should not be nil")
}

func TestLoggerLifecycle(t *testing.T) {

	//Setup
	var (
		testFile        *os.File
		testLogger      *log.Logger
		testLogFileName = "testLog"
		testVerboseFlag = true
	)

	//Initialize
	testLogger, err := InitializeLogger(testFile, testLogFileName, new(log.TextFormatter), testVerboseFlag)
	assert.NoError(t, err, "Initialization error should be nil")
	assert.NotNil(t, testLogger, "Logger must have been initialized")

	//GetInstance
	loggerInstance, err := GetInstance()
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
