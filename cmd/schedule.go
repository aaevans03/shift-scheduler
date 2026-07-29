package main

import (
	"fmt"
	"log"
	"math"
	"reflect"
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
	TotalTime float64
}

type Week struct {
	Days      []Day
	TotalTime float64
}

type Schedule struct {
	SubmittedWeek   Week
	ApprovalStatus  string
	ApprovalMessage string
}

var userSchedule Schedule
var dayNames = []string{"Mon", "Tues", "Wed", "Thurs", "Fri"}

func initializeSchedule() Schedule {

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

	userSchedule = Schedule{Week{week, 0.0}, "Not Submitted", ""}
	return userSchedule
}

func getScheduleFromMemory() Schedule {
	var data Schedule

	if reflect.ValueOf(userSchedule).IsZero() {
		data = initializeSchedule()
	} else {
		data = userSchedule
	}
	return data
}

func convertSliceStrToInt(stringSlice []string) []int {

	var convertedDay []int

	for _, time := range stringSlice {
		convertedTime, err := strconv.Atoi(time)
		if err != nil {
			log.Print("Error converting string array to int array")
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

func minutesToHours(minutes int) float64 {
	return math.Round((float64(minutes)/60.0)*100.0) / 100
}

func validateSchedule(selected map[string][]int) []string {
	findings := []string{}

	totalTime := 0
	for day, times := range selected {
		dayTime := (len(times) * 10)

		shiftStartTimes := calculateShiftStartTimes(times)

		// Calculate shift lengths
		for _, shiftStart := range shiftStartTimes {
			length := calculateShiftLength(shiftStart, times)
			if length < 180 {
				findings = append(
					findings,
					fmt.Sprintf("<strong>%s</strong> has a shift of %.2f hours. Shifts must be at least 3 hours.", day, minutesToHours(length)),
				)
			}
		}

		// Max daily work: 9 hours
		if dayTime > 540 {
			findings = append(
				findings,
				fmt.Sprintf("<strong>%s</strong> has a total work time of length %.2f hours. Do not exceed 9 hours in one day.", day, minutesToHours(dayTime)),
			)
		}

		totalTime += dayTime
	}

	// Weekly total is 20-40 hours
	if totalTime < 1200 || totalTime > 2400 {
		findings = append(
			findings,
			fmt.Sprintf("Total work time of the week is %.2f hours. It must be 20-40 hours.", minutesToHours(totalTime)),
		)
	}

	return findings
}

func updateWeek(selected map[string][]int) {
	weekTotalBlocks := 0
	for dayIndex := range userSchedule.SubmittedWeek.Days {
		day := &userSchedule.SubmittedWeek.Days[dayIndex]

		dayTotalBlocks := 0
		for hourIndex := range day.Hours {
			hour := &day.Hours[hourIndex]

			if slices.Contains(selected[day.DayName], hour.Time) {
				hour.Active = true
				dayTotalBlocks++
				weekTotalBlocks++
			} else {
				hour.Active = false
			}
		}
		day.TotalTime = minutesToHours(dayTotalBlocks * 10)
	}
	userSchedule.SubmittedWeek.TotalTime = minutesToHours(weekTotalBlocks * 10)
	userSchedule.ApprovalStatus = "Pending Approval"
}
