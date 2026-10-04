package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/api"
	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/models"
)

func TestSearchCharactersSinglePage(t *testing.T) {
	var srv *httptest.Server
	srv = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/people/" {
			t.Errorf("path = %q, want /people/", r.URL.Path)
		}
		if got := r.URL.Query().Get("search"); got != "luke sky" {
			t.Errorf("search = %q, want %q", got, "luke sky")
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page = %q, want 1", got)
		}
		writeJSON(t, w, map[string]any{
			"count": 1,
			"results": []map[string]any{{
				"name": "Luke Skywalker", "height": "172", "mass": "77",
				"birth_year": "19BBY", "homeworld": srv.URL + "/planets/1/",
				"films": []string{"ignored"},
			}},
		})
	})

	got, err := api.NewClient(srv.URL).SearchCharacters(context.Background(), "luke sky")
	if err != nil {
		t.Fatalf("SearchCharacters: %v", err)
	}
	want := []models.Character{{
		Name: "Luke Skywalker", Height: "172", Mass: "77",
		BirthYear: "19BBY", Homeworld: srv.URL + "/planets/1/",
	}}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestSearchCharactersMultiPage(t *testing.T) {
	const total, pageSize = 25, 10
	var requests atomic.Int32
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		var results []models.Character
		for i := (page-1)*pageSize + 1; i <= min(page*pageSize, total); i++ {
			results = append(results, models.Character{Name: fmt.Sprintf("Char %d", i)})
		}
		writeJSON(t, w, map[string]any{"count": total, "results": results})
	})

	got, err := api.NewClient(srv.URL).SearchCharacters(context.Background(), "char")
	if err != nil {
		t.Fatalf("SearchCharacters: %v", err)
	}
	if len(got) != total {
		t.Fatalf("got %d characters, want %d", len(got), total)
	}
	for i, c := range got {
		if want := fmt.Sprintf("Char %d", i+1); c.Name != want {
			t.Errorf("got[%d].Name = %q, want %q", i, c.Name, want)
		}
	}
	if n := requests.Load(); n != 3 {
		t.Errorf("made %d requests, want 3", n)
	}
}

func TestSearchCharactersNotFound(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{"count": 0, "results": []any{}})
	})

	_, err := api.NewClient(srv.URL).SearchCharacters(context.Background(), "nobody")
	if !errors.Is(err, api.ErrNotFound) {
		t.Errorf("err = %v, want api.ErrNotFound", err)
	}
}

func TestSearchCharactersStatusError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	_, err := api.NewClient(srv.URL).SearchCharacters(context.Background(), "luke")
	var se *api.StatusError
	if !errors.As(err, &se) || se.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("err = %v, want StatusError with status 503", err)
	}
}

func TestSearchCharactersLaterPageError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "2" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(t, w, map[string]any{"count": 11, "results": make([]models.Character, 10)})
	})

	_, err := api.NewClient(srv.URL).SearchCharacters(context.Background(), "a")
	var se *api.StatusError
	if !errors.As(err, &se) || se.StatusCode != http.StatusInternalServerError {
		t.Errorf("err = %v, want StatusError with status 500", err)
	}
}

func TestSearchCharactersCanceledContext(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not reach the server")
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := api.NewClient(srv.URL).SearchCharacters(ctx, "luke")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestGetPlanetInvalidJSON(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "not json")
	})

	_, err := api.NewClient(srv.URL).GetPlanet(context.Background(), srv.URL+"/planets/1/")
	if err == nil || !strings.Contains(err.Error(), "decode") {
		t.Errorf("err = %v, want decode error", err)
	}
}

func TestGetPlanetCachesSuccess(t *testing.T) {
	var requests atomic.Int32
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		writeJSON(t, w, models.Planet{Name: "Tatooine"})
	})
	c := api.NewClient(srv.URL)

	for range 3 {
		p, err := c.GetPlanet(context.Background(), srv.URL+"/planets/1/")
		if err != nil || p.Name != "Tatooine" {
			t.Fatalf("GetPlanet = %+v, %v", p, err)
		}
	}
	if n := requests.Load(); n != 1 {
		t.Errorf("made %d requests, want 1", n)
	}
}

func TestGetPlanetDoesNotCacheFailure(t *testing.T) {
	var requests atomic.Int32
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(t, w, models.Planet{Name: "Tatooine"})
	})
	c := api.NewClient(srv.URL)
	u := srv.URL + "/planets/1/"

	if _, err := c.GetPlanet(context.Background(), u); err == nil {
		t.Fatal("first GetPlanet: want error")
	}
	p, err := c.GetPlanet(context.Background(), u)
	if err != nil || p.Name != "Tatooine" {
		t.Errorf("second GetPlanet = %+v, %v; want Tatooine", p, err)
	}
}

func TestGetPlanetsKeepsOrderAndDedupes(t *testing.T) {
	names := map[string]string{"/planets/1/": "Tatooine", "/planets/2/": "Alderaan"}
	var mu sync.Mutex
	hits := map[string]int{}
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hits[r.URL.Path]++
		mu.Unlock()
		writeJSON(t, w, models.Planet{Name: names[r.URL.Path]})
	})

	urls := []string{srv.URL + "/planets/1/", srv.URL + "/planets/2/", srv.URL + "/planets/1/"}
	got, err := api.NewClient(srv.URL).GetPlanets(context.Background(), urls)
	if err != nil {
		t.Fatalf("GetPlanets: %v", err)
	}
	want := []string{"Tatooine", "Alderaan", "Tatooine"}
	for i := range want {
		if got[i].Name != want[i] {
			t.Errorf("got[%d].Name = %q, want %q", i, got[i].Name, want[i])
		}
	}
	if hits["/planets/1/"] != 1 || hits["/planets/2/"] != 1 {
		t.Errorf("hits = %v, want one request per planet", hits)
	}
}

func TestGetPlanetsFetchesInParallel(t *testing.T) {
	const n = 3
	var arrived sync.WaitGroup
	arrived.Add(n)
	allArrived := make(chan struct{})
	go func() { arrived.Wait(); close(allArrived) }()

	// Each handler waits until all n requests are in flight, which only
	// happens if the client sends them concurrently.
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		arrived.Done()
		select {
		case <-allArrived:
		case <-time.After(2 * time.Second):
			t.Errorf("%s waited alone: requests are not concurrent", r.URL.Path)
		}
		writeJSON(t, w, models.Planet{Name: r.URL.Path})
	})

	urls := make([]string, n)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/planets/%d/", srv.URL, i+1)
	}
	if _, err := api.NewClient(srv.URL).GetPlanets(context.Background(), urls); err != nil {
		t.Fatalf("GetPlanets: %v", err)
	}
}
