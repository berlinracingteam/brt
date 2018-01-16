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
	mux := mux.New()
	mux.Handler("GET", "/", index(tmpl))
	mux.Handler("GET", "/rennen.ics", calendar(client, tmpl))
	mux.Handler("GET", "/rennen", redirect("/"))
	mux.Handler("GET", "/team", redirect("/"))
	mux.Handler("GET", "/kontakt", redirect("/"))
	mux.Handler("GET", "/news", redirect("/"))

	h := http.NewServeMux()
	h.Handle("/css/", http.FileServer(http.Dir("./static/")))
	h.Handle("/img/", http.FileServer(http.Dir("./static/")))
	h.Handle("/", mux)

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
