package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/api"
	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/cli"
	"github.com/jkapsalis/SwAPI-cli-search-tool/Go/internal/models"
)

// Exit codes documented in the README.
const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

// newFakeSWAPI serves a tiny subset of SWAPI: searches for "luke", "sky",
// "nobody" and "boom", plus planet 1.
func newFakeSWAPI(t *testing.T) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		tatooine := srv.URL + "/planets/1/"
		luke := models.Character{Name: "Luke Skywalker", Height: "172", Mass: "77", BirthYear: "19BBY", Homeworld: tatooine}
		anakin := models.Character{Name: "Anakin Skywalker", Height: "188", Mass: "84", BirthYear: "41.9BBY", Homeworld: tatooine}

		var body any
		switch {
		case r.URL.Path == "/planets/1/":
			body = models.Planet{Name: "Tatooine", Population: "200000", RotationPeriod: "23", OrbitalPeriod: "304"}
		case r.URL.Path == "/people/":
			switch r.URL.Query().Get("search") {
			case "luke":
				body = map[string]any{"count": 1, "results": []models.Character{luke}}
			case "sky":
				body = map[string]any{"count": 2, "results": []models.Character{luke, anakin}}
			case "boom":
				w.WriteHeader(http.StatusInternalServerError)
				return
			default:
				body = map[string]any{"count": 0, "results": []models.Character{}}
			}
		default:
			http.NotFound(w, r)
			return
		}
		writeJSON(t, w, body)
	})
	return srv
}

func runCLI(t *testing.T, srv *httptest.Server, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = cli.Run(context.Background(), args, &out, &errOut, api.NewClient(srv.URL))
	return code, out.String(), errOut.String()
}

func run(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	return runCLI(t, newFakeSWAPI(t), args...)
}

func TestRunWorldOutput(t *testing.T) {
	code, stdout, stderr := run(t, "search", "luke", "--world")
	if code != exitOK {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	want := `Name: Luke Skywalker
Height: 172 cm
Mass: 77 kg
Birth Year: 19BBY

Homeworld: Tatooine
Population: 200000
Rotation Period: 23 hours
Orbital Period: 304 days

Day vs Earth: 0.96x
Year vs Earth: 0.83x
`
	if stdout != want {
		t.Errorf("stdout =\n%s\nwant\n%s", stdout, want)
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout []string
		notStdout  []string
		wantStderr string
	}{
		{
			name:       "search without world",
			args:       []string{"search", "luke"},
			wantCode:   exitOK,
			wantStdout: []string{"Name: Luke Skywalker", "Height: 172 cm"},
			notStdout:  []string{"Homeworld:"},
		},
		{
			name:       "world flag before name",
			args:       []string{"search", "--world", "luke"},
			wantCode:   exitOK,
			wantStdout: []string{"Homeworld: Tatooine", "Day vs Earth: 0.96x"},
		},
		{
			name:       "single dash flag",
			args:       []string{"search", "luke", "-world"},
			wantCode:   exitOK,
			wantStdout: []string{"Homeworld: Tatooine"},
		},
		{
			name:       "name is trimmed",
			args:       []string{"search", "  luke  "},
			wantCode:   exitOK,
			wantStdout: []string{"Name: Luke Skywalker"},
		},
		{
			name:       "multiple matches",
			args:       []string{"search", "sky", "--world"},
			wantCode:   exitOK,
			wantStdout: []string{"Name: Luke Skywalker", "----", "Name: Anakin Skywalker", "Homeworld: Tatooine"},
		},
		{
			name:       "character not found",
			args:       []string{"search", "nobody"},
			wantCode:   exitError,
			wantStderr: "Character not found.",
		},
		{
			name:       "api error",
			args:       []string{"search", "boom"},
			wantCode:   exitError,
			wantStderr: "unexpected status 500",
		},
		{
			name:       "no arguments",
			args:       nil,
			wantCode:   exitUsage,
			wantStderr: "Usage:",
		},
		{
			name:       "unknown command",
			args:       []string{"find", "luke"},
			wantCode:   exitUsage,
			wantStderr: `unknown command "find"`,
		},
		{
			name:       "missing name",
			args:       []string{"search"},
			wantCode:   exitUsage,
			wantStderr: "exactly one non-empty name",
		},
		{
			name:       "blank name",
			args:       []string{"search", "   "},
			wantCode:   exitUsage,
			wantStderr: "exactly one non-empty name",
		},
		{
			name:       "unquoted two-word name",
			args:       []string{"search", "luke", "sky"},
			wantCode:   exitUsage,
			wantStderr: "exactly one non-empty name",
		},
		{
			name:       "unknown flag",
			args:       []string{"search", "luke", "--nope"},
			wantCode:   exitUsage,
			wantStderr: "flag provided but not defined",
		},
		{
			name:       "help command",
			args:       []string{"help"},
			wantCode:   exitOK,
			wantStdout: []string{"Usage:"},
		},
		{
			name:       "search help flag",
			args:       []string{"search", "-h"},
			wantCode:   exitOK,
			wantStdout: []string{"Usage:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := run(t, tt.args...)
			if code != tt.wantCode {
				t.Errorf("exit code = %d, want %d (stderr = %q)", code, tt.wantCode, stderr)
			}
			for _, s := range tt.wantStdout {
				if !strings.Contains(stdout, s) {
					t.Errorf("stdout missing %q:\n%s", s, stdout)
				}
			}
			for _, s := range tt.notStdout {
				if strings.Contains(stdout, s) {
					t.Errorf("stdout should not contain %q:\n%s", s, stdout)
				}
			}
			if !strings.Contains(stderr, tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tt.wantStderr)
			}
		})
	}
}

// TestTimeRatio checks the day and year ratios through the --world output,
// since the ratio helper is unexported.
func TestTimeRatio(t *testing.T) {
	tests := []struct {
		name   string
		planet models.Planet
		want   string
	}{
		{"tatooine", models.Planet{RotationPeriod: "23", OrbitalPeriod: "304"}, "Day vs Earth: 0.96x\nYear vs Earth: 0.83x\n"},
		{"earth-like", models.Planet{RotationPeriod: "24", OrbitalPeriod: "365"}, "Day vs Earth: 1.00x\nYear vs Earth: 1.00x\n"},
		{"zero periods", models.Planet{RotationPeriod: "0", OrbitalPeriod: "0"}, "Day vs Earth: 0.00x\nYear vs Earth: 0.00x\n"},
		{"unknown rotation", models.Planet{RotationPeriod: "unknown", OrbitalPeriod: "304"}, "Unknown time ratio\n"},
		{"empty orbital", models.Planet{RotationPeriod: "23", OrbitalPeriod: ""}, "Unknown time ratio\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var srv *httptest.Server
			srv = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/planets/1/" {
					writeJSON(t, w, tt.planet)
					return
				}
				writeJSON(t, w, map[string]any{
					"count":   1,
					"results": []models.Character{{Name: "Tester", Homeworld: srv.URL + "/planets/1/"}},
				})
			})

			code, stdout, stderr := runCLI(t, srv, "search", "tester", "--world")
			if code != exitOK {
				t.Fatalf("exit code = %d, stderr = %q", code, stderr)
			}
			if !strings.HasSuffix(stdout, "\n"+tt.want) {
				t.Errorf("stdout =\n%s\nwant it to end with\n%s", stdout, tt.want)
			}
		})
	}
}
