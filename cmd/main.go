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
	mux.HandleFunc("DELETE /schedule", deleteSchedule)

	mux.HandleFunc("GET /admin", getAdmin)
	mux.HandleFunc("GET /admin/view", getAdminView)
	mux.HandleFunc("PUT /admin/approve", putAdminApprove)
	mux.HandleFunc("PUT /admin/reject", putAdminReject)

	log.Print("starting server on http://localhost:4001")

	err := http.ListenAndServe(":4001", mux)
	log.Fatal(err)
}
