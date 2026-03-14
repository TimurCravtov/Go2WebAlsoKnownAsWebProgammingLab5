package main

import (
	"fmt"
	"go2web/internal/connect"
)

func main() {
	// cmd.Execute()
	

	response, err := connect.Get("https://point.md", nil, make(map[string]string))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response:\n%s\n", string(response))
	
}
