package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(env string, logPath string) (*logrus.Logger, error) {
	logger := logrus.New()

	if env == "prod" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(f)
	} else {
		logger.SetOutput(os.Stdout)
	}

	return logger, nil
}

func SetSource(logger *logrus.Logger, source string) *logrus.Entry {
	return logger.WithField("source", source)
}
