package brt_test

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/berlinracingteam/brt"
	"github.com/berlinracingteam/brt/mock"
)

func newHandler(c *mock.Client) http.Handler {
	if c == nil {
		c = mock.NewClient()
	}
	t := template.Must(template.ParseGlob("./views/*.tmpl"))
	h := brt.New(c, t)

	return h
}

func TestIndex(t *testing.T) {
	h := newHandler(nil)

	t.Run("test index", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		h.ServeHTTP(rec, req)

		equals(t, 200, rec.Code)
		equals(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
		equals(t, true, strings.Contains(rec.Body.String(), "Berlin Racing Team"))
	})

	t.Run("test redirects", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/rennen", nil)
		h.ServeHTTP(rec, req)

		equals(t, 301, rec.Code)
	})
}

func TestCalendar(t *testing.T) {
	now := time.Now()
	c := mock.NewClient()
	c.Event.EventsFn = func(year int) ([]brt.Event, error) {
		return []brt.Event{
			brt.Event{
				ID:        42,
				Title:     "Rund um die Wurst",
				Date:      &now,
				CreatedAt: &now,
				CreatedBy: &brt.Person{
					Name:  "Anonymous",
					Email: "anon@example.com",
				},
				Website:  "",
				Distance: 42,
				People:   []*brt.Person{},
			},
		}, nil
	}
	h := newHandler(c)

	t.Run("test calendar", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/rennen.ics", nil)
		h.ServeHTTP(rec, req)

		equals(t, 200, rec.Code)
		equals(t, "text/calendar", rec.Header().Get("Content-Type"))
		equals(t, true, strings.Contains(rec.Body.String(), "Rund um die Wurst"))
	})
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
