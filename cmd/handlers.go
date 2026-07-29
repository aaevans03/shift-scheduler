package main

import (
	"log"
	"net/http"
	"text/template"
)

type ScheduleFrontendData struct {
	Week            Week
	EditMode        bool
	ApprovalStatus  string
	ApprovalMessage string

	CurrentUser UserRole
	IsAdmin     bool
}

type AdminFrontendData struct {
	EditMode        bool   // needed for HTML template
	ApprovalStatus  string // needed for HTML template
	CurrentUser     UserRole
	IsAdmin         bool
	StudentSchedule Schedule
}

type StatusMessage struct {
	StatusType string
	Message    string
	Bullets    []string
}

func postLoginAdmin(writer http.ResponseWriter, request *http.Request) {
	clearSession(writer, request)
	createSession(writer, AdminUser)

	// Admin starter screen (you can select)

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	files := []string{
		"./templates/admin-dashboard.html",
		"./templates/status-message.html",
		"./templates/header-status.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "admin-dashboard", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "status-message-clear-oob", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postLoginStudent(writer http.ResponseWriter, request *http.Request) {
	clearSession(writer, request)
	createSession(writer, StudentUser)

	// Default screen

	memorySchedule := getScheduleFromMemory()

	data := ScheduleFrontendData{
		Week:            memorySchedule.SubmittedWeek,
		EditMode:        false,
		ApprovalStatus:  memorySchedule.ApprovalStatus,
		ApprovalMessage: memorySchedule.ApprovalMessage,

		CurrentUser: StudentUser,
		IsAdmin:     false,
	}

	files := []string{
		"./templates/view-schedule.html",
		"./templates/week-view.html",
		"./templates/status-message.html",
		"./templates/header-status.html",
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

	err = template.ExecuteTemplate(writer, "status-message-clear-oob", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

}

func getHome(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		user = StudentUser
		createSession(writer, user)
	}

	memorySchedule := getScheduleFromMemory()

	var data any
	if user == AdminUser {
		data = AdminFrontendData{
			EditMode:        false,
			ApprovalStatus:  "",
			CurrentUser:     AdminUser,
			IsAdmin:         true,
			StudentSchedule: getScheduleFromMemory(),
		}

	} else {
		data = ScheduleFrontendData{
			Week:            memorySchedule.SubmittedWeek,
			EditMode:        false,
			ApprovalStatus:  memorySchedule.ApprovalStatus,
			ApprovalMessage: memorySchedule.ApprovalMessage,

			CurrentUser: user,
			IsAdmin:     user == AdminUser,
		}
	}

	files := []string{
		"./templates/base.html",
		"./templates/header-status.html",
		"./templates/view-schedule.html",
		"./templates/week-view.html",
		"./templates/admin-dashboard.html",
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
	user, ok := currentUser(request)
	if !ok {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memorySchedule := getScheduleFromMemory()

	data := ScheduleFrontendData{
		Week:            memorySchedule.SubmittedWeek,
		EditMode:        false,
		ApprovalStatus:  memorySchedule.ApprovalStatus,
		ApprovalMessage: memorySchedule.ApprovalMessage,

		CurrentUser: user,
		IsAdmin:     user == AdminUser,
	}

	files := []string{
		"./templates/view-schedule.html",
		"./templates/week-view.html",
		"./templates/status-message.html",
		"./templates/header-status.html",
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

	err = template.ExecuteTemplate(writer, "status-message-clear-oob", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getScheduleEdit(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memorySchedule := getScheduleFromMemory()

	data := ScheduleFrontendData{
		Week:            memorySchedule.SubmittedWeek,
		EditMode:        true,
		ApprovalStatus:  memorySchedule.ApprovalStatus,
		ApprovalMessage: memorySchedule.ApprovalMessage,

		CurrentUser: user,
		IsAdmin:     user == AdminUser,
	}

	files := []string{
		"./templates/edit-schedule.html",
		"./templates/week-view.html",
		"./templates/status-message.html",
		"./templates/header-status.html",
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

	err = template.ExecuteTemplate(writer, "status-message-clear-oob", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postScheduleSubmit(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	statusTemplate, err := template.ParseFiles("./templates/status-message.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if userSchedule.ApprovalStatus == "Approved" {
		data := StatusMessage{
			"failure",
			"Failure",
			[]string{"Schedule has already been approved. Reset to submit a new one"},
		}

		writer.WriteHeader(http.StatusConflict)
		err = statusTemplate.ExecuteTemplate(writer, "status-message", data)
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	if len(findings) > 0 {
		data := StatusMessage{
			"failure",
			"Failure",
			findings,
		}

		writer.WriteHeader(http.StatusUnprocessableEntity)
		err = statusTemplate.ExecuteTemplate(writer, "status-message", data)
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	// Update/save week in server after validation
	updateWeek(selected)

	memorySchedule := getScheduleFromMemory()

	weekData := ScheduleFrontendData{
		Week:            memorySchedule.SubmittedWeek,
		EditMode:        false,
		ApprovalStatus:  memorySchedule.ApprovalStatus,
		ApprovalMessage: memorySchedule.ApprovalMessage,

		CurrentUser: user,
		IsAdmin:     user == AdminUser,
	}

	statusData := StatusMessage{
		"success",
		"Success",
		[]string{"Schedule has been submitted for admin review"},
	}

	viewScheduleTemplate, err := template.ParseFiles(
		"./templates/view-schedule.html",
		"./templates/week-view.html",
		"./templates/status-message.html",
		"./templates/header-status.html",
	)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = viewScheduleTemplate.ExecuteTemplate(writer, "view-schedule", weekData)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = viewScheduleTemplate.ExecuteTemplate(writer, "status-message-oob", statusData)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = viewScheduleTemplate.ExecuteTemplate(writer, "header-status-oob", weekData)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getAdmin(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		user = StudentUser
		createSession(writer, user)
	} else if user != AdminUser {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	files := []string{
		"./templates/admin-dashboard.html",
		"./templates/header-status.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "admin-dashboard", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getAdminView(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	} else if user != AdminUser {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memorySchedule := getScheduleFromMemory()

	data := ScheduleFrontendData{
		Week:            memorySchedule.SubmittedWeek,
		EditMode:        false,
		ApprovalStatus:  memorySchedule.ApprovalStatus,
		ApprovalMessage: memorySchedule.ApprovalMessage,

		CurrentUser: user,
		IsAdmin:     user == AdminUser,
	}

	files := []string{
		"./templates/view-schedule.html",
		"./templates/week-view.html",
		"./templates/status-message.html",
		"./templates/header-status.html",
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

	err = template.ExecuteTemplate(writer, "status-message-clear-oob", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func putAdminApprove(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		user = StudentUser
		createSession(writer, user)
	} else if user != AdminUser {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Invalid form data", http.StatusBadRequest)
		return
	}

	approveSchedule(request.Form["approval-comments"][0])

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	files := []string{
		"./templates/admin-dashboard.html",
		"./templates/header-status.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "admin-dashboard", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func putAdminReject(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		user = StudentUser
		createSession(writer, user)
	} else if user != AdminUser {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Invalid form data", http.StatusBadRequest)
		return
	}

	rejectSchedule(request.Form["approval-comments"][0])

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	files := []string{
		"./templates/admin-dashboard.html",
		"./templates/header-status.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = template.ExecuteTemplate(writer, "admin-dashboard", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(writer, "header-status-oob", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
