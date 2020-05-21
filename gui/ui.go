package gui

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var SystemCursorEnabled = true

type Base struct {
	EventListeners EventListener
}

func SetSystemCursorEnabled(enabled bool) {
	SystemCursorEnabled = enabled
	if SystemCursorEnabled {
		engo.SetCursor(engo.CursorArrow)
	} else {
		engo.SetCursor(common.CursorNone)
		engo.SetCursorVisibility(false)
	}
}

type EventListener map[string][]func()

func (b *Base) DispatchEvents(event engo.Message) {
	for _, f := range b.EventListeners[event.Type()] {
		f()
	}
}
