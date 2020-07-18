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
	RockColor       int
	RockType        int
	RockSize        int
	ProjectileColor int
)

var (
	ShipColors = []ShipColor{ShipColorRed, ShipColorBlue, ShipColorOrange}
	RockColors = []RockColor{RockColorBrown, RockColorGrey}
	RockTypes  = []RockType{RockTypeOne, RockTypeTwo, RockTypeThree, RockTypeFour}
	RockSizes  = []RockSize{RockSizeBig, RockSizeMedium, RockSizeSmall, RockSizeTiny}
)

const (
	ShipColorRed ShipColor = iota
	ShipColorBlue
	ShipColorOrange
)

const (
	RockColorBrown RockColor = iota
	RockColorGrey
)

const (
	RockTypeOne RockType = iota + 1
	RockTypeTwo
	RockTypeThree
	RockTypeFour
)

const (
	RockSizeBig RockSize = iota
	RockSizeMedium
	RockSizeSmall
	RockSizeTiny
)

var GlobalScale = engo.Point{X: 0.75, Y: 0.75}

type ScoreChanged struct{}

func (ScoreChanged) Type() string { return "ScoreChanged" }

func (s ShipColor) String() string {
	return [...]string{"red", "blue", "orange"}[s]
}

func (r RockColor) String() string {
	return [...]string{"brown", "grey", "orange"}[r]
}

func (r RockSize) String() string {
	return [...]string{"big", "medium", "small", "tiny"}[r]
}

func (r RockType) String() string {
	return [...]string{"1", "2", "3", "4"}[r]
}
