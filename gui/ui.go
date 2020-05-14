package gui

import (
	"github.com/EngoEngine/engo"
)

type Base struct {
	EventListeners EventListener
}

type EventListener map[string][]func()

func (b *Base) DispatchEvents(event engo.Message) {
	for _, f := range b.EventListeners[event.Type()] {
		f()
	}
}
