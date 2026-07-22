package main

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"slices"
	"strconv"
)

type Block struct {
	Active  bool
	DayName string
	Time    string
}

type Hour struct {
	Blocks []Block // 6 blocks to an hour
}

type Day struct {
	DayName   string
	Hours     []Hour
	TotalTime float32
}

type Week struct {
	Days []Day
}

type Schedule struct {
	SubmittedWeek   Week
	Approved        bool
	ApprovalMessage string
}

var userSchedule Schedule
var dayNames = []string{"Mon", "Tues", "Wed", "Thurs", "Fri"}

func initializeSchedule() Week {

	// Loop through all days
	var week []Day
	for _, value := range dayNames {
		var hourList []Hour

		// Loop through all hours
		for hour := 800; hour < 1800; hour += 100 {
			var blockList []Block
			for minutes := 0; minutes <= 50; minutes += 10 {
				blockList = append(blockList, Block{false, value, strconv.Itoa(hour + minutes)})
			}
			hourList = append(hourList, Hour{blockList})
		}
		week = append(week, Day{value, hourList, 0})
	}

	userSchedule = Schedule{Week{week}, true, ""}
	return Week{week}
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

func updateWeek(selected map[string][]string) {

	for dayIndex := range userSchedule.SubmittedWeek.Days {
		day := &userSchedule.SubmittedWeek.Days[dayIndex]

		for hourIndex := range day.Hours {
			hour := &day.Hours[hourIndex]

			for blockIndex := range hour.Blocks {
				block := &hour.Blocks[blockIndex]

				block.Active = slices.Contains(selected[day.DayName], block.Time)
			}
		}
	}
	userSchedule.Approved = false
}

func postSubmit(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Parse validations
	selected := map[string][]string{
		"Mon":   request.Form["Mon"],
		"Tues":  request.Form["Tues"],
		"Wed":   request.Form["Wed"],
		"Thurs": request.Form["Thurs"],
		"Fri":   request.Form["Fri"],
	}

	updateWeek(selected)

	// TODO: Schedule validation
	// Save schedule in server
	// Re-send schedule to frontend?
	// Pop-up that says it was a success

}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", getHome)
	mux.HandleFunc("POST /submit", postSubmit)

	log.Print("starting server on http://localhost:4001")

	err := http.ListenAndServe(":4001", mux)
	log.Fatal(err)
}
