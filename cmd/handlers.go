package main

import (
	"log"
	"net/http"
	"reflect"
	"text/template"
)

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

	// Schedule validation
	validateSchedule(selected)

	// Update/save week in server after validation
	updateWeek(selected)

	// TODO:
	// Re-send schedule to frontend?
	// Pop-up that says it was a success

}
