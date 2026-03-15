package connect

// kindof decorator
type GetFunc func(url string, body []byte, headers map[string]string) (*HttpResponse, error)
