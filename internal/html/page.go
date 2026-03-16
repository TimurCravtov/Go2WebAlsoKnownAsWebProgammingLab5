package html

import (
	"bytes"
	"go2web/internal/connect"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

func ParsePage(pageURL string, get connect.GetFunc) (string, *connect.HttpResponse, error) {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	res, err := get(pageURL, nil, headers)
	if err != nil {
		return "", nil, err
	}

	reader := bytes.NewReader(res.Body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", nil, err
	}

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return "", nil, err
	}

	doc.Find("script, style, iframe, noscript, nav, footer, .sidebar, .menu").Remove()

	// Convert relative hrefs to absolute
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			linkURL, err := url.Parse(href)
			if err == nil {
				absoluteURL := baseURL.ResolveReference(linkURL)
				s.SetAttr("href", absoluteURL.String())
			}
		}
	})

	selection := doc.Find("main")
	if selection.Length() == 0 {
		selection = doc.Find("body")
	}

	htmlContent, err := selection.Html()
	if err != nil {
		return "", nil, err
	}

	text, err := html2text.FromString(htmlContent, html2text.Options{
		PrettyTables: true,
		OmitLinks:    false,
	})
	if err != nil {
		return "", nil, err
	}

	coloredText := colorizeURLs(text)

	return strings.TrimSpace(coloredText), res, nil
}
