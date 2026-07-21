package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Block struct {
	Active bool
	Time   string
}

type Hour struct {
	Blocks []Block // 6 blocks to an hour
}

type Day struct {
	DayName string
	Hours   []Hour
}

type Week struct {
	Days []Day
}

func blankWeek() Week {
	dayNames := []string{"Mon", "Tues", "Wed", "Thurs", "Fri"}

	// Loop through all days
	var week []Day
	for _, value := range dayNames {
		var hourList []Hour

		// Loop through all hours
		for hour := 800; hour < 1800; hour += 100 {
			var blockList []Block
			for minutes := 0; minutes <= 50; minutes += 10 {
				blockList = append(blockList, Block{false, strconv.Itoa(hour + minutes)})
			}
			hourList = append(hourList, Hour{blockList})
		}
		week = append(week, Day{value, hourList})
	}

	return Week{week}
}

func getHome(writer http.ResponseWriter, request *http.Request) {
	data := blankWeek()

	files := []string{
		"./templates/base.html",
		"./templates/schedule.html",
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

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", getHome)

	log.Print("starting server on http://localhost:4001")

	err := http.ListenAndServe(":4001", mux)
	log.Fatal(err)
}
