package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML(".time_schedule_table", func(e *colly.HTMLElement) {
		fmt.Print(e.Attr("class"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
	c.Visit("https://www.nyurban.com/open-play-registration-vb/?id=35&gametypeid=1&filter_id=1")
}
