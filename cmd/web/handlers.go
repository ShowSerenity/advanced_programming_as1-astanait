package main

import (
	"errors"
	"fmt"
	"net/http"
	"showserenity.net/astanait/pkg/forms"
	"showserenity.net/astanait/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.Path != "/" {
		app.notFound(w)
		return
	}*/

	s, err := app.newses.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{Newses: s})

}

func (app *application) showNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.newses.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.html", &templateData{News: s})
}

func (app *application) showAccountant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.accountants.GetAccountantById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "accountantshow.page.html", &templateData{Accountant: s})
}

func (app *application) pageAccountant(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.Path != "/" {
		app.notFound(w)
		return
	}*/

	s, err := app.accountants.LatestAccountants()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "accountantlist.page.html", &templateData{Accountants: s})

}

func (app *application) showStudentNews(w http.ResponseWriter, r *http.Request) {
	app.showNewsByType(w, r, "students")
}

func (app *application) showStaffNews(w http.ResponseWriter, r *http.Request) {
	app.showNewsByType(w, r, "staff")
}

func (app *application) showApplicantNews(w http.ResponseWriter, r *http.Request) {
	app.showNewsByType(w, r, "applicants")
}

func (app *application) showResearcherNews(w http.ResponseWriter, r *http.Request) {
	app.showNewsByType(w, r, "researchers")
}

func (app *application) showNewsByType(w http.ResponseWriter, r *http.Request, newsType string) {

	newsList, err := app.newses.GetByType(newsType)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "news.page.html", &templateData{NewsType: newsType, Newses: newsList})
}

func (app *application) createNewsForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}*/

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires", "newsType")
	form.PermittedValues("newsType", "students", "staff", "applicants", "researchers")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.html", &templateData{Form: form})
		return
	}

	id, err := app.newses.Insert(form.Get("title"), form.Get("content"), form.Get("expires"), form.Get("newsType"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "News created successfully!")

	http.Redirect(w, r, fmt.Sprintf("/news/%d", id), http.StatusSeeOther)
}

func (app *application) createAccountantForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "accountant.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createAccountant(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}*/

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "sname", "age")
	form.MaxLength("name", 100)
	form.MaxLength("sname", 100)

	if !form.Valid() {
		app.render(w, r, "accountant.page.html", &templateData{Form: form})
		return
	}

	id, err := app.newses.InsertAccountant(form.Get("name"), form.Get("sname"), form.Get("age"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/accountant/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.html", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.html", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.html", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/news/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
