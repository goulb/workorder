package main

import (
	"fmt"
	"net/http"
	//"time"
	"workorder/data"
)

var useforstr = []string{"生产用车", "检修用车"}

type displayOrder struct {
	Id         int
	Num        string
	Department string
	DateBegin  string
	DateEnd    string
	Provider   string
	CarType    string
	CarNum     string
	UseFor     string
	Submit     bool
	Locked     bool
	CreatedAt  string
}

func orders(writer http.ResponseWriter, request *http.Request) {
	curUser, err := getUserBySession(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	var orders []data.Order

	if curUser.Privileges&data.CanEditAll == data.CanEditAll ||
		curUser.Privileges&data.CanBroweAll == data.CanBroweAll {
		orders, _ = data.Orders()
		fmt.Println(0, orders)
	} else {
		orders, _ = data.OrdersByDepartmentID(curUser.DepartmentId)
		fmt.Println(1, orders)
	}
	var dos []displayOrder

	deptmap := map[int]string{}
	for _, dept := range depts {
		deptmap[dept.Id] = dept.Name
	}
	providermap := map[int]string{}
	for _, prov := range provs {
		providermap[prov.Id] = prov.Name
	}
	cartypemap := map[int]string{}
	for _, cartype := range cartypes {
		s := ""
		if cartype.Weight == 0 {
			s = cartype.TypeName
		} else {
			s = fmt.Sprintf("%dT%s", cartype.Weight, cartype.TypeName)
		}
		cartypemap[cartype.Id] = s
	}
	for _, order := range orders {
		do := displayOrder{
			Id:         order.Id,
			Num:        fmt.Sprintf("GXJCWL%07d", 20000+order.Id),
			Department: deptmap[order.DepartmentId],
			DateBegin:  order.DateBegin,
			DateEnd:    order.DateEnd,
			Provider:   providermap[order.ProviderId],
			CarType:    cartypemap[order.CarTypeId],
			CarNum:     order.CarNum,
			CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
			UseFor:     useforstr[order.UseFor],
			Submit:     order.Submit,
			Locked:     order.Locked,
		}
		dos = append(dos, do)
	}
	navbar := "private.navbar"
	if curUser.Privileges&data.CanAdmin == data.CanAdmin {
		navbar = "admin.navbar"
	}
	info := struct {
		CanEdit       bool
		DisplayOrders []displayOrder
	}{
		curUser.Privileges&data.CanEditAll == data.CanEditAll,
		dos,
	}
	generateHTML(writer, info, "layout", navbar, "orders")
}

func newOrder(writer http.ResponseWriter, request *http.Request) {
	curUser, err := getUserBySession(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	if curUser.Privileges&data.CanEditAll != data.CanEditAll {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户没有添加任务单权限！"), 302)
		return
	}
	info := struct {
		Depts     []data.Department
		Providers []data.Provider
		CarTypes  []data.CarType
	}{
		depts,
		provs,
		cartypes,
	}
	generateHTML(writer, info, "login.layout", "public.navbar", "neworder")
}
func editOrder(writer http.ResponseWriter, request *http.Request) {
	curUser, err := getUserBySession(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	if curUser.Privileges&data.CanEditAll != data.CanEditAll {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户没有添加任务单权限！"), 302)
		return
	}
	vals := request.URL.Query()
	orderid := 0
	fmt.Sscan(vals.Get("id"), &orderid)

	order, _ := data.OrderByID(orderid)
	if order.Locked {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前任务单已锁定！"), 302)
		return
	}
	info := struct {
		Depts     []data.Department
		Providers []data.Provider
		CarTypes  []data.CarType
		OrderNum  string
		Order     data.Order
	}{
		depts,
		provs,
		cartypes,
		fmt.Sprintf("GXJCWL%07d", 20000+order.Id),
		order,
	}
	generateHTML(writer, info, "login.layout", "public.navbar", "editorder")
}

func deleteOrder(writer http.ResponseWriter, request *http.Request) {
	curUser, err := getUserBySession(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	if curUser.Privileges&data.CanEditAll != data.CanEditAll {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户没有添加任务单权限！"), 302)
		return
	}

	http.Redirect(writer, request, "/orders", 302)
}
func updateOrder(writer http.ResponseWriter, request *http.Request) {

	dept := data.Department{Name: request.PostFormValue("name")}
	err := dept.Create()
	if err != nil {
		fmt.Println(err)
	}
	depts, err = data.Departments()
	http.Redirect(writer, request, "/orders", 302)
}
func createOrder(writer http.ResponseWriter, request *http.Request) {
	curUser, err := getUserBySession(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
		return
	}
	if curUser.Privileges&data.CanEditAll != data.CanEditAll {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户没有添加任务单权限！"), 302)
		return
	}
	request.ParseForm()
	var deptid, providerid, usefor int
	fmt.Sscan(request.PostFormValue("department"), &deptid)
	fmt.Sscan(request.PostFormValue("provider"), &providerid)
	fmt.Sscan(request.PostFormValue("usefor"), &usefor)
	datebegin := request.PostFormValue("datebegin")
	dateend := request.PostFormValue("dateend")
	cartypeStr := request.PostFormValue("cartype")
	carnum := request.PostFormValue("carnum")
	var cartypeid int
	for _, ct := range cartypes {
		if fmt.Sprintf("%dT%s", ct.Weight, ct.TypeName) == cartypeStr {
			cartypeid = ct.Id
		}
	}
	if cartypeid == 0 {
		var w int
		var t string
		fmt.Sscanf(cartypeStr, "%dT%s", &w, &t)
		ct := data.CarType{Weight: w, TypeName: t}
		ct.Create()
		cartypeid = ct.Id
	}

	order := data.Order{DepartmentId: deptid, DateBegin: datebegin,
		DateEnd: dateend, ProviderId: providerid, CarTypeId: cartypeid,
		UseFor: usefor, CarNum: carnum}
	order.Create()

	http.Redirect(writer, request, "/orders", 302)
}
