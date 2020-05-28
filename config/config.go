package config

import "github.com/EngoEngine/engo"

var GlobalScale = engo.Point{X: 0.75, Y: 0.75}
var SpeedY = 400 * GlobalScale.Y

type ScoreChanged struct {}
func (ScoreChanged) Type() string { return "ScoreChanged" }