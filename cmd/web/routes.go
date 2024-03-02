package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice" // New import
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/article/:id", dynamicMiddleware.ThenFunc(app.showArticle))
	mux.Get("/article/create/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createArticleForm))
	mux.Post("/article/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createArticle))

	mux.Get("/students", dynamicMiddleware.ThenFunc(app.students))
	mux.Get("/teachers", dynamicMiddleware.ThenFunc(app.teachers))

	//SIGN UP AND LOGIN AND LOGOUT
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	mux.Get("/admin", dynamicMiddleware.ThenFunc(app.admin))

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)

}
