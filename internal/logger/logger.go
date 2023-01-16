package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	return logger
}

func SetSource(logger *logrus.Logger, source string) *logrus.Entry {
	return logger.WithField("source", source)
}
