package main

import (
	"go2web/internal/html"
	"fmt"
	"github.com/0magnet/calvin"
)

func main() {
	
	engine := html.NewBingSearchEngine("https://www.bing.com/search?q=")
	response, err := engine.Search("cats")

	fmt.Println(calvin.AsciiFont("BING")) 

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, result := range response {
		fmt.Printf("Title: %s\n", result.Title)
		fmt.Printf("URL: %s\n", result.URL)
		fmt.Println("--------------------------------------------------")
	}

}
