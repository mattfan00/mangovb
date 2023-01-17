package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(env string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	if env == "prod" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	return logger
}

func SetSource(logger *logrus.Logger, source string) *logrus.Entry {
	return logger.WithField("source", source)
}
