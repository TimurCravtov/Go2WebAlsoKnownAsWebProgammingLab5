package cli

import (
    "fmt"
    "go2web/internal/connect"
    "go2web/internal/html"
    "math"
    "net/url"
    "regexp"
    "strings"

    "github.com/spf13/cobra"

    "image"
    _ "image/jpeg"
    _ "image/png"
    "net/http"

    "github.com/0magnet/calvin"
    "github.com/gookit/color"
    _ "github.com/mat/besticon/ico"
    "golang.org/x/image/draw"
)

func HandleUrlRequest(cmd *cobra.Command, args []string) {
    rawURL, _ := cmd.Flags().GetString("url")
    urlStr := rawURL
    if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
        urlStr = "https://" + urlStr
    }

    var getter connect.GetFunc = connect.Get
    noCache, _ := cmd.Flags().GetBool("no-cache")
    if !noCache {
        cache := connect.NewFileCache("cache")
        getter = cache.WithCache(getter)
    }

    redirectCount, _ := cmd.Flags().GetInt("max-redirects")
    if redirectCount < 0 {
        redirectCount = math.MaxInt
    }
    if redirectCount >= 0 {
        getter = connect.WithRedirects(getter, redirectCount)
    }

    response, httpRes, err := html.ParsePage(urlStr, getter)
    if err != nil {
        fmt.Printf("Error fetching page: %v\n", err)
        return
    }

    heroText := buildWebsiteHero(httpRes, urlStr)
    if heroText != "" {
        fmt.Println(heroText)
    }

    fmt.Println(response)
}

func buildWebsiteHero(response *connect.HttpResponse, rootUrl string) string {
    faviconUrl := getFavicoLink(response, rootUrl)
    asciiFavicon, err := generateColoredFaviconASCII(faviconUrl)
    var sb strings.Builder

    u, _ := url.Parse(rootUrl)
    // Hostname() strictly returns the domain (e.g., example.com) without ports or paths
    websiteName := u.Hostname() 
    
    // Generate the ASCII title and split it into lines
    asciiTitle := calvin.AsciiFont(strings.ToUpper(websiteName))
    titleLines := strings.Split(strings.TrimRight(asciiTitle, "\n"), "\n")

    var iconLines []string
    boxWidth := 24 
    
    if err == nil {
        rawIconLines := strings.Split(strings.TrimRight(asciiFavicon, "\n"), "\n")
        
        iconLines = append(iconLines, "╭"+strings.Repeat("─", boxWidth-2)+"╮")
        for _, line := range rawIconLines {
            iconLines = append(iconLines, "│ "+line+" │")
        }
        iconLines = append(iconLines, "╰"+strings.Repeat("─", boxWidth-2)+"╯")
    }

    iconHeight := len(iconLines)
    titleHeight := len(titleLines)
    
    maxLines := iconHeight
    if titleHeight > maxLines {
        maxLines = titleHeight
    }

    // Calculate vertical starting offsets for true centering
    iconOffset := (maxLines - iconHeight) / 2
    titleOffset := (maxLines - titleHeight) / 2

    emptyIconPadding := strings.Repeat(" ", boxWidth)

    for i := 0; i < maxLines; i++ {
        // Determine the icon row (or pad with spaces if above/below the icon bounds)
        iconPart := emptyIconPadding
        if len(iconLines) > 0 && i >= iconOffset && i < iconOffset+iconHeight {
            iconPart = iconLines[i-iconOffset]
        } else if len(iconLines) == 0 {
            iconPart = ""
        }

        titlePart := ""
        if i >= titleOffset && i < titleOffset+titleHeight {
            titlePart = titleLines[i-titleOffset]
        }

        if iconPart != "" && titlePart != "" {
            sb.WriteString(iconPart + "   " + titlePart + "\n")
        } else if iconPart != "" {
            sb.WriteString(iconPart + "\n")
        } else {
            sb.WriteString(titlePart + "\n")
        }
    }

    return sb.String()
}

func getFavicoLink(response *connect.HttpResponse, rootUrl string) string {
    baseURL, err := url.Parse(rootUrl)
    if err != nil {
        return ""
    }

    bodyStr := string(response.Body)

    linkRegex := regexp.MustCompile(`(?i)<link[^>]+>`)
    hrefRegex := regexp.MustCompile(`(?i)href\s*=\s*["']([^"']+)["']`)
    relRegex := regexp.MustCompile(`(?i)rel\s*=\s*["']([^"']+)["']`)

    links := linkRegex.FindAllString(bodyStr, -1)

    for _, linkTag := range links {
        relMatch := relRegex.FindStringSubmatch(linkTag)
        if len(relMatch) > 1 {
            relVal := strings.ToLower(relMatch[1])

            if strings.Contains(relVal, "icon") {
                hrefMatch := hrefRegex.FindStringSubmatch(linkTag)
                if len(hrefMatch) > 1 {
                    rawHref := hrefMatch[1]

                    hrefURL, err := url.Parse(rawHref)
                    if err == nil {
                        // Resolve relative URLs against the base URL
                        resolvedURL := baseURL.ResolveReference(hrefURL)
                        return resolvedURL.String()
                    }
                }
            }
        }
    }

    fallbackURL, _ := url.Parse("/favicon.ico")
    resolvedFallback := baseURL.ResolveReference(fallbackURL)

    return resolvedFallback.String()
}

func generateColoredFaviconASCII(iconURL string) (string, error) {
    if iconURL == "" {
        return "", fmt.Errorf("no favicon URL provided")
    }

    // 1. Fetch the image
    resp, err := http.Get(iconURL)
    if err != nil {
        return "", fmt.Errorf("failed to fetch favicon: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("bad status fetching favicon: %s", resp.Status)
    }

    // 2. Decode the image (Accepts PNG/JPG)
    img, _, err := image.Decode(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to decode image: %v", err)
    }

    // 3. Resize the image to 5x5
    destRect := image.Rect(0, 0, 10, 10)
    dest := image.NewRGBA(destRect)
    draw.ApproxBiLinear.Scale(dest, dest.Bounds(), img, img.Bounds(), draw.Src, nil)

    // 4. Map pixels to COLOURED ASCII characters
    // Characters ordered by density. Used as a fallback or combined shape.
    asciiChars := []rune{'@', '%', '#', '*', '+', '=', '-', ':', '.', ' '}
    var sb strings.Builder

    for y := 0; y < 10; y++ {
        for x := 0; x < 10; x++ {
            c := dest.At(x, y)
            // Go's image package returns alpha-premultiplied values (0-65535)
            r16, g16, b16, _ := c.RGBA()

            // Convert to standard 8-bit RGB (0-255) required for terminal colors
            r := uint8(r16 >> 8)
            g := uint8(g16 >> 8)
            b := uint8(b16 >> 8)

            // Calculate luminance to select the appropriate ASCII density char
            lum := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

            // Map luminance (0-255) to the index of the asciiChars array
            idx := int((lum * float64(len(asciiChars)-1)) / 255.0)

            // Clamp index
            if idx < 0 {
                idx = 0
            }
            if idx >= len(asciiChars) {
                idx = len(asciiChars) - 1
            }
            char := asciiChars[idx]

            coloredChar := color.RGB(r, g, b).Sprintf("%c ", char)
            sb.WriteString(coloredChar)
        }
        sb.WriteString("\n")
    }

    return sb.String(), nil
}