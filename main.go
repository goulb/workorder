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
	mux.HandleFunc("/orders/new", changeOrder(newOrder))
	mux.HandleFunc("/orders/create", changeOrder(createOrder))
	mux.HandleFunc("/orders/update", changeOrder(updateOrder))
	mux.HandleFunc("/orders/edit", changeOrder(editOrder))
	mux.HandleFunc("/orders/delete", changeOrder(deleteOrder))

	mux.HandleFunc("/providers", doAdmin(providers))
	mux.HandleFunc("/providers/new", doAdmin(editProvider))
	mux.HandleFunc("/providers/edit", doAdmin(editProvider))
	mux.HandleFunc("/providers/create", doAdmin(createProvider))
	mux.HandleFunc("/providers/update", doAdmin(updateProvider))

	mux.HandleFunc("/password", changePassword)
	mux.HandleFunc("/password/update", updatePassword)
	mux.HandleFunc("/password/reset", resetPassword)
	mux.HandleFunc("/workitems", workitems)
	mux.HandleFunc("/workitems/new", newWorkitem)
	mux.HandleFunc("/workitems/create", createWorkitem)
	mux.HandleFunc("/workitems/update", updateWorkitem)
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
