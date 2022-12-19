package vb

import "time"

type Event struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Location    string    `db:"location"`
	StartDate   time.Time `db:"start_date"`
	StartTime   string    `db:"start_time"`
	EndTime     string    `db:"end_time"`
	SkillLevel  string
	Price       float64   `db:"price"`
	IsAvailable bool      `db:"is_available"`
	SpotsLeft   int       `db:"spots_left"`
	Url         string    `db:"url"`
	Source      string    `db:"source"`
	ScrapedOn   time.Time `db:"scraped_on"`
}
