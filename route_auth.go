// route_auth.go
package main

import (
	"fmt"
	"net/http"
	"workorder/data"
)

func login(writer http.ResponseWriter, request *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(writer, nil)
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create the user account
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	user := data.User{
		Name:     request.PostFormValue("name"),
		Password: request.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	http.Redirect(writer, request, "/login", 302)
}
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
func authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, _ := data.UserByName(r.PostFormValue("name"))
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		if r.PostFormValue("password") == "password" {
			http.Redirect(w, r, "/password", 302)
		}
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
func changePassword(writer http.ResponseWriter, request *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "password")
	t.Execute(writer, nil)
}
func updatePassword(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	user, err := getUserBySession(writer, request)
	if err != nil {
		return
	}
	opw, npw, rpw := request.PostFormValue("oldpassword"),
		request.PostFormValue("newpassword"), request.PostFormValue("replypassword")
	if user.Password == data.Encrypt(opw) && npw == rpw {
		user.Password = npw
		user.Update()
		http.Redirect(writer, request, "/", 302)
	} else {
		http.Redirect(writer, request, "/password", 302)
	}

}
func privilegeHandle(hf func(http.ResponseWriter, *http.Request, data.User),
	privilege int) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		curUser, err := getUserBySession(writer, request)
		if err != nil {
			http.Redirect(writer, request, "/login", 302)
			return
		}
		if curUser.Privileges&privilege != privilege {
			http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户没有该项操作权限！"), 302)
			return
		}
		hf(writer, request, curUser)
	}
}
