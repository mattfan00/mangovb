package api

import (
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"
	vb "github.com/mattfan00/mangovb"
	"github.com/rs/cors"
)

func (a *Api) routes() *chi.Mux {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
	})
	r.Use(c.Handler)

	r.Get("/events", a.getEvents)
	r.Get("/filters", a.getFilters)

	return r
}

func (a *Api) getEvents(w http.ResponseWriter, r *http.Request) {
	events, err := a.eventStore.GetLatest(true)
	if err != nil {
		renderError(w, ErrInternalServer(err))
		return
	}

	render(w, http.StatusOK, events)
}

type FilterResponse struct {
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
	res := FilterResponse{}

	sourceFilters := []FilterEntry{}
	for k, v := range vb.EventSourceMap {
		filter := FilterEntry{
			Value: int(k),
			Text:  v,
		}
		sourceFilters = append(sourceFilters, filter)
	}
	sort.Sort(ByValue(sourceFilters))
	res.Source = sourceFilters

	skillLevelFilters := []FilterEntry{}
	for k, v := range vb.EventSkillLevelMap {
		filter := FilterEntry{
			Value: int(k),
			Text:  v,
		}
		skillLevelFilters = append(skillLevelFilters, filter)
	}
	sort.Sort(ByValue(skillLevelFilters))
	res.SkillLevel = skillLevelFilters

	res.Spots = []FilterEntry{
		{Value: 0, Text: "Filled"},
		{Value: 1, Text: "Available"},
	}

	render(w, http.StatusOK, res)
}
