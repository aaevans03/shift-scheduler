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

	renderTemplates(
		writer,
		[]string{
			"./templates/admin-dashboard.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"admin-dashboard", data},
		TemplateRender{"status-message-clear-oob", nil},
		TemplateRender{"header-status-oob", data},
	)
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

	renderTemplates(
		writer,
		[]string{
			"./templates/view-schedule.html",
			"./templates/week-view.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"view-schedule", data},
		TemplateRender{"status-message-clear-oob", nil},
		TemplateRender{"header-status-oob", data},
	)
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

	renderTemplates(
		writer,
		[]string{
			"./templates/base.html",
			"./templates/header-status.html",
			"./templates/view-schedule.html",
			"./templates/week-view.html",
			"./templates/admin-dashboard.html",
		},
		TemplateRender{"base", data},
	)
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

	renderTemplates(
		writer,
		[]string{
			"./templates/view-schedule.html",
			"./templates/week-view.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"view-schedule", data},
		TemplateRender{"status-message-clear-oob", nil},
		TemplateRender{"header-status-oob", data},
	)
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

	renderTemplates(
		writer,
		[]string{
			"./templates/edit-schedule.html",
			"./templates/week-view.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"edit-schedule", data},
		TemplateRender{"status-message-clear-oob", nil},
		TemplateRender{"header-status-oob", data},
	)
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

	renderTemplates(
		writer,
		[]string{
			"./templates/view-schedule.html",
			"./templates/week-view.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"view-schedule", weekData},
		TemplateRender{"status-message-clear-oob", statusData},
		TemplateRender{"header-status-oob", weekData},
	)
}

func deleteSchedule(writer http.ResponseWriter, request *http.Request) {
	user, ok := currentUser(request)
	if !ok {
		log.Print("401 Unauthorized Access Attempted")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newSchedule := initializeSchedule()

	data := ScheduleFrontendData{
		Week:            newSchedule.SubmittedWeek,
		EditMode:        false,
		ApprovalStatus:  newSchedule.ApprovalStatus,
		ApprovalMessage: newSchedule.ApprovalMessage,

		CurrentUser: user,
		IsAdmin:     user == AdminUser,
	}

	renderTemplates(
		writer,
		[]string{
			"./templates/view-schedule.html",
			"./templates/week-view.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"view-schedule", data},
		TemplateRender{"status-message-clear-oob", nil},
		TemplateRender{"header-status-oob", data},
	)
}

func getAdmin(writer http.ResponseWriter, request *http.Request) {
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

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	renderTemplates(
		writer,
		[]string{
			"./templates/admin-dashboard.html",
			"./templates/header-status.html",
		},
		TemplateRender{"admin-dashboard", data},
		TemplateRender{"header-status-oob", data},
	)
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

	renderTemplates(
		writer,
		[]string{
			"./templates/view-schedule.html",
			"./templates/week-view.html",
			"./templates/status-message.html",
			"./templates/header-status.html",
		},
		TemplateRender{"view-schedule", data},
		TemplateRender{"status-message-clear-oob", nil},
		TemplateRender{"header-status-oob", data},
	)
}

func putAdminApprove(writer http.ResponseWriter, request *http.Request) {
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

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Invalid form data", http.StatusBadRequest)
		return
	}

	approveSchedule(request.FormValue("approval-comments"))

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	renderTemplates(
		writer,
		[]string{
			"./templates/admin-dashboard.html",
			"./templates/header-status.html",
		},
		TemplateRender{"admin-dashboard", data},
		TemplateRender{"header-status-oob", data},
	)
}

func putAdminReject(writer http.ResponseWriter, request *http.Request) {
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

	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Invalid form data", http.StatusBadRequest)
		return
	}

	rejectSchedule(request.FormValue("approval-comments"))

	data := AdminFrontendData{
		EditMode:        false,
		ApprovalStatus:  "",
		CurrentUser:     AdminUser,
		IsAdmin:         true,
		StudentSchedule: getScheduleFromMemory(),
	}

	renderTemplates(
		writer,
		[]string{
			"./templates/admin-dashboard.html",
			"./templates/header-status.html",
		},
		TemplateRender{"admin-dashboard", data},
		TemplateRender{"header-status-oob", data},
	)
}
