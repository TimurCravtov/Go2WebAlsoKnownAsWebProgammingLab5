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

type MojeekSearchEngine struct {
	searchURL string
}

func NewMojeekSearchEngine(searchURL string) *MojeekSearchEngine {
	return &MojeekSearchEngine{searchURL: searchURL}
}

func (m *MojeekSearchEngine) Search(query string, page int, get connect.GetFunc) ([]html.SearchResult, error) {
	var headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Accept-Language": "en-US,en;q=0.9",
	}

	reqUrl := m.searchURL + url.QueryEscape(query)
	if page > 1 {
		s := (page-1)*10 + 1
		reqUrl += fmt.Sprintf("&s=%d", s)
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
	doc.Find("ul.results-standard li").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2 a.title").Text())
		href, _ := s.Find("h2 a.title").Attr("href")
		href = strings.TrimSpace(href)

		snippet := strings.TrimSpace(s.Find("p.s").Text())

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
