package api

import (
	"fmt"
	"net/http"

	"github.com/mattfan00/mangovb/internal/store"

	"github.com/sirupsen/logrus"
)

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

func (a *Api) Start(port int) {
	a.logger.Info("starting server on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), a.routes())
}
