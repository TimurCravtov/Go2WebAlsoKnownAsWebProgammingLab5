package main

import (
	"go2web/internal/html"
	"fmt"
)

func main() {
	
	response, err := html.ParsePage("https://www.utm.md")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %s\n", response)
}
