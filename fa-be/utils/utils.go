package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 " +
		"Safari/537.36"
)

type AppEntry struct {
	// Id        string
	Img       string `json:"imgSrc"`
	Thumbnail string `json:"thumb"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Rating    string `json:"rating"`
	Link      string `json:"link"`
}

type Page struct {
	Locale string
	Data   *goquery.Document
}

type Result struct {
	Locale string     `json:"locale"`
	Apps   []AppEntry `json:"apps"`
}

var urlParts = [2]string{
	"https://play.google.com/store/search?c=apps&hl=en&gl=",
	"&q=",
}

func getPage(url string, chFailedUrls chan string, chPageData chan Page, chIsFinished chan bool, locale string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)

	defer func() {
		chIsFinished <- true
	}()

	if err != nil || resp.StatusCode != 200 {
		chFailedUrls <- url
		return
	}

	html, parseErr := goquery.NewDocumentFromReader(resp.Body)

	if parseErr != nil {
		chFailedUrls <- url
		return
	}

	page := Page{locale, html}

	chPageData <- page
}

func ScrapeData(locales []string, searchTerm string, keyword string) ([]Result, []string) {
	chFailedUrls := make(chan string)
	chPageData := make(chan Page)
	chIsFinished := make(chan bool)

	for _, locale := range locales {
		url := urlParts[0] + locale + urlParts[1] + searchTerm
		fmt.Println(url)
		go getPage(url, chFailedUrls, chPageData, chIsFinished, locale)
	}

	failedUrls := make([]string, 0)
	pageData := make([]Page, 0)

	for i := 0; i < len(locales); {
		select {
		case data := <-chPageData:
			pageData = append(pageData, data)
		case url := <-chFailedUrls:
			failedUrls = append(failedUrls, url)
		case <-chIsFinished:
			i++
		}
	}

	results := make([]Result, 0)

	for i := range pageData {
		parsed := Parse(pageData[i].Data, keyword)
		if len(parsed) > 0 {
			results = append(results, Result{Locale: pageData[i].Locale, Apps: parsed})
		}
	}

	return results, failedUrls
}

func Parse(data *goquery.Document, keyword string) []AppEntry {
	apps := make([]AppEntry, 0)
	data.Find("div[role='listitem']").Each(func(i int, s *goquery.Selection) {
		if i > 1 {
			if keyword != "" {
				if strings.Contains(strings.ToLower(s.Find("span").Text()), keyword) {
					ae := populateResults(s)
					apps = append(apps, ae)
				}
			} else {
				ae := populateResults(s)
				apps = append(apps, ae)
			}
		}
	})
	return apps
}

func populateResults(s *goquery.Selection) AppEntry {
	ae := AppEntry{}
	spans := s.Find("span")
	imgs := s.Find("img")
	link, _ := s.Find("a").Eq(0).Attr("href")
	ae.Name = spans.Eq(0).Text()
	ae.Author = spans.Eq(1).Text()
	ae.Rating = spans.Eq(2).Text()
	ae.Img, _ = imgs.Eq(0).Attr("src")
	ae.Thumbnail, _ = imgs.Eq(1).Attr("src")
	ae.Link = "https://play.google.com" + link
	return ae
}
