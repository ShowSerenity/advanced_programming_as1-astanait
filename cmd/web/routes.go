package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/news/students", dynamicMiddleware.ThenFunc(app.showStudentNews))
	mux.Get("/news/staff", dynamicMiddleware.ThenFunc(app.showStaffNews))
	mux.Get("/news/applicants", dynamicMiddleware.ThenFunc(app.showApplicantNews))
	mux.Get("/news/researchers", dynamicMiddleware.ThenFunc(app.showResearcherNews))
	mux.Get("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNewsForm))
	mux.Post("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNews))
	mux.Get("/news/:id", dynamicMiddleware.ThenFunc(app.showNews))

	mux.Get("/accountant/list", dynamicMiddleware.ThenFunc(app.pageAccountant))
	mux.Get("/accountant/create", dynamicMiddleware.ThenFunc(app.createAccountantForm))
	mux.Post("/accountant/create", dynamicMiddleware.ThenFunc(app.createAccountant))
	mux.Get("/accountant/:id", dynamicMiddleware.ThenFunc(app.showAccountant))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
