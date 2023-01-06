package logger

import "github.com/sirupsen/logrus"

func New() *logrus.Logger {
	logger := logrus.New()

	return logger
}

func SetSource(logger *logrus.Logger, source string) *logrus.Entry {
	return logger.WithField("source", source)
}
