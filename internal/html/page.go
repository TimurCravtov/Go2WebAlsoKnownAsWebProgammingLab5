package html

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"go2web/internal/connect"
	"golang.org/x/net/html" 
	"regexp"
)

func ParsePage(url string) (string, error) {

    var headers = map[string]string{
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    }

    res, err := connect.Get(url, nil, headers)
    if err != nil {
        return "", err
    }

	reader := bytes.NewReader(res.Body)

    doc, err := goquery.NewDocumentFromReader(reader)
    if err != nil {
        return "", err
    }

    doc.Find("script, style, iframe, noscript, nav, footer, .sidebar, .menu").Remove()

    var builder strings.Builder
    
    selection := doc.Find("main")
    if selection.Length() == 0 {
        selection = doc.Find("body")
    }

    selection.Each(func(i int, s *goquery.Selection) {
        extractText(s.Get(0), &builder)
    })

    rawText := builder.String()

    reMultiLine := regexp.MustCompile(`\n{3,}`)
    cleanText := reMultiLine.ReplaceAllString(rawText, "\n\n")

    return strings.TrimSpace(cleanText), nil
}

func extractText(n *html.Node, builder *strings.Builder) {
    switch n.Type {

	case html.TextNode:
        text := strings.TrimSpace(n.Data)
        if text != "" {
            builder.WriteString(text + " ")
        }
    case html.ElementNode:

		blockElements := map[string]bool{
            "p": true, "div": true, "h1": true, "h2": true, 
            "h3": true, "li": true, "article": true, "br": true,
        }
        
        if blockElements[n.Data] {
            builder.WriteString("\n")
        }

        for c := n.FirstChild; c != nil; c = c.NextSibling {
            extractText(c, builder)
        }

        if blockElements[n.Data] {
            builder.WriteString("\n")
        }
    }
}