package search_engines

import (
	"bytes"
	"net/url"
	"strings"

	"go2web/internal/connect"
	"go2web/internal/html"

	"github.com/PuerkitoBio/goquery"
)

type BingSearchEngine struct {
	searchURL string
}

func NewBingSearchEngine(searchURL string) *BingSearchEngine {
	return &BingSearchEngine{searchURL: searchURL}
}

func (b *BingSearchEngine) Search(query string, get connect.GetFunc) ([]html.SearchResult, error) {
	var headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Accept-Language": "en-US,en;q=0.9",
	}

	res, err := get(b.searchURL+url.QueryEscape(query), nil, headers)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	var results []html.SearchResult
	doc.Find("li.b_algo").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2 a").Text())
		href, _ := s.Find("h2 a").Attr("href")
		href = strings.TrimSpace(href)

		snippet := strings.TrimSpace(s.Find(".b_caption p").Text())
		if snippet == "" {
			snippet = strings.TrimSpace(s.Find(".b_algoSlug").Text())
		}

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
