package cli

import (
	"fmt"
	"go2web/internal/connect"
	"go2web/internal/html"
	"go2web/internal/html/search_engines"
	"strings"

	"github.com/0magnet/calvin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// ── Messages ──────────────────────────────────────────────────────────────────

type searchResultsMsg struct {
	results []html.SearchResult
	err     error
}

type selectedMsg struct {
	result html.SearchResult
}

// ── Model ─────────────────────────────────────────────────────────────────────

type searchModel struct {
	// config
	engineName string
	query      string
	engine     html.Search

	// state
	results  []html.SearchResult
	cursor   int
	loading  bool
	err      error
	selected *html.SearchResult
	height   int

	// hero banner (rendered once)
	hero string
}

func HandleSearch(cmd *cobra.Command, args []string) {
	searchQuery, _ := cmd.Flags().GetString("search")
	if searchQuery == "" {
		return
	}
	// select engine
	engineName, _ := cmd.Flags().GetString("engine")
	var engine html.Search
	switch engineName {
	case "startpage":
		engine = search_engines.NewStartpageSearchEngine("https://www.startpage.com/sp/search?query=")
	case "mojeek":
		engine = search_engines.NewMojeekSearchEngine("https://www.mojeek.com/search?q=")
	default:
		fmt.Printf("Unknown search engine: %s\n", engineName)
		return
	}
	fmt.Println(buildHero(engineName, searchQuery))
	// execute
	results, err := engine.Search(searchQuery, 1, connect.Get)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// print results
	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", html.Colorize(result.URL, html.ColorBlue))
		fmt.Println("   " + "─" + "─" + "─" + "─" + "─" + "─" + "─" + "─" + "─" + "─" + "─" + "─" + "─")
	}
}

func newSearchModel(engineName, query string, engine html.Search) searchModel {
	return searchModel{
		engineName: engineName,
		query:      query,
		engine:     engine,
		loading:    true,
		hero:       buildHero(engineName, query),
	}
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m searchModel) Init() tea.Cmd {
	return fetchResults(m.engine, m.query)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height

	case searchResultsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.results = msg.results
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.results) > 0 {
				selected := m.results[m.cursor]
				m.selected = &selected
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m searchModel) View() string {
	var sb strings.Builder

	sb.WriteString(m.hero)
	sb.WriteString("\n")

	switch {
	case m.loading:
		sb.WriteString("  Searching...\n")

	case m.err != nil:
		sb.WriteString(fmt.Sprintf("  Error: %v\n", m.err))

	case m.selected != nil:
		// ── Placeholder "page" after selection ──────────────────────────────
		sb.WriteString(renderSelectedPlaceholder(m.selected))

	default:
		visibleResults := 5
		heroLines := strings.Count(m.hero, "\n") + 1
		if m.height > 0 {
			available := m.height - heroLines - 4
			if available > 3 {
				visibleResults = available / 3
			} else {
				visibleResults = 1
			}
		}

		start := m.cursor - visibleResults/2
		if start < 0 {
			start = 0
		}
		end := start + visibleResults
		if end > len(m.results) {
			end = len(m.results)
			start = end - visibleResults
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			r := m.results[i]
			cursor := "  "
			titleColor := html.ColorReset
			if i == m.cursor {
				cursor = html.Colorize("▶ ", html.ColorCyan)
				titleColor = html.ColorCyan
			}

			sb.WriteString(fmt.Sprintf(
				"%s%d. %s\n",
				cursor, i+1,
				html.Colorize(r.Title, titleColor),
			))
			sb.WriteString(fmt.Sprintf(
				"     %s\n",
				html.Colorize(r.URL, html.ColorBlue),
			))
			sb.WriteString("   " + strings.Repeat("─", 44) + "\n")
		}

		sb.WriteString("\n  ↑/↓  navigate   enter  open   q  quit\n")
	}

	return sb.String()
}

func fetchResults(engine html.Search, query string) tea.Cmd {
	return func() tea.Msg {
		results, err := engine.Search(query, 1, connect.Get)
		return searchResultsMsg{results: results, err: err}
	}
}

func buildHero(engineName, query string) string {
	title := calvin.AsciiFont(strings.ToUpper(engineName))
	box := fmt.Sprintf(
		"╭───────────────────────────────────────────────╮\n│ %-43s ⌕ │\n╰───────────────────────────────────────────────╯",
		query,
	)
	return title + "\n" + box
}

func renderSelectedPlaceholder(r *html.SearchResult) string {
	separator := strings.Repeat("═", 50)
	return fmt.Sprintf(
		"\n  %s\n\n"+
			"  %s\n"+
			"  %s\n\n"+
			"  [ Placeholder — page content would load here ]\n\n"+
			"  %s\n",
		html.Colorize("OPENED RESULT", html.ColorCyan),
		html.Colorize(r.Title, html.ColorReset),
		html.Colorize(r.URL, html.ColorBlue),
		separator,
	)
}

// ── Entry point ───────────────────────────────────────────────────────────────

func HandleSearchDynamic(cmd *cobra.Command, args []string) {
	searchQuery, _ := cmd.Flags().GetString("search")
	if searchQuery == "" {
		fmt.Println("No search query provided.")
		return
	}

	engineName, _ := cmd.Flags().GetString("engine")
	var engine html.Search
	switch engineName {
	case "startpage":
		engine = search_engines.NewStartpageSearchEngine("https://www.startpage.com/sp/search?query=")
	case "mojeek":
		engine = search_engines.NewMojeekSearchEngine("https://www.mojeek.com/search?q=")
	default:
		fmt.Printf("Unknown search engine: %s\n", engineName)
		return
	}

	m := newSearchModel(engineName, searchQuery, engine)
	p := tea.NewProgram(m)

	final, err := p.Run()
	if err != nil {
		fmt.Printf("TUI error: %v\n", err)
		return
	}

	if fm, ok := final.(searchModel); ok && fm.selected != nil {

		page, httpRes, err := html.ParsePage(fm.selected.URL, connect.Get)
		if err != nil {
			fmt.Printf("Error fetching page: %v\n", err)
			return
		}

		heroText := buildWebsiteHero(httpRes, fm.selected.URL)
		if heroText != "" {
			fmt.Println(heroText)
		}

		fmt.Printf("%s\n\n",
			page,
		)
	}
}
