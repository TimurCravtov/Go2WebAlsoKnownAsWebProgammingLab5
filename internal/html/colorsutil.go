package html

import "regexp"

func colorizeURLs(text string) string {

	reURL := regexp.MustCompile(`https?://[^\s)]+`)

	coloredText := reURL.ReplaceAllStringFunc(text, func(match string) string {
		return Colorize(match, ColorBlue)
	})
	return coloredText
}

const (
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorReset  = "\033[0m" // Required to revert to default terminal text
)

func Colorize(text string, colorCode string) string {
	return colorCode + text + ColorReset
}
