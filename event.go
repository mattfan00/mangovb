package vb

import "time"

type Event struct {
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
}
