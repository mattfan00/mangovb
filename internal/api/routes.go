package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
    "github.com/rs/cors"
)

func (a *Api) routes() *chi.Mux {
	r := chi.NewRouter()

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
    })
    r.Use(c.Handler)

	r.Get("/events", a.getEvents)

	return r
}

func (a *Api) getEvents(w http.ResponseWriter, r *http.Request) {
	events, err := a.eventStore.GetAll()
	if err != nil {
		renderError(w, ErrInternalServer(err))
		return
	}

	render(w, http.StatusOK, events)
}
