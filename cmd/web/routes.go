package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/news", app.showNews)
	mux.HandleFunc("/news/students", app.showStudentNews)
	mux.HandleFunc("/news/staff", app.showStaffNews)
	mux.HandleFunc("/news/applicants", app.showApplicantNews)
	mux.HandleFunc("/news/researchers", app.showResearcherNews)
	mux.HandleFunc("/news/create", app.createNews)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
