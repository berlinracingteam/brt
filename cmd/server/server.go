package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/berlinracingteam/brt"
	"github.com/berlinracingteam/brt/pg"
	"github.com/getsentry/raven-go"
	_ "github.com/lib/pq"
)

func main() {
	client, err := pg.New(os.Getenv("DATABASE"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	defer client.Close() // nolint: errcheck

	tmpl := template.Must(template.ParseGlob("./views/*.tmpl"))
	handler := brt.New(client, tmpl)
	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    0,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
