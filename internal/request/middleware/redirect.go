package middleware

import (
	"fmt"
	"go2web/internal/request"
	"log/slog"
	urlpkg "net/url"
	"slices"
)

// WithRedirects wraps a GetFunc to automatically follow HTTP redirects up to maxRedirects times.
func WithRedirects(next request.GetFunc, maxRedirects int) request.GetFunc {
	return func(url string, body []byte, headers map[string]string) (*request.HttpResponse, error) {
		currentURL := url
		visited := make(map[string]struct{})

		for i := 0; i <= maxRedirects; i++ {
			if _, seen := visited[currentURL]; seen {
				return nil, fmt.Errorf("redirect loop detected to %s", currentURL)
			}
			visited[currentURL] = struct{}{}

			resp, err := next(currentURL, body, headers)
			if err != nil {
				return nil, err
			}

			if slices.Contains([]int{301, 302, 303, 307, 308}, resp.StatusCode) {
				location, ok := resp.Headers["location"]
				if !ok {
					return resp, nil
				}
				slog.Info("Redirecting initiated: ", "from", currentURL, "to", location, "status_code", resp.StatusCode)
				// Try to resolve relative Location headers against the current URL
				prevURL, err := urlpkg.Parse(currentURL)
				if err == nil {
					locURL, err := urlpkg.Parse(location)
					if err == nil {
						currentURL = prevURL.ResolveReference(locURL).String()
						continue
					}
				}

				// Fallback: use raw location value
				currentURL = location
				continue
			}

			return resp, nil
		}

		return nil, fmt.Errorf("stopped after %d redirects", maxRedirects)
	}
}
