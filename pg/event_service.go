package pg

import (
	"database/sql"

	"github.com/berlinracingteam/brt"
)

// Ensure EventService implements brt.EventService
var _ brt.EventService = &EventService{}

type EventService struct {
	Client *Client
}

func (s *EventService) Events(from int) ([]brt.Event, error) {
	query := `
        SELECT
                e.id, e.title, e.date, e.created_at, e.url, e.distance,
                p.first_name || ' ' || p.last_name AS name, p.email,
                o.first_name || ' ' || o.last_name AS oname, o.email AS oemail
        FROM
                events AS e
                LEFT JOIN people AS o ON e.person_id = o.id
                LEFT JOIN participations AS t ON e.id = t.event_id
                LEFT JOIN people AS p ON t.person_id = p.id
        WHERE
                e.date > make_Date($1, 1, 1)
        ORDER BY
                e.date, e.id, p.first_name, p.last_name
	`
	rows, err := s.Client.Query(query, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // nolint: errcheck

	var events []brt.Event
	var lastEvent brt.Event
	var name, email sql.NullString

	for rows.Next() {
		var event brt.Event
		var person brt.Person

		if err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.CreatedAt,
			&event.Website,
			&event.Distance,
			&name,
			&email,
			&person.Name,
			&person.Email,
		); err != nil {
			return nil, err
		}
		event.CreatedBy = &person

		if lastEvent.ID == 0 {
			lastEvent = event
		}

		if event.ID != lastEvent.ID {
			events = append(events, lastEvent)
			lastEvent = event
		}

		if name.Valid && email.Valid {
			lastEvent.People = append(lastEvent.People, &brt.Person{
				Name:  name.String,
				Email: email.String,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
