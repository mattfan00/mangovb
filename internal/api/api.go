package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const ADDR = ":8080"

func Start(logger *logrus.Entry) {
	logger.Info("starting server on ", ADDR)
	http.ListenAndServe(ADDR, routes())
}
