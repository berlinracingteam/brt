package pg

import (
	"database/sql"

	"github.com/berlinracingteam/brt"
)

const driver = "postgres"

// Ensure Client implements brt.Client
var _ brt.Client = &Client{}

// Client contains the links to the database and the current implemented
// services.
type Client struct {
	*sql.DB

	eventService EventService
}

// New creates and returns a new Client struct.
func New(dataSource string) (*Client, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	c := &Client{DB: db}
	c.eventService.Client = c

	return c, nil
}

// EventService returns the linked service.
func (c *Client) EventService() brt.EventService {
	return &c.eventService
}
