package pg

import (
	"database/sql"

	"github.com/berlinracingteam/brt"
)

const driver = "postgres"

// Ensure Client implements brt.Client
var _ brt.Client = &Client{}

type Client struct {
	*sql.DB

	eventService EventService
}

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

func (c *Client) EventService() brt.EventService {
	return &c.eventService
}
