package connect

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GetFunc func(url string, body []byte, headers map[string]string) (*HttpResponse, error)

type CacheEntry struct {
	ValidUntil time.Time     `json:"valid_until"`
	Response   *HttpResponse `json:"response"`
}

type FileCache struct {
	CacheDir string
}

func NewFileCache(dir string) *FileCache {
	os.MkdirAll(dir, os.ModePerm)
	return &FileCache{
		CacheDir: dir,
	}
}

func (c *FileCache) WithCache(next GetFunc) GetFunc {
	return func(url string, body []byte, headers map[string]string) (*HttpResponse, error) {

		cachePath := c.getCachePath(url)
		if cachedResp := c.tryGet(cachePath); cachedResp != nil {
			return cachedResp, nil
		}

		resp, err := next(url, body, headers)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {
			c.doCache(cachePath, resp)
		}

		return resp, nil
	}
}

func (c *FileCache) getCachePath(url string) string {
	cleanURL := simplifyURL(url)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(cleanURL)))
	return filepath.Join(c.CacheDir, hash+".json")
}

func (c *FileCache) tryGet(cacheFile string) *HttpResponse {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil
	}

	// Check if expired
	if time.Now().After(entry.ValidUntil) {
		os.Remove(cacheFile) // Clean up expired cache
		return nil
	}

	slog.Info(fmt.Sprintf("Cache hit for %s, valid until %s", cacheFile, entry.ValidUntil.Format(time.RFC1123)))
	return entry.Response
}

func (c *FileCache) doCache(cacheFile string, resp *HttpResponse) {
	headers := resp.Headers

	// Default: do not cache unless we find a directive
	var duration time.Duration = 0

	if val, ok := headers["Cache-Control"]; ok {
		// Look for max-age=<seconds>
		if strings.Contains(val, "no-store") || strings.Contains(val, "private") {
			return // Respect "no-store" or "private" by not caching
		}

		// Use a regex or string splitting to find max-age
		re := regexp.MustCompile(`max-age=(\d+)`)
		matches := re.FindStringSubmatch(val)
		if len(matches) > 1 {
			seconds, _ := strconv.Atoi(matches[1])
			duration = time.Duration(seconds) * time.Second
		}
	}

	// Fallback to Expires header if max-age is missing
	if duration == 0 {
		if expVal, ok := headers["Expires"]; ok {
			if expTime, err := http.ParseTime(expVal); err == nil {
				duration = time.Until(expTime)
			}
		}
	}

	if duration <= 0 {
		return // Don't cache if no valid duration is found
	}

	entry := CacheEntry{
		ValidUntil: time.Now().Add(duration),
		Response:   resp,
	}

	if data, err := json.Marshal(entry); err == nil {
		os.WriteFile(cacheFile, data, 0644)
	}
}

func simplifyURL(url string) string {
	prefix := "://"
	if idx := strings.Index(url, prefix); idx != -1 {
		url = url[idx+len(prefix):]
	}
	url = strings.TrimSuffix(url, "/")
	return strings.ToLower(url)
}
