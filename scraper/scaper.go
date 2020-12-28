package scraper

import (
	"fmt"
	"github.com/filipwtf/filips-installer/ui"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

const (
	sk1erClub   = "https://sk1er.club/"
	agent       = "Filip's Mod Manager"
	moduleTitle = "#masonry > div:nth-child(?) > section > div > a"
	downloadCSS = "#Form-? > div > p:nth-child(10) > a"
)

type Mod struct {
	Name string
	Link string
}

func ScrapeSk1erMods() *[]Mod {
	c := colly.NewCollector(
		colly.UserAgent(agent),
	)

	var links []string
	for i := 0; i <= 39; i++ {
		current := strings.Replace(moduleTitle, "?", strconv.FormatInt(int64(i), 10), 1)
		c.OnHTML(current, func(e *colly.HTMLElement) {
			links = append(links, e.Attr("href"))
		})
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(sk1erClub)
	if err != nil {
		ui.Log(err.Error())
	}
	c.Wait()

	var mods []Mod
	c = colly.NewCollector(
		colly.UserAgent(agent),
		colly.Async(true),
	)
	for _, url := range links {
		modName := strings.SplitAfter(url, "/")[4]
		current := strings.Replace(downloadCSS, "?", modName, 1)
		c.OnHTML(current, func(e *colly.HTMLElement) {
			fmt.Println(e.DOM.Text())
			fmt.Println(e.Attr("href"))
			mods = append(mods, Mod{
				Name: e.DOM.Text(),
				Link: e.Attr("href"),
			})
		})
		err := c.Visit(url)
		if err != nil {
			ui.Log(err.Error())
		}
	}

	return &mods
}
