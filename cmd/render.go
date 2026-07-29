package main

import (
	"log"
	"net/http"
	"text/template"
)

type TemplateRender struct {
	name string
	data any
}

func renderTemplates(writer http.ResponseWriter, files []string, renders ...TemplateRender) bool {
	template, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	for _, render := range renders {
		if err := template.ExecuteTemplate(writer, render.name, render.data); err != nil {
			log.Print(err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return false
		}
	}
	return true
}
