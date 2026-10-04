# SWAPI CLI Search Tool

A command-line tool that searches Star Wars characters and their homeworlds using the [Star Wars API (SWAPI)](https://swapi.dev).

The project started as a Python script and has been rebuilt in Go. The Go version is the current one. The original Python code is kept in [`Python_old_Version/`](Python_old_Version) for reference.

## Repository layout

| Folder | Contents |
|---|---|
| [`Go/`](Go) | Current version: CLI source, unit tests, build scripts |
| [`Python_old_Version/`](Python_old_Version) | Original Python prototype (no longer maintained) |

## Quick start

Requires [Go 1.25+](https://go.dev/dl/).

Install the binary:

```bash
go install github.com/jkapsalis/SwAPI-cli-search-tool/Go/cmd/swapi@latest
swapi search "luke sky" --world
```

Or build from source:

```bash
git clone https://github.com/jkapsalis/SwAPI-cli-search-tool.git
cd SwAPI-cli-search-tool/Go
make build                              # creates bin/swapi
./bin/swapi search "luke sky" --world
```

## Usage

```text
swapi search <name> [--world]
```

| Argument | Description |
|---|---|
| `<name>` | Full or partial character name. Quote names that contain spaces. |
| `--world` | Also show each character's homeworld and compare its day and year length to Earth's. Can be placed before or after the name. |

Example:

```text
$ swapi search "luke sky" --world
Name: Luke Skywalker
Height: 172 cm
Mass: 77 kg
Birth Year: 19BBY

Homeworld: Tatooine
Population: 200000
Rotation Period: 23 hours
Orbital Period: 304 days

Day vs Earth: 0.96x
Year vs Earth: 0.83x
```

A partial name such as `sky` returns every match (Luke, Anakin and Shmi Skywalker).

| Exit code | Meaning |
|---|---|
| `0` | Success |
| `1` | Character not found, or SWAPI returned an error |
| `2` | Invalid input (unknown command, missing name, unknown flag) |

## Why the remake

The Python version worked well as a first prototype: it was quick to write and easy to read. Moving it toward a tool other people can install and rely on exposed some limits:

- Running it required a Python installation plus `pip install -r requirements.txt`, which pins 16 packages even though only `requests` is used.
- It showed only the first match from the first page of results.
- Every HTTP request ran one after another.
- API data was handled as untyped dictionaries, so a wrong field name failed only at runtime.
- An API failure ended in a Python traceback, and there were no tests.

The remake treats the project as a small piece of production software instead of a script: clear package boundaries, typed data, explicit error handling, automated tests, and reproducible builds.

## Why Go over Python

| | Python (old) | Go (current) |
|---|---|---|
| Distribution | Python interpreter plus pip packages | One static binary, nothing else to install |
| Dependencies | 16 pinned packages | Go standard library only |
| Search results | First match only | All matches across all result pages |
| Concurrency | Sequential requests | Result pages and homeworlds fetched in parallel; repeated requests served from a cache |
| Data model | Dictionaries | Typed `Character` and `Planet` structs |
| Error handling | Generic exceptions, bare `except` | Typed errors and distinct exit codes |
| Testing | None | Unit tests run with the race detector, 97% statement coverage |
| Platforms | Wherever Python is installed | Cross-compiled for Linux and macOS (amd64, arm64) and Windows (amd64) |

The main reasons for choosing Go:

1. **Simple distribution.** Go compiles to a single binary with no runtime, so users download one file and run it.
2. **Built-in concurrency.** Goroutines make it straightforward to fetch several pages and planets at the same time.
3. **Static typing.** JSON responses are decoded into structs, so many mistakes are caught when the code compiles.
4. **Explicit errors.** Every failure is returned as a value and handled on purpose, which leads to clear messages and exit codes.
5. **Tooling included.** Formatting, vetting, testing, the race detector and cross-compilation all ship with the Go toolchain.

Python is still a strong choice for quick scripts and data work. For a small CLI that should be easy to install and safe to run anywhere, Go is the better fit.

## Project structure

```text
Go/
├── cmd/swapi/          Entry point: wires the API client into the CLI
├── internal/cli/       Argument parsing and output formatting
├── internal/api/       SWAPI client: HTTP calls, parallel fetching, cache
├── internal/models/    Character and Planet types
├── tests/              Unit tests for the api and cli packages
└── Makefile            build, test and release targets
```

Each package has one job, and the CLI receives its API client as a parameter. This lets the tests point the CLI at a local fake SWAPI server instead of the real API.

## Development

Run these from the `Go/` folder:

| Command | What it does |
|---|---|
| `make build` | Builds a static binary for your platform in `bin/` |
| `make test` | Runs all unit tests offline with the race detector |
| `make release` | Builds binaries for Linux, macOS and Windows in `dist/` |
| `make clean` | Removes `bin/` and `dist/` |

The tests live in `Go/tests/` and use Go's `testing` package and `net/http/httptest`. They test each package only through its exported API, the same way a user of the package would. They cover search across several pages, not-found and HTTP errors, request cancellation, caching, parallel fetching, CLI flags, exit codes and the exact output format.

## Roadmap

All goals of the Go rewrite are complete:

- [x] CLI in Go using the standard `flag` package, with the same `search` command and `--world` flag
- [x] HTTP requests with `net/http`
- [x] Typed structs for characters and planets
- [x] Parallel fetching with goroutines and an in-memory cache
- [x] Explicit error handling for not found, API errors and invalid input
- [x] Modular packages: `api`, `models`, `cli`
- [x] Standalone cross-platform binaries
- [x] Unit tests for API calls and search
