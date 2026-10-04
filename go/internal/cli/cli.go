// Package cli parses arguments and prints SWAPI search results.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jkapsalis/SwAPI-cli-search-tool/go/internal/api"
	"github.com/jkapsalis/SwAPI-cli-search-tool/go/internal/models"
)

const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2

	earthDayHours = 24
	earthYearDays = 365
)

const usage = `Usage: swapi search <name> [--world]

Search Star Wars characters by name.

Flags:
  --world   also show each character's homeworld
`

// Run executes the CLI and returns the process exit code.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer, client *api.Client) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return exitUsage
	}
	switch args[0] {
	case "search":
		return runSearch(ctx, args[1:], stdout, stderr, client)
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return exitOK
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
		return exitUsage
	}
}

func runSearch(ctx context.Context, args []string, stdout, stderr io.Writer, client *api.Client) int {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	world := fs.Bool("world", false, "show homeworld info")

	// flag stops at the first positional argument; keep parsing so
	// --world works both before and after the name.
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				fmt.Fprint(stdout, usage)
				return exitOK
			}
			fmt.Fprintf(stderr, "%v\n\n%s", err, usage)
			return exitUsage
		}
		args = fs.Args()
		if len(args) == 0 {
			break
		}
		positional = append(positional, args[0])
		args = args[1:]
	}

	if len(positional) != 1 || strings.TrimSpace(positional[0]) == "" {
		fmt.Fprintf(stderr, "search needs exactly one non-empty name (quote names with spaces)\n\n%s", usage)
		return exitUsage
	}

	chars, err := client.SearchCharacters(ctx, strings.TrimSpace(positional[0]))
	if errors.Is(err, api.ErrNotFound) {
		fmt.Fprintln(stderr, "Character not found.")
		return exitError
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return exitError
	}

	var planets []models.Planet
	if *world {
		urls := make([]string, len(chars))
		for i, c := range chars {
			urls[i] = c.Homeworld
		}
		planets, err = client.GetPlanets(ctx, urls)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return exitError
		}
	}

	for i, c := range chars {
		if i > 0 {
			fmt.Fprintln(stdout, "----------------------------------------")
		}
		fmt.Fprint(stdout, formatCharacter(c))
		if *world {
			fmt.Fprint(stdout, "\n", formatHomeworld(planets[i]), "\n", timeRatio(planets[i]))
		}
	}
	return exitOK
}

func formatCharacter(c models.Character) string {
	return fmt.Sprintf("Name: %s\nHeight: %s cm\nMass: %s kg\nBirth Year: %s\n",
		c.Name, c.Height, c.Mass, c.BirthYear)
}

func formatHomeworld(p models.Planet) string {
	return fmt.Sprintf("Homeworld: %s\nPopulation: %s\nRotation Period: %s hours\nOrbital Period: %s days\n",
		p.Name, p.Population, p.RotationPeriod, p.OrbitalPeriod)
}

// timeRatio compares a planet's day and year length to Earth's.
func timeRatio(p models.Planet) string {
	rotation, err1 := strconv.ParseFloat(p.RotationPeriod, 64)
	orbital, err2 := strconv.ParseFloat(p.OrbitalPeriod, 64)
	if err1 != nil || err2 != nil {
		return "Unknown time ratio\n"
	}
	return fmt.Sprintf("Day vs Earth: %.2fx\nYear vs Earth: %.2fx\n",
		rotation/earthDayHours, orbital/earthYearDays)
}
