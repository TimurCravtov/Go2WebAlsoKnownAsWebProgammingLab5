package cli

import (
	"fmt"
	"go2web/internal/request"
	"go2web/internal/html"
	"go2web/internal/html/negociation"
    "go2web/internal/cli/printer"
	"math"
    "log/slog"
	"strings"
    "go2web/internal/request/middleware"
	_ "github.com/mat/besticon/ico"
	"github.com/spf13/cobra"
)

func HandleUrlRequest(cmd *cobra.Command, args []string) {

    rawURL, _ := cmd.Flags().GetString("url")
    urlStr := rawURL
    if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
        urlStr = "https://" + urlStr
    }

    var getter request.GetFunc = request.Get
    noCache, _ := cmd.Flags().GetBool("no-cache")
    if !noCache {
        cache := middleware.NewFileCache("cache")
        getter = cache.WithCache(getter)
    }

    redirectCount, _ := cmd.Flags().GetInt("max-redirects")
    if redirectCount < 0 {
        redirectCount = math.MaxInt
    }
    if redirectCount >= 0 {
        getter = middleware.WithRedirects(getter, redirectCount)
    }


    languages, _ := cmd.Flags().GetStringArray("lang")
    charsets, _ := cmd.Flags().GetStringArray("charset")
    types, _ := cmd.Flags().GetStringArray("type")

    if len(languages) > 0 {
        slog.Debug("Adding Accept-Language header", "values ", languages)
        getter = html.WithHeaders(
            negociation.AcceptLanguages(languages),
        )(getter)
    }

    if len(charsets) > 0 {
        slog.Debug("Adding Accept-Charset header with values ", "values", charsets)
        getter = html.WithHeaders(
            negociation.AcceptCharsets(charsets),
        )(getter)
    }

    if len(types) > 0 {
        slog.Debug("Adding Accept-Content-Type header with values ", "values", types)
        getter = html.WithHeaders(
            negociation.AcceptContentTypes(types),
        )(getter)
    }

    
    response, err := getter(urlStr, nil, nil)
    if err != nil {
        slog.Error("Error fetching page", "error", err)
        return
    }

    var basePrinter printer.HttpResponsePrinter

    contentType, err := html.GetContentType(response)

    if err != nil {
        slog.Error("Error determining content type", "error", err)
        return
    }

    switch contentType {
    case html.TypeHTML:
        basePrinter = printer.HtmlResponseParser
    case html.TypeJSON:

        basePrinter = printer.JsonPrinter
    case html.TypePNG, html.TypeJPEG, html.TypeGIF:
        basePrinter = printer.ImagePrinter
    default:
        basePrinter = printer.TextPrinter
    }

    printer := printer.WithStatusLine(printer.WithHeaders(printer.WithHero(basePrinter)))
    
    str, _ := printer(urlStr, response);
    
    fmt.Println(str)

}


