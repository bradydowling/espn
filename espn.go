package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/manifoldco/promptui"
)

type headline struct {
	Title string
	Link  string
}

const (
	espnBaseURL    = "https://www.espn.com"
	headlineSel    = ".col-three .headlineStack ul li"
	articleContent = "#article-feed p"
)

func main() {
	headlines, err := scrapeHeadlines()
	if err != nil {
		fmt.Printf("Error scraping headlines: %s", err)
		os.Exit(1)
	}

	// Create prompt to select an article to read
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Title | cyan }}",
		Inactive: "  {{ .Title | cyan }}",
		Selected: "\U0001F449 {{ .Title | red | cyan }}",
		Details: `
--------- Article Details ----------
{{ "Title:" | faint }}  {{ .Title }}
{{ "Link:" | faint }}   {{ .Link }}`,
	}
	promptItems := append(headlines, headline{Title: "Quit"})
	prompt := promptui.Select{
		Label:     "Select an article to read",
		Items:     promptItems,
		Templates: templates,
		Size:      len(promptItems),
	}
	for {
		// Show prompt to select article
		i, _, err := prompt.Run()
		if i == len(headlines) {
			fmt.Println("Program terminated by user")
			os.Exit(0)
		}
		if err != nil {
			log.Fatalf("Error selecting article: %s", err)
		}
		selectedHeadline := headlines[i]

		articleContents, err := getArticleContents(selectedHeadline.Link)
		if err != nil {
			fmt.Printf("Error fetching article content: %s", err)
			os.Exit(1)
		}

		articleContents.Each(func(i int, s *goquery.Selection) {
			printStringByLines(s.Text(), 80)
		})

		// Prompt to select another article or quit
		promptItems = append(headlines, headline{Title: "Quit"})
		prompt = promptui.Select{
			Label:     "Select an article to read or quit",
			Items:     promptItems,
			Templates: templates,
			Size:      len(promptItems),
		}
	}
}

func printStringByLines(longString string, maxLength int) {
	longString = strings.TrimSpace(longString)

	var lines []string
	line := ""
	for i := 0; i < len(longString); i++ {
		line += string(longString[i])
		if len(line) == maxLength {
			lines = append(lines, line)
			line = ""
		}
	}
	if len(line) > 0 {
		lines = append(lines, line)
	}

	for _, l := range lines {
		fmt.Println(l)
	}
}

func scrapeHeadlines() ([]headline, error) {
	// Get the HTML content of ESPN homepage
	resp, err := http.Get(espnBaseURL)
	if err != nil {
		log.Fatalf("Error fetching ESPN homepage: %s", err)
	}
	defer resp.Body.Close()

	// Load the HTML content into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error loading ESPN homepage into goquery: %s", err)
	}

	// Select the headlines
	var headlines []headline
	doc.Find(headlineSel).Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text()
		link, _ := s.Find("a").Attr("href")
		headlines = append(headlines, headline{
			Title: title,
			Link:  espnBaseURL + link,
		})
	})

	return headlines, nil
}

func getArticleContents(url string) (*goquery.Selection, error) {
	// Get the article content
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching article content: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error loading article content into goquery: %s", err)
	}

	// Get the article text
	articleContents := doc.Find(articleContent)
	if articleContents.Text() == "" {
		return nil, fmt.Errorf("Could not find article content. Please select another article.")
	}

	// Print the article text
	return articleContents, nil
}
