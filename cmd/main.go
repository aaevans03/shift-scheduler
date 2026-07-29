package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("POST /login/admin", postLoginAdmin)
	mux.HandleFunc("POST /login/student", postLoginStudent)

	mux.HandleFunc("GET /{$}", getHome)
	mux.HandleFunc("GET /schedule", getSchedule)
	mux.HandleFunc("GET /schedule/edit", getScheduleEdit)
	mux.HandleFunc("POST /schedule/submit", postScheduleSubmit)
	// TODO: implement DELETE /submit endpoint

	log.Print("starting server on http://localhost:4001")

	err := http.ListenAndServe(":4001", mux)
	log.Fatal(err)
}
