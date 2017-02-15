// handler.go
package main

import (
	"fmt"
	//"html/template"
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
	if request.URL.Path != "/" {
		error_message(writer, request, fmt.Sprintf("页面%s未找到！", request.URL.Path))
	} else {
		if _, err := getUserBySession(writer, request); err == nil {
			http.Redirect(writer, request, "/orders", 302)
		}
	}
}
