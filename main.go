package main

import (
	"html/template"
	"log"
	"net/http"
)

func getHome(writer http.ResponseWriter, request *http.Request) {

	files := []string{
		"./templates/base.html",
		"./templates/schedule.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(writer, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", getHome)

	log.Print("starting server on http://localhost:4001")

	err := http.ListenAndServe(":4001", mux)
	log.Fatal(err)
}
