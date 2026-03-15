package search_engines

import (
	"bytes"
	"fmt"
	"go2web/internal/connect"
	"go2web/internal/html"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type DuckSearchEngine struct {
	searchURL string
}

func NewDuckSearchEngine(searchURL string) *DuckSearchEngine {
	return &DuckSearchEngine{searchURL: searchURL}
}

func (d *DuckSearchEngine) Search(query string, get connect.GetFunc) ([]html.SearchResult, error) {
	// 'kl=uk-en' forces United Kingdom results
	reqUrl := fmt.Sprintf("https://duckduckgo.com/lite/?q=%s&kl=uk-en", url.QueryEscape(query))

	// DDG Lite doesn't even need complex headers
	res, err := get(reqUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	var results []html.SearchResult
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body))
	// DDG Lite uses a simple table structure:
	// Results are usually in 'a.result-link'
	doc.Find("a.result-link").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		href, _ := s.Attr("href")
		results = append(results, html.SearchResult{Title: title, URL: href})
	})
	return results, nil
}
