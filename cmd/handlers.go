package main

import (
	"log"
	"net/http"
	"text/template"
)

type WeekViewData struct {
	Week     Week
	EditMode bool
}

type StatusMessage struct {
	StatusType string
	Message    string
	Bullets    []string
}

func getHome(writer http.ResponseWriter, request *http.Request) {

	data := WeekViewData{
		Week:     getScheduleFromMemory(),
		EditMode: false,
	}

	files := []string{
		"./templates/base.html",
		"./templates/view-schedule.html",
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

func getSchedule(writer http.ResponseWriter, request *http.Request) {

	data := WeekViewData{
		Week:     getScheduleFromMemory(),
		EditMode: false,
	}

	files := []string{
		"./templates/view-schedule.html",
		"./templates/week-view.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "view-schedule", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getScheduleEdit(writer http.ResponseWriter, request *http.Request) {

	data := WeekViewData{
		Week:     getScheduleFromMemory(),
		EditMode: true,
	}

	files := []string{
		"./templates/edit-schedule.html",
		"./templates/week-view.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "edit-schedule", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

}

func postScheduleSubmit(writer http.ResponseWriter, request *http.Request) {
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
