package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Character struct {
	name      string
	url       string
	full_name string
	image_url string
}

func main() {

	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	c := colly.NewCollector(colly.AllowedDomains("myanimelist.net"))

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Name", "Full Name", "Image URL"}
	writer.Write(headers)

	c.OnHTML("table.js-anime-character-table", func(e *colly.HTMLElement) {
		e.ForEach("td:nth-child(2)>.spaceit_pad>a", func(_ int, el *colly.HTMLElement) {
			url := strings.TrimSpace(el.Attr("href"))
			name := strings.TrimSpace(el.ChildText(".spaceit_pad > a > h3"))
			getInfo(url, name, writer)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		println("Request URL:", r.Request.URL, "failed with", err)
	})

	c.OnScraped(func(r *colly.Response) {
		println("Finished", r.Request.URL.String())
	})

	/*
		Start scraping for this MAL anime character list
		'insert url of the characters page here'
	*/
	c.Visit("https://myanimelist.net/anime/21/One_Piece/characters")
}

func getInfo(url string, name string, writer *csv.Writer) {

	character := Character{}
	character.url = strings.TrimSpace(url)
	character.name = strings.TrimSpace(name)
	c := colly.NewCollector(colly.AllowedDomains("myanimelist.net"))
	c.OnHTML("div#contentWrapper", func(e *colly.HTMLElement) {
		character.full_name = strings.TrimSpace(e.ChildText("div.edit-info>.h1-title>.title-name"))
	})

	c.OnHTML(".portrait-225x350", func(e *colly.HTMLElement) {
		character.image_url = e.Attr("data-src")
	})

	c.OnScraped(func(r *colly.Response) {
		println(character.full_name, character.image_url)
		row := []string{character.name, character.full_name, character.image_url}
		writer.Write(row)
	})

	c.Visit(url)
}
