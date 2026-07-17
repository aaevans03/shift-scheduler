package main

import (
	"html/template"
	"log"
	"net/http"
)

func getHome(writer http.ResponseWriter, request *http.Request) {

	ts, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(writer, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", getHome)

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
