package gui

type EventMouseOver struct{}

func (EventMouseOver) Type() string {
	return "mouse_over"
}

type EventMouseClicked struct{}

func (EventMouseClicked) Type() string {
	return "mouse_clicked"
}
