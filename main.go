package main

import (
	"time"
)

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

func main() {
	collector := initCollector()

	collector.Visit("https://www.nyurban.com/open-play-registration-vb/?id=35&gametypeid=1&filter_id=1")
}
