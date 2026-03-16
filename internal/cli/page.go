package cli

import (
	"fmt"
	"go2web/internal/connect"
	"go2web/internal/html"
	"strings"

	"github.com/spf13/cobra"
)

func HandleUrlRequest(cmd *cobra.Command, args []string) {

	rawURL, _ := cmd.Flags().GetString("url")
	url := rawURL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	var getter connect.GetFunc = connect.Get
	noCache, _ := cmd.Flags().GetBool("no-cache")
	if !noCache {
		cache := connect.NewFileCache("cache")
		getter = cache.WithCache(getter)
	}

	getter = connect.WithRedirects(getter)

	response, err := html.ParsePage(url, getter)
	if err != nil {
		fmt.Printf("Error fetching page: %v\n", err)
		return
	}

	fmt.Println(response)

}
