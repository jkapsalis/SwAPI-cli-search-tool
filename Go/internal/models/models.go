// Package models defines the SWAPI resources used by the CLI.
package models

// Character is a SWAPI person. Numeric fields stay strings because the API
// returns values such as "unknown" or "1,358".
type Character struct {
	Name      string `json:"name"`
	Height    string `json:"height"`
	Mass      string `json:"mass"`
	BirthYear string `json:"birth_year"`
	Homeworld string `json:"homeworld"`
}

// Planet is a SWAPI planet, used as a character's homeworld.
type Planet struct {
	Name           string `json:"name"`
	Population     string `json:"population"`
	RotationPeriod string `json:"rotation_period"`
	OrbitalPeriod  string `json:"orbital_period"`
}
