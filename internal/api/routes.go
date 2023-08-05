package api

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	vb "github.com/mattfan00/mangovb"
	"github.com/mattfan00/mangovb/internal/store"
	"github.com/mattfan00/mangovb/pkg/util"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func (a *Api) routes() *chi.Mux {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
		},
	})
	r.Use(c.Handler)

	r.Use(a.logRequestMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/events", a.getEvents)
	r.Get("/filters", a.getFilters)
	r.Get("/healthcheck", a.healthcheck)

	return r
}

func (a *Api) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t1 := time.Now()
		defer func() {
			a.logger.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     r.URL,
				"from":     r.RemoteAddr,
				"duration": time.Since(t1),
			}).Info("received request")
		}()

		next.ServeHTTP(w, r)
	})
}

func parseQuery(query string) []int {
	if query == "" {
		return []int{}
	}

	arr, _ := util.SliceAtoi(strings.Split(query, "|"))
	return arr
}

func (a *Api) getEvents(w http.ResponseWriter, r *http.Request) {
	filters := store.EventQueryFilters{
		Source:     parseQuery(r.URL.Query().Get("source")),
		SkillLevel: parseQuery(r.URL.Query().Get("skillLevel")),
		Spots:      parseQuery(r.URL.Query().Get("spots")),
	}

	events, err := a.eventStore.GetLatest(true, filters)
	if err != nil {
		renderError(w, ErrInternalServer(err))
		return
	}

	render(w, http.StatusOK, events)
}

type GetFiltersResponse struct {
	Source     []FilterEntry `json:"source"`
	SkillLevel []FilterEntry `json:"skillLevel"`
	Spots      []FilterEntry `json:"spots"`
}

type FilterEntry struct {
	Value int         `json:"value"`
	Text  interface{} `json:"text"`
}

type ByValue []FilterEntry

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (a *Api) getFilters(w http.ResponseWriter, r *http.Request) {
	res := GetFiltersResponse{}

	source := []FilterEntry{}
	for k, v := range vb.EventSourceMap {
		filter := FilterEntry{
			Value: int(k),
			Text:  v,
		}
		source = append(source, filter)
	}
	sort.Sort(ByValue(source))
	res.Source = source

	skillLevel := []FilterEntry{}
	for k, v := range vb.EventSkillLevelMap {
		filter := FilterEntry{
			Value: int(k),
			Text:  v,
		}
		skillLevel = append(skillLevel, filter)
	}
	sort.Sort(ByValue(skillLevel))
	res.SkillLevel = skillLevel

	res.Spots = []FilterEntry{
		{Value: 0, Text: "Filled"},
		{Value: 1, Text: "Available"},
	}

	render(w, http.StatusOK, res)
}

func (a *Api) healthcheck(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "ok")
}
