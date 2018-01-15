package brt

import (
	"fmt"
	"time"
)

type EventService interface {
	Events(int) ([]Event, error)
}

type Client interface {
	EventService() EventService
}

// Attendee defines an "attendee" within a calendar component.
type Attendee interface {
	CN() string
}

// Ensure Person implements Attendee
var _ Attendee = &Person{}

// Person contains the data structure of a single person.
type Person struct {
	Name  string
	Email string
}

func (p *Person) CN() string {
	return fmt.Sprintf("%s:mailto:%s", p.Name, p.Email)
}

// VEvent provides the grouping of component properties that describe the
// event.
type VEvent interface {
	DTStamp() string
	DTStart() string
	DTEnd() string
	Summary() string
	URL() string
	Organizer() *Person
	Attendees() []*Person
}

// Ensure Event implements VEvent
var _ VEvent = &Event{}

// Event contains the structure for an event in the database.
type Event struct {
	ID        int
	Title     string
	Date      *time.Time
	CreatedAt *time.Time
	CreatedBy *Person
	Website   string
	Distance  int
	People    []*Person
}

// UID defines the persistent, globally unique identifier for the calendar
// component.
func (e *Event) UID() string {
	return fmt.Sprintf("%s-%d", e.DTStamp(), e.ID)
}

// DTStamp specifies the date and time the instance of the iCalendar object
// was created.
func (e *Event) DTStamp() string {
	return e.CreatedAt.Format("20060102T150405Z")
}

// DTStart returns the start event date formatted for the use in .ics calendar
// format.
func (e *Event) DTStart() string {
	return e.Date.Format("20060102")
}

// DTEnd returns the next day of the event. For use in .ics calendar format
// to mark the ending date.
func (e *Event) DTEnd() string {
	return e.Date.AddDate(0, 0, 1).Format("20060102")
}

// Summary defines a short summary for the calendar component.
func (e *Event) Summary() string {
	return fmt.Sprintf("%s\\, %dkm", e.Title, e.Distance)
}

// URL defines an URL associated with the iCalendar object.
func (e *Event) URL() string {
	return e.Website
}

// Organizer defines th eorganizer of the calendar component.
func (e *Event) Organizer() *Person {
	return e.CreatedBy
}

// Attendees defines a list of "attendees" within a calendar component.
func (e *Event) Attendees() []*Person {
	return e.People
}
