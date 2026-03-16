package printer

import "go2web/internal/connect"

type HttpResponsePrinter func(url string, response *connect.HttpResponse) (string, error)
