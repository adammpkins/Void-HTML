package events

// ShoutEvent is an interface representing the data needed for a shout event.
type ShoutEvent interface {
	GetShoutID() uint
	GetContent() string
	GetUserID() uint
	GetUsername() string
}
