package main

import (
	"log"
	"net/http"
	"reflect"
	"text/template"
)

type StatusMessage struct {
	StatusType string
	Message    string
	Bullets    []string
}

func getHome(writer http.ResponseWriter, request *http.Request) {

	var data Week

	if reflect.ValueOf(userSchedule).IsZero() {
		data = initializeSchedule()
	} else {
		data = userSchedule.SubmittedWeek
	}

	files := []string{
		"./templates/base.html",
		"./templates/schedule.html",
		"./templates/week-view.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postSubmit(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse validations
	selected := map[string][]int{
		"Mon":   convertSliceStrToInt(request.Form["Mon"]),
		"Tues":  convertSliceStrToInt(request.Form["Tues"]),
		"Wed":   convertSliceStrToInt(request.Form["Wed"]),
		"Thurs": convertSliceStrToInt(request.Form["Thurs"]),
		"Fri":   convertSliceStrToInt(request.Form["Fri"]),
	}

	findings := []string{}

	// Schedule validation
	findings = append(findings, validateSchedule(selected)...)

	var data StatusMessage

	if userSchedule.Approved == true {
		data = StatusMessage{"failure", "Failure", []string{"Schedule has already been approved. Reset to submit a new one"}}
	} else if len(findings) > 0 {
		data = StatusMessage{"failure", "Failure", findings}
	} else {
		// Update/save week in server after validation
		updateWeek(selected)
		data = StatusMessage{"success", "Success", []string{"Schedule has been submitted for admin review"}}
	}

	template, err := template.ParseFiles("./templates/status-message.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "status-message", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	// TODO: Re-send schedule to frontend?
}
