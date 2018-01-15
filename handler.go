package brt

import (
	"html/template"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/pinub/mux"
)

func NewHandler(client Client, tmpl *template.Template) http.Handler {
	mux := mux.New()
	mux.Handler("GET", "/", index(tmpl))
	mux.Handler("GET", "/rennen.ics", calendar(client, tmpl))
	mux.Handler("GET", "/rennen", http.RedirectHandler("/", http.StatusMovedPermanently))
	mux.Handler("GET", "/team", http.RedirectHandler("/", http.StatusMovedPermanently))
	mux.Handler("GET", "/kontakt", http.RedirectHandler("/", http.StatusMovedPermanently))
	mux.Handler("GET", "/news", http.RedirectHandler("/", http.StatusMovedPermanently))

	h := http.NewServeMux()
	h.Handle("/css/", http.FileServer(http.Dir("./static/")))
	h.Handle("/img/", http.FileServer(http.Dir("./static/")))
	h.Handle("/", mux)

	return h
}

func index(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl.ExecuteTemplate(w, "index.html.tmpl", struct {
			Year int
		}{
			Year: time.Now().Year(),
		})
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

		tmpl.ExecuteTemplate(w, "rennen.ics.tmpl", events)
	})
}
