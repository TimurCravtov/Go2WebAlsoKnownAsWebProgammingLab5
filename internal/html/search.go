package html

type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

type Search interface {
	Search(query string) ([]SearchResult, error)
}
