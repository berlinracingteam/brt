package brt

import (
	"html/template"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/pinub/mux"
)

// NewHandler creates a new `http.Handler` to be used to serve the content.
func NewHandler(client Client, tmpl *template.Template) http.Handler {
	m := mux.New()
	m.Get("/", index(tmpl))
	m.Get("/rennen.ics", calendar(client, tmpl))
	m.Get("/rennen", redirect("/"))
	m.Get("/team", redirect("/"))
	m.Get("/kontakt", redirect("/"))
	m.Get("/news", redirect("/"))

	h := http.NewServeMux()
	h.Handle("/css/", http.FileServer(http.Dir("./static/")))
	h.Handle("/img/", http.FileServer(http.Dir("./static/")))
	h.Handle("/", m)

	return h
}

func index(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := struct {
			Year int
		}{
			Year: time.Now().Year(),
		}
		if err := tmpl.ExecuteTemplate(w, "index.html.tmpl", data); err != nil {
			raven.CaptureError(err, nil)
		}
	})
}

func calendar(client Client, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/calendar")

		events, err := client.EventService().Events(time.Now().Year())
		if err != nil {
			raven.CaptureError(err, nil)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := tmpl.ExecuteTemplate(w, "rennen.ics.tmpl", events); err != nil {
			raven.CaptureError(err, nil)
		}
	})
}

func redirect(to string) http.Handler {
	return http.RedirectHandler(to, http.StatusMovedPermanently)
}
