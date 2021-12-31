package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/29-FYI/twentynine"
	"github.com/go-chi/chi/v5"
)

type twentyninefyi struct {
	cli    *twentynine.Client
	logger *log.Logger
}

func (tnfyi *twentyninefyi) IndexHandler(w http.ResponseWriter, r *http.Request) {
	links, err := tnfyi.cli.Links()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, links)
}

func (tnfyi *twentyninefyi) StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))).ServeHTTP(w, r)
}

func (tnfyi *twentyninefyi) FormHandler(w http.ResponseWriter, r *http.Request) {
	headline := r.PostFormValue("headline")
	url := r.PostFormValue("url")

	if err := tnfyi.cli.LinkLink(twentynine.Link{
		Headline: headline,
		URL:      url,
	}); err != nil {
		if err, ok := err.(twentynine.Error); ok {
			http.Error(w, err.Message, err.Code)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (tnfyi *twentyninefyi) Handler() http.Handler {
	rtr := chi.NewRouter()
	rtr.Get("/", tnfyi.IndexHandler)
	rtr.Get("/static/*", tnfyi.StaticHandler)
	rtr.Post("/submit", tnfyi.FormHandler)
	return rtr
}

func main() {
	// logger
	logger := log.New(os.Stderr, "", log.Lshortfile)

	// client and url
	url := twentynine.TwentyNineApiUrl
	if envUrl, ok := os.LookupEnv("29_FYI_API_URL"); ok {
		url = envUrl
	}
	cli := twentynine.NewClient(twentynine.OptURL(url))

	// handler
	tnfyi := twentyninefyi{
		cli:    cli,
		logger: logger,
	}
	hndlr := tnfyi.Handler()

	// HTTP
	if err := http.ListenAndServe(":http", hndlr); err != nil {
		logger.Fatalf("HTTP server: %s", err)
	}
}
