package search_engines

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"go2web/internal/connect"
	"go2web/internal/html"

	"github.com/PuerkitoBio/goquery"
)

type StartpageSearchEngine struct {
	searchURL string
}

func NewStartpageSearchEngine(searchURL string) *StartpageSearchEngine {
	return &StartpageSearchEngine{searchURL: searchURL}
}

func (s *StartpageSearchEngine) Search(query string, page int, get connect.GetFunc) ([]html.SearchResult, error) {
	var headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
	}

	reqUrl := s.searchURL + url.QueryEscape(query)
	if page > 1 {
		reqUrl += fmt.Sprintf("&page=%d", page)
	}

	res, err := get(reqUrl, nil, headers)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	var results []html.SearchResult
	doc.Find(".w-gl .result").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find("a.result-title h2").Text())
		href, _ := sel.Find("a.result-title").Attr("href")
		href = strings.TrimSpace(href)

		snippet := strings.TrimSpace(sel.Find("p.description").Text())

		if title != "" && href != "" && strings.HasPrefix(href, "http") {
			results = append(results, html.SearchResult{
				Title:   title,
				URL:     href,
				Snippet: snippet,
			})
		}
	})

	return results, nil
}
