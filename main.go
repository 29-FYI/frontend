package main

import (
	"html/template"
	"net/http"

	"github.com/29-FYI/twentynine"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	links, err := twentynine.GetLinks()
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

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))).ServeHTTP(w, r)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/submit" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	headline := r.PostFormValue("headline")
	url := r.PostFormValue("url")

	if err := twentynine.PostLink(twentynine.Link{
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

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/submit", formHandler)
	http.ListenAndServe(":6970", nil)
}
