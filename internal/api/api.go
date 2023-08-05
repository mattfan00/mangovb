package api

import (
	"fmt"
	"net/http"

	"github.com/mattfan00/mangovb/internal/store"
	configPkg "github.com/mattfan00/mangovb/pkg/config"

	"github.com/sirupsen/logrus"
)

type Api struct {
	logger     *logrus.Entry
	eventStore *store.EventStore
	config     *configPkg.Config
}

func New(logger *logrus.Entry, eventStore *store.EventStore, config *configPkg.Config) *Api {
	return &Api{
		logger:     logger,
		eventStore: eventStore,
		config:     config,
	}
}

func (a *Api) Start(port int) {
	a.logger.Info("starting server on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), a.routes())
}
