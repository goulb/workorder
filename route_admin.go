package main

import (
	"fmt"
	"net/http"
	"workorder/data"
)

type displayUser struct {
	Id             int
	Name           string
	Department     string
	Privileges     []string
	PasswordStatus string
}

var privilegeStrings = map[int]string{
	data.CanEdit:     "维护内容",
	data.CanBroweAll: "查看单据",
	data.CanEditAll:  "维护单据",
	data.CanAdmin:    "系统管理",
}

func admin(writer http.ResponseWriter, request *http.Request) {
	/*curUser, err := getUserBySession(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	if curUser.Privileges&data.CanAdmin != data.CanAdmin {
		http.Redirect(writer, request, "/", 302)
		return
	}*/
	//fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
	adminHandles[request.URL.Path](writer, request)
}

func users(writer http.ResponseWriter, request *http.Request) {

	users, err := data.Users()
	if err != nil {
		fmt.Println(err)
		return
	}

	deptmap := map[int]string{}
	for _, dept := range depts {
		deptmap[dept.Id] = dept.Name
	}

	displayUsers := []displayUser{}
	for _, user := range users {
		du := displayUser{
			Id:             user.Id,
			Name:           user.Name,
			PasswordStatus: "已更改",
		}
		for k, v := range privilegeStrings {
			if user.Privileges&k == k {
				du.Privileges = append(du.Privileges, v)
			}
		}
		du.Department = deptmap[user.DepartmentId]
		if user.Password == data.Encrypt("password") {

			du.PasswordStatus = "初始密码：password"
		}

		displayUsers = append(displayUsers, du)
	}
	generateHTML(writer, displayUsers, "layout", "admin.navbar", "users")
}
func departments(writer http.ResponseWriter, request *http.Request) {

	generateHTML(writer, depts, "layout", "admin.navbar", "departments")
}
func newDepartment(writer http.ResponseWriter, request *http.Request) {

	t := parseTemplateFiles("login.layout", "public.navbar", "newdepartment")
	t.Execute(writer, nil)
}
func createDepartment(writer http.ResponseWriter, request *http.Request) {

	dept := data.Department{Name: request.PostFormValue("name")}
	err := dept.Create()
	if err != nil {
		danger(err)
	}
	depts, err = data.Departments()
	http.Redirect(writer, request, "/departments", 302)
}
func providers(writer http.ResponseWriter, request *http.Request) {
	provs, err := data.Providers()
	if err != nil {
		danger(err)
		return
	}
	generateHTML(writer, provs, "layout", "admin.navbar", "providers")
}
func editProvider(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	provider := data.Provider{}
	id := 0
	fmt.Sscanln(vals.Get("id"), &id)
	provider, _ = data.ProviderByID(id)

	generateHTML(writer, provider, "login.layout", "public.navbar",
		"editprovider")
}
func updateProvider(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	id := 0
	fmt.Sscan(request.PostFormValue("id"), &id)
	provider, err := data.ProviderByID(id)
	if err != nil {
		error_message(writer, request, fmt.Sprintln(err))
		return
	}
	provider.Name = request.PostFormValue("name")
	provider.FullName = request.PostFormValue("fullname")
	err = provider.Update()
	if err != nil {
		error_message(writer, request, fmt.Sprintln(err))
		return
	}
	http.Redirect(writer, request, "/providers", 302)
}
func createProvider(writer http.ResponseWriter, request *http.Request) {
	provider := data.Provider{
		Name:     request.PostFormValue("name"),
		FullName: request.PostFormValue("fullname"),
	}
	err := provider.Create()
	if err != nil {
		error_message(writer, request, fmt.Sprintln(err))
		return
	}
	http.Redirect(writer, request, "/providers", 302)
}
func createUser(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	var deptId, canEdit, canBroweAll, canEditAll, canAdmin int
	fmt.Sscan(request.PostFormValue("department"), &deptId)
	fmt.Sscan(request.PostFormValue("canEdit"), &canEdit)
	fmt.Sscan(request.PostFormValue("canBroweAll"), &canBroweAll)
	fmt.Sscan(request.PostFormValue("canEditAll"), &canEditAll)
	fmt.Sscan(request.PostFormValue("canAdmin"), &canAdmin)
	p(deptId, canEdit, canBroweAll, canEditAll, canAdmin)
	user := data.User{Name: request.PostFormValue("name"), DepartmentId: deptId,
		Password: "password", Privileges: canEdit | canBroweAll | canEditAll | canAdmin}
	err := user.Create()
	if err != nil {
		error_message(writer, request, fmt.Sprintln(err))
		return
	}
	http.Redirect(writer, request, "/users", 302)
}
func updateUser(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	var id, deptId, canEdit, canBroweAll, canEditAll, canAdmin int
	fmt.Sscan(request.PostFormValue("id"), &id)
	fmt.Sscan(request.PostFormValue("department"), &deptId)
	fmt.Sscan(request.PostFormValue("canEdit"), &canEdit)
	fmt.Sscan(request.PostFormValue("canBroweAll"), &canBroweAll)
	fmt.Sscan(request.PostFormValue("canEditAll"), &canEditAll)
	fmt.Sscan(request.PostFormValue("canAdmin"), &canAdmin)

	user, err := data.UserByID(id)
	if err != nil {
		error_message(writer, request, fmt.Sprintln(err))
		return
	}
	user.Name = request.PostFormValue("name")
	user.DepartmentId = deptId
	user.Privileges = canEdit | canBroweAll | canEditAll | canAdmin
	user.Update()
	http.Redirect(writer, request, "/users", 302)
}
func resetPassword(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("reset password")
	vals := request.URL.Query()
	var id int
	fmt.Sscan(vals.Get("id"), &id)

	user, err := data.UserByID(id)
	if err != nil {
		error_message(writer, request, "Cannot read thread")
	} else {
		user.Password = "password"
		user.Update()
		http.Redirect(writer, request, "/users", 302)
	}
}
func newUser(writer http.ResponseWriter, request *http.Request) {
	items := struct {
		User        data.User
		Departments []data.Department
	}{
		User:        data.User{},
		Departments: depts,
	}
	generateHTML(writer, items, "login.layout", "public.navbar", "edituser")
}
func editUser(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	var id int
	fmt.Sscan(vals.Get("id"), &id)

	user, err := data.UserByID(id)
	if err != nil {
		error_message(writer, request, "Cannot read thread")
	} else {
		items := struct {
			User        data.User
			Departments []data.Department
		}{
			User:        user,
			Departments: depts,
		}
		generateHTML(writer, items, "login.layout", "public.navbar", "edituser")
	}
}
func doAdmin(hf func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		curUser, err := getUserBySession(writer, request)
		if err != nil {
			warning(err)
			http.Redirect(writer, request, "/login", 302)
			return
		}
		if curUser.Privileges&data.CanAdmin != data.CanAdmin {
			http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户权限！"), 302)
			return
		}
		hf(writer, request)
	}
}
