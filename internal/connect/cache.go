package connect

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	Duration time.Duration
}

func NewFileCache(dir string, duration time.Duration) *FileCache {
	os.MkdirAll(dir, os.ModePerm)
	return &FileCache{
		CacheDir: dir,
		Duration: duration,
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

	return entry.Response
}

func (c *FileCache) doCache(cacheFile string, resp *HttpResponse) {
	entry := CacheEntry{
		ValidUntil: time.Now().Add(c.Duration),
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
