package main

import (
	"fmt"
	"go2web/internal/connect"
	"go2web/internal/html"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	// engine := search_engines.NewStartpageSearchEngine("https://www.startpage.com/sp/search?query=")
	// response, err := engine.Search("cats", 1)

	// fmt.Println(calvin.AsciiFont("STARTPAGE"))

	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// for _, result := range response {
	// 	fmt.Printf("Title: %s\n", result.Title)
	// 	fmt.Printf("URL: %s\n", result.URL)
	// 	fmt.Println("--------------------------------------------------")
	// }

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	}))

	slog.SetDefault(logger)

	cache := connect.NewFileCache("cache")

	cachedGet := cache.WithCache(connect.Get)
	redirectGet := connect.WithRedirects(cachedGet)

	fmt.Println(html.ParsePage("https://point.md", redirectGet))

}
