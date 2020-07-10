package config

import "github.com/EngoEngine/engo"

type (
	// ShipConfig represents the settings that change the look and functionality of the ship
	ShipConfig struct {
		// Color of the ship (red, blue, orange)
		Color ShipColor
		// Type of the ship (1, 2, 3)
		Type string
	}

	GameConfig struct {
		VSync           bool
		Sound           bool
		BackgroundColor string
		ShipStyle       ShipConfig
	}

	ShipColor       int
	ProjectileColor int
)

const (
	ShipColorRed ShipColor = iota
	ShipColorBlue
	ShipColorOrange
)

var GlobalScale = engo.Point{X: 0.75, Y: 0.75}

type ScoreChanged struct{}

func (ScoreChanged) Type() string { return "ScoreChanged" }

func (s ShipColor) String() string {
	return [...]string{"red", "blue", "orange"}[s]
}
