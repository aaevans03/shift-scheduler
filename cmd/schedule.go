package main

import (
	"log"
	"slices"
	"strconv"
)

type Block struct {
	Active  bool
	DayName string
	Time    int
}

type Day struct {
	DayName   string
	Hours     []Block
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
		var hourList []Block

		// Loop through all minutes in a workday
		for minutes := 480; minutes < 1080; minutes += 10 {
			hourList = append(hourList, Block{false, value, minutes})
		}
		week = append(week, Day{value, hourList, 0})
	}

	userSchedule = Schedule{Week{week}, true, ""}
	return Week{week}
}

func convertSliceStrToInt(stringSlice []string) []int {

	var convertedDay []int

	for _, time := range stringSlice {
		convertedTime, err := strconv.Atoi(time)
		if err != nil {
			log.Print("error")
		}
		convertedDay = append(convertedDay, convertedTime)
	}
	return convertedDay
}

func validateSchedule(selected map[string][]int) {
	// TODO: implement this
	// Weekly total is 20-40 hours (sum up all blocks)
	// Max daily work: 9 hours
	// Min daily shift length: 3 hours
}

func updateWeek(selected map[string][]int) {
	// TODO: If schedule is already approved, you can't edit it, unless you reset it

	for dayIndex := range userSchedule.SubmittedWeek.Days {
		day := &userSchedule.SubmittedWeek.Days[dayIndex]

		for hourIndex := range day.Hours {
			hour := &day.Hours[hourIndex]

			hour.Active = slices.Contains(selected[day.DayName], hour.Time)
		}
	}
	userSchedule.Approved = false
}
