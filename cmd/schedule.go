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

func calculateShiftStartTimes(times []int) []int {
	shiftStartTimes := []int{}

	for _, time := range times {
		if !slices.Contains(times, time-10) {
			shiftStartTimes = append(shiftStartTimes, time)
		}
	}
	return shiftStartTimes
}

func calculateShiftLength(startTime int, times []int) int {
	shiftLength := 10
	curTime := startTime
	for {
		curTime += 10
		if !slices.Contains(times, curTime) {
			break
		}
		shiftLength += 10
	}
	return shiftLength
}

func validateSchedule(selected map[string][]int) {
	totalTime := 0
	for day, times := range selected {
		dayTime := (len(times) * 10)

		shiftStartTimes := calculateShiftStartTimes(times)

		log.Print("Shift start times for ", day, ": ", shiftStartTimes)

		// Calculate shift lengths
		for _, shiftStart := range shiftStartTimes {
			length := calculateShiftLength(shiftStart, times)
			if length < 180 {
				log.Print("One shift must be at least 180 minutes long. A shift on ", day, " is: ", length)
			}
		}

		// Max daily work: 9 hours
		if dayTime > 540 {
			log.Print("Daily time must not exceed 540 minutes. Minutes for ", day, ": ", dayTime)
		}

		totalTime += dayTime
	}

	// Weekly total is 20-40 hours
	if totalTime < 1200 || totalTime > 2400 {
		log.Print("Total time per week must be between 1200 and 2400 minutes. Current minutes: ", totalTime)
	}
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
