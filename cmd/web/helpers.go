package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	// Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	// Add the flash message to the template data, if one exists.
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	td.IsAdmin = app.isAdmin(r)
	td.IsTeacher = app.isTeacher(r)

	return td
}
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	// Execute the template set, passing the dynamic data with the current
	// year injected.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}
func (app *application) isAdmin(r *http.Request) bool {
	// Получите ID пользователя из сессии
	userID := app.session.GetInt(r, "authenticatedUserID")
	if userID == 0 {
		// Если ID пользователя не существует, верните false
		return false
	}

	// Используйте метод вашей модели для проверки роли пользователя
	role, err := app.users.IsAdmin(userID)
	if err != nil {
		// Обработайте ошибку по вашему усмотрению, например, логирование или возврат false
		app.errorLog.Println(err)
		return false
	}

	// Проверьте, является ли роль "Admin"
	return role == "Admin"
}

func (app *application) isTeacher(r *http.Request) bool {
	userID := app.session.GetInt(r, "authenticatedUserID")
	if userID == 0 {
		return false
	}

	role, err := app.users.IsTeacher(userID)
	if err != nil {
		app.errorLog.Println(err)
		return false
	}
	return role == "Teacher"
}
