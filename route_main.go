// handler.go
package main

import (
	"html/template"
	"net/http"
	"workorder/data"
)

func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}
func getUserBySession(writer http.ResponseWriter, request *http.Request) (user data.User, err error) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		user, err = sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
			http.Redirect(writer, request, "/login", 302)
		}
	}
	return
}

func index(writer http.ResponseWriter, request *http.Request) {
	//fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])

	getUserBySession(writer, request)
	files := []string{"templates/layout.html",
		"templates/navbar.html",
		"templates/index.html"}
	templates := template.Must(template.ParseFiles(files...))
	//worlorders := []workorder{{"12345", "xx"}, {"67890", "yy"}}
	workorders := data.GetWorkorder()
	//threads, err := data.Threads()
	//if err == nil {

	templates.ExecuteTemplate(writer, "layout", workorders)
	//}
}
