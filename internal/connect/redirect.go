package connect

import "slices"

func WithRedirects(next GetFunc) GetFunc {
	return func(url string, body []byte, headers map[string]string) (*HttpResponse, error) {

		resp, err := next(url, body, headers)

		if err != nil {
			return nil, err
		}

		if slices.Contains([]int{301, 302, 303, 307, 308}, resp.StatusCode) {
			location, ok := resp.Headers["location"]
			if ok {
				return next(location, body, headers)
			}
		}
		return resp, nil
	}
}
