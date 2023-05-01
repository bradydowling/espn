package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Make HTTP GET request to ESPN.com
	res, err := http.Get("https://www.espn.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Parse the HTML using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find all the frontpage headlines st
	headlines := doc.Find(".col-three .headlineStack ul li")

	// Print out the headlines
	fmt.Println("ESPN Frontpage Headlines:")
	headlines.Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%d. %s\n", i+1, s.Text())
	})
}
