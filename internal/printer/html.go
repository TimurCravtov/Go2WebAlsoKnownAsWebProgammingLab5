package printer

import (
	"bytes"
	"net/url"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
	"go2web/internal/connect"
	"go2web/internal/html"
)

func HtmlResponseParser(urlPath string, response *connect.HttpResponse) (string, error) {

    reader := bytes.NewReader(response.Body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	baseURL, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}

	doc.Find("script, style, iframe, noscript, nav, footer, .sidebar, .menu").Remove()

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
		return "", err
	}

	text, err := html2text.FromString(htmlContent, html2text.Options{
		PrettyTables: true,
		OmitLinks:    false,
	})
	if err != nil {
		return "", err
	}

	coloredText := html.ColorizeURLs(text)

	return strings.TrimSpace(coloredText), nil
}

