package main

import (
	"fmt"
	"log/slog"
	"go2web/internal/html/search_engines"
	"os"
	"time"
	"github.com/0magnet/calvin"
	"github.com/lmittmann/tint"
	"go2web/internal/html"
	"go2web/internal/connect"
)

func main() {

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	}))

	slog.SetDefault(logger)


	engine := search_engines.NewStartpageSearchEngine("https://www.startpage.com/sp/search?query=")

	query := "cats"
	response, err := engine.Search(query, 1, connect.Get)

	fmt.Println(calvin.AsciiFont("STARTPAGE"))

	fmt.Println("╭-----------------------------------------------╮")
	fmt.Printf("| %-42s ⌕  |\n", query);
	fmt.Println("╰-----------------------------------------------╯")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, result := range response {
		fmt.Printf("Title: %s\n", result.Title)
		fmt.Printf("URL: %s\n", html.Colorize(result.URL, html.ColorBlue))
		fmt.Println("--------------------------------------------------")
	}


	// cache := connect.NewFileCache("cache")

	// cachedGet := cache.WithCache(connect.Get)
	// redirectGet := connect.WithRedirects(cachedGet)

	// fmt.Println(html.ParsePage("https://point.md", redirectGet))

}
