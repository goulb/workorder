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
	mux.HandleFunc("/orders/new", privilegeHandle(newOrder, data.CanEditAll))
	mux.HandleFunc("/orders/create", privilegeHandle(createOrder, data.CanEditAll))
	mux.HandleFunc("/orders/update", privilegeHandle(updateOrder, data.CanEditAll))
	mux.HandleFunc("/orders/edit", privilegeHandle(editOrder, data.CanEditAll))
	mux.HandleFunc("/orders/delete", privilegeHandle(deleteOrder, data.CanEditAll))
	mux.HandleFunc("/workitems", privilegeHandle(workitems, 0))
	mux.HandleFunc("/workitems/new", privilegeHandle(newWorkitem, data.CanEdit))
	mux.HandleFunc("/workitems/edit", privilegeHandle(editWorkitem, data.CanEdit))
	mux.HandleFunc("/workitems/create", privilegeHandle(createWorkitem, data.CanEdit))
	mux.HandleFunc("/workitems/update", privilegeHandle(updateWorkitem, data.CanEdit))

	mux.HandleFunc("/providers", privilegeHandle(providers, data.CanAdmin))
	mux.HandleFunc("/providers/new", privilegeHandle(editProvider, data.CanAdmin))
	mux.HandleFunc("/providers/edit", privilegeHandle(editProvider, data.CanAdmin))
	mux.HandleFunc("/providers/create", privilegeHandle(createProvider, data.CanAdmin))
	mux.HandleFunc("/providers/update", privilegeHandle(updateProvider, data.CanAdmin))
	mux.HandleFunc("/users", privilegeHandle(users, data.CanAdmin))
	mux.HandleFunc("/users/new", privilegeHandle(newUser, data.CanAdmin))
	mux.HandleFunc("/users/create", privilegeHandle(createUser, data.CanAdmin))
	mux.HandleFunc("/users/edit", privilegeHandle(editUser, data.CanAdmin))
	mux.HandleFunc("/users/update", privilegeHandle(updateUser, data.CanAdmin))
	mux.HandleFunc("/departments", privilegeHandle(departments, data.CanAdmin))
	mux.HandleFunc("/departments/new", privilegeHandle(newDepartment, data.CanAdmin))
	mux.HandleFunc("/departments/create", privilegeHandle(createDepartment, data.CanAdmin))

	mux.HandleFunc("/password", changePassword)
	mux.HandleFunc("/password/update", updatePassword)
	mux.HandleFunc("/password/reset", privilegeHandle(resetPassword, data.CanAdmin))

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
