package utils

import (
	"encoding/json"
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
}

var urlParts = [1]string{
	"https://play.google.com/store/search?q=sportybet&c=apps&hl=en&gl=",
}

func getPage(url string, chFailedUrls chan string, chPageData chan *goquery.Document, chIsFinished chan bool) {
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

	fmt.Println(html)

	chPageData <- html
}

func ScrapeData(locales []string) ([]byte, []string) {
	chFailedUrls := make(chan string)
	chPageData := make(chan *goquery.Document)
	chIsFinished := make(chan bool)

	for _, locale := range locales {
		url := urlParts[0] + locale
		go getPage(url, chFailedUrls, chPageData, chIsFinished)
	}

	failedUrls := make([]string, 0)
	pageData := make([]*goquery.Document, 0)

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

	parsed := Parse(pageData)

	return parsed, failedUrls
}

func Parse(data []*goquery.Document) []byte {
	results := make([]AppEntry, 0)
	for i := range data {
		data[i].Find("div[role='listitem']").Each(func(i int, s *goquery.Selection) {
			if i > 1 {
				ae := AppEntry{}
				s.Find("span").Each(func(j int, sp *goquery.Selection) {
					if j == 0 && strings.Contains(strings.ToLower(sp.Text()), "sportybet") {
						ae.Name = sp.Text()
						s.Find("img").Each(func(k int, img *goquery.Selection) {
							src, ok := img.Attr("src")
							if ok && k == 0 {
								ae.Img = src
							}
							if ok && k == 1 {
								ae.Thumbnail = src
							}
						})
					}

					if ae.Name != "" {
						if j == 1 {
							ae.Author = sp.Text()
						}
						if j == 2 {
							ae.Rating = sp.Text()
						}
					}
				})
				if ae.Name != "" {
					results = append(results, ae)
				}
			}
		})
	}

	byteArr, _ := json.Marshal(results)
	return byteArr
}
