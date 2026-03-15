package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func GetLogger() *logrus.Logger {
	// logger := logrus.New()
	file, err := os.OpenFile("./logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	multipleWriter := io.MultiWriter(file, os.Stdout)
	logger.SetOutput(multipleWriter)
	logger.SetLevel(logrus.DebugLevel)
	// defer file.Close()
	return logger
}
