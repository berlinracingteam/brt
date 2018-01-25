package mock

import "github.com/berlinracingteam/brt"

// Ensure EventService implements brt.EventService
var _ brt.EventService = &EventService{}

// EventService description.
type EventService struct {
	Client *Client

	EventsFn func(int) ([]brt.Event, error)

	EventsInvoked bool
}

// Events description.
func (s *EventService) Events(year int) ([]brt.Event, error) {
	s.EventsInvoked = true
	return s.EventsFn(year)
}
