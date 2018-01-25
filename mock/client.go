package mock

import "github.com/berlinracingteam/brt"

// Ensure Client implements brt.Client
var _ brt.Client = &Client{}

// Client description.
type Client struct {
	Event EventService
}

// NewClient description.
func NewClient() *Client {
	c := &Client{}
	c.Event.Client = c

	return c
}

// EventService description.
func (c *Client) EventService() brt.EventService {
	return &c.Event
}
