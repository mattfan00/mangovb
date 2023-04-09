package api

import (
	"net/http"

	"github.com/mattfan00/mangovb/internal/store"
	"github.com/sirupsen/logrus"
)

const ADDR = ":8080"

type Api struct {
	logger     *logrus.Entry
	eventStore *store.EventStore
}

func New(logger *logrus.Entry, eventStore *store.EventStore) *Api {
	return &Api{
		logger:     logger,
		eventStore: eventStore,
	}
}

func (a *Api) Start() {
	a.logger.Info("starting server on ", ADDR)
	http.ListenAndServe(ADDR, a.routes())
}
