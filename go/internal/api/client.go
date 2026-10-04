// Package api is a small SWAPI client with an in-memory response cache.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jkapsalis/SwAPI-cli-search-tool/go/internal/models"
)

const DefaultBaseURL = "https://swapi.dev/api"

// ErrNotFound is returned when a search has no matches.
var ErrNotFound = errors.New("character not found")

// StatusError is returned when SWAPI answers with a non-200 status.
type StatusError struct {
	URL        string
	StatusCode int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("swapi: GET %s: unexpected status %d", e.URL, e.StatusCode)
}

type peoplePage struct {
	Count   int                `json:"count"`
	Results []models.Character `json:"results"`
}

// cacheEntry lets concurrent callers for the same URL share one request.
type cacheEntry struct {
	once sync.Once
	body []byte
	err  error
}

type Client struct {
	baseURL string
	http    *http.Client

	mu    sync.Mutex
	cache map[string]*cacheEntry
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 10 * time.Second},
		cache:   make(map[string]*cacheEntry),
	}
}

// SearchCharacters returns every character whose name matches. Pages after
// the first are fetched in parallel.
func (c *Client) SearchCharacters(ctx context.Context, name string) ([]models.Character, error) {
	first, err := c.searchPage(ctx, name, 1)
	if err != nil {
		return nil, err
	}
	if first.Count == 0 || len(first.Results) == 0 {
		return nil, ErrNotFound
	}

	pageSize := len(first.Results)
	pages := (first.Count + pageSize - 1) / pageSize
	results := make([][]models.Character, pages)
	errs := make([]error, pages)
	results[0] = first.Results

	var wg sync.WaitGroup
	for p := 2; p <= pages; p++ {
		wg.Go(func() {
			page, err := c.searchPage(ctx, name, p)
			results[p-1], errs[p-1] = page.Results, err
		})
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	var all []models.Character
	for _, r := range results {
		all = append(all, r...)
	}
	return all, nil
}

// GetPlanet fetches a planet by its full SWAPI URL.
func (c *Client) GetPlanet(ctx context.Context, planetURL string) (models.Planet, error) {
	var p models.Planet
	err := c.getJSON(ctx, planetURL, &p)
	return p, err
}

// GetPlanets fetches planets in parallel and returns them in input order.
// Repeated URLs are requested once thanks to the cache.
func (c *Client) GetPlanets(ctx context.Context, planetURLs []string) ([]models.Planet, error) {
	planets := make([]models.Planet, len(planetURLs))
	errs := make([]error, len(planetURLs))

	var wg sync.WaitGroup
	for i, u := range planetURLs {
		wg.Go(func() {
			planets[i], errs[i] = c.GetPlanet(ctx, u)
		})
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return planets, nil
}

func (c *Client) searchPage(ctx context.Context, name string, page int) (peoplePage, error) {
	q := url.Values{"search": {name}, "page": {strconv.Itoa(page)}}
	var p peoplePage
	err := c.getJSON(ctx, c.baseURL+"/people/?"+q.Encode(), &p)
	return p, err
}

func (c *Client) getJSON(ctx context.Context, u string, v any) error {
	body, err := c.get(ctx, u)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("swapi: decode %s: %w", u, err)
	}
	return nil
}

func (c *Client) get(ctx context.Context, u string) ([]byte, error) {
	c.mu.Lock()
	e, ok := c.cache[u]
	if !ok {
		e = &cacheEntry{}
		c.cache[u] = e
	}
	c.mu.Unlock()

	e.once.Do(func() { e.body, e.err = c.fetch(ctx, u) })

	if e.err != nil {
		// Drop failed entries so a later call can retry.
		c.mu.Lock()
		if c.cache[u] == e {
			delete(c.cache, u)
		}
		c.mu.Unlock()
	}
	return e.body, e.err
}

func (c *Client) fetch(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("swapi: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("swapi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &StatusError{URL: u, StatusCode: resp.StatusCode}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("swapi: read %s: %w", u, err)
	}
	return body, nil
}
