package html

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"golang.org/x/net/html" 
	"regexp"
)

func ParsePage(url string) (string, error) {
    res, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    doc, err := goquery.NewDocumentFromReader(res.Body)
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