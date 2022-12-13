package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML("div.time_schedule_table", func(tableDiv *colly.HTMLElement) {
		location := strings.TrimSuffix(tableDiv.ChildText("h3 > span"), ":")
		fmt.Println(location)

		tableDiv.ForEach("table tr", func(rowNum int, row *colly.HTMLElement) {
			if rowNum == 0 {
				return
			}

			fmt.Printf("%s\t", row.ChildText("td:nth-child(2)"))
			fmt.Printf("%s\t", row.ChildText("td:nth-child(3)"))
			fmt.Printf("%s\t", row.ChildText("td:nth-child(4)"))
			fmt.Printf("%s\t", row.ChildText("td:nth-child(5)"))
			fmt.Printf("%s\n", row.ChildText("td:nth-child(6)"))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
	c.Visit("https://www.nyurban.com/open-play-registration-vb/?id=35&gametypeid=1&filter_id=1")
}
