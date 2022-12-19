package vb

import "time"

type Event struct {
	Id          string
	Name        string
	Location    string
	StartDate   time.Time
	StartTime   string
	EndTime     string
	SkillLevel  string
	Price       float64
	IsAvailable bool
	SpotsLeft   int
	Url         string
	Source      string
	ScrapedOn   time.Time
}
