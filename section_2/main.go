package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type page struct {
	Name     string
	DBStatus bool
}

type searchResult struct {
	Title  string
	Author string
	Year   string
	ID     string
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "dev.db")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := page{Name: "Gopher"}

		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}

		p.DBStatus = db.Ping() == nil

		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		results := []searchResult{
			searchResult{"Moby-Dick", "Herman Melville", "1851", "22222"},
			searchResult{"The Adventures of Huckleberry Finn", "Mark Twain", "1884", "44444"},
			searchResult{"The Catcher in the Rye", "JD Salinger", "1951", "33333"},
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}
