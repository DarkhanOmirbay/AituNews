package main

import (
	"darkhanomirbay.net/aitunews/pkg/forms"
	"darkhanomirbay.net/aitunews/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	a, err := app.articles.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Articles: a,
	})

}
func (app *application) admin(w http.ResponseWriter, r *http.Request) {

	u, err := app.users.GetAllUsers()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "admin.page.tmpl", &templateData{
		Users: u,
	})

}
func (app *application) students(w http.ResponseWriter, r *http.Request) {
	a, err := app.articles.GetStudents()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "students.page.tmpl", &templateData{
		Articles: a,
	})
}
func (app *application) teachers(w http.ResponseWriter, r *http.Request) {
	a, err := app.articles.GetTeachers()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "teachers.page.tmpl", &templateData{
		Articles: a,
	})
}

func (app *application) showArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)

		return
	}
	a, err := app.articles.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "show.page.tmpl", &templateData{
		Article: a,
		Flash:   flash,
	})

}

func (app *application) createArticleForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createArticle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//title := r.PostForm.Get("title")
	//content := r.PostForm.Get("content")
	//category := r.PostForm.Get("category")
	form := forms.New(r.PostForm)
	form.Required("title", "content", "category")
	form.MaxLength("title", 100)
	form.PermittedValues("category", "For students", "For teachers")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.articles.Insert(form.Get("title"), form.Get("content"), form.Get("category"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/article/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
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
	form.Required("fullname", "email", "password", "role")
	form.MaxLength("fullname", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	form.MaxLength("role", 25)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form})
		return
	}

	err = app.users.Insert(form.Get("fullname"), form.Get("email"), form.Get("password"), form.Get("role"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
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
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been logged out succesfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
