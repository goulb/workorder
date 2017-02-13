// workorder project main.go
package main

import (
	"net/http"
	"time"
	"workorder/data"
)

var depts []data.Department
var provs []data.Provider
var cartypes []data.CarType

var adminHandles map[string]func(http.ResponseWriter, *http.Request)

/*const adminHandles = map[string]func(http.ResponseWriter, *http.Request){
	"users":         users,
	"users/new":     newUser,
	"/users/create": createUser,
	"/users/edit":   editUser,
	"/users/update": updateUser,
}*/

func main() {
	p("WorkOrder", version(), "started at", config.Address)
	adminHandles = map[string]func(http.ResponseWriter, *http.Request){
		"/users":              users,
		"/users/new":          newUser,
		"/users/create":       createUser,
		"/users/edit":         editUser,
		"/users/update":       updateUser,
		"/departments":        departments,
		"/departments/new":    newDepartment,
		"/departments/create": createDepartment,
		"/providers":          providers,
		"/providers/new":      newProvider,
		"/providers/create":   createProvider,
	}
	mux := http.NewServeMux()
	depts, _ = data.Departments()
	provs, _ = data.Providers()
	cartypes, _ = data.CarTypes()

	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/err", err)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/orders", orders)
	mux.HandleFunc("/orders/new", newOrder)
	mux.HandleFunc("/orders/create", createOrder)
	mux.HandleFunc("/orders/update", updateOrder)
	mux.HandleFunc("/orders/edit", editOrder)
	mux.HandleFunc("/orders/delete", deleteOrder)
	mux.HandleFunc("/password", changePassword)
	mux.HandleFunc("/password/update", updatePassword)
	mux.HandleFunc("/password/reset", resetPassword)
	for key, _ := range adminHandles {
		mux.HandleFunc(key, admin)
	}

	mux.HandleFunc("/authenticate", authenticate)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
