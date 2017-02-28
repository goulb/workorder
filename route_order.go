package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"workorder/data"

	"encoding/json"

	"github.com/signintech/gopdf"
)

type displayWorkItem struct {
	WorkItem data.WorkItem
	UnitStr  string
}

var useforstr = []string{"生产用车", "检修用车"}
var unitstrs = []string{"台班", "吨"}

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
	} else {
		orders, _ = data.OrdersByDepartmentID(curUser.DepartmentId)
	}
	var dos []displayOrder

	for _, order := range orders {
		do := displayOrder{
			Id:         order.Id,
			Num:        fmt.Sprintf("GXJCWL%07d", 20000+order.Id),
			Department: deptMap()[order.DepartmentId],
			DateBegin:  order.DateBegin,
			DateEnd:    order.DateEnd,
			Provider:   providerMap()[order.ProviderId],
			CarType:    cartypeMap()[order.CarTypeId],
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

func newOrder(writer http.ResponseWriter, request *http.Request, user data.User) {

	info := struct {
		Depts     []data.Department
		Providers []data.Provider
		CarTypes  []data.CarType
		OrderNum  string
		CarType   string
		Order     data.Order
	}{
		depts,
		provs,
		cartypes,
		"",
		"",
		data.Order{},
	}
	generateHTML(writer, info, "login.layout", "public.navbar", "editorder")
}
func editOrder(writer http.ResponseWriter, request *http.Request, user data.User) {

	vals := request.URL.Query()
	orderid := 0
	fmt.Sscan(vals.Get("id"), &orderid)

	order, _ := data.OrderByID(orderid)
	if order.Locked {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前任务单已锁定！"), 302)
		return
	}
	workitems, _ := order.WorkItems()
	if len(workitems) > 0 {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前任务单已填写作业内容！"), 302)
		return
	}
	info := struct {
		Depts     []data.Department
		Providers []data.Provider
		CarTypes  []data.CarType
		OrderNum  string
		CarType   string
		Order     data.Order
	}{
		depts,
		provs,
		cartypes,
		fmt.Sprintf("GXJCWL%07d", 20000+order.Id),
		cartypeMap()[order.CarTypeId],
		order,
	}
	generateHTML(writer, info, "login.layout", "public.navbar", "editorder")
}
func deleteOrder(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	orderid := 0
	fmt.Sscan(vals.Get("id"), &orderid)

	order, _ := data.OrderByID(orderid)
	if order.Locked {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前任务单已锁定！"), 302)
		return
	}
	workitems, _ := order.WorkItems()
	if len(workitems) > 0 {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前任务单已填写作业内容！"), 302)
		return
	}
	order.Delete()
	http.Redirect(writer, request, "/orders", 302)
}

func updateOrder(writer http.ResponseWriter, request *http.Request, user data.User) {
	request.ParseForm()

	var id, deptid, providerid, usefor int
	fmt.Sscan(request.PostFormValue("id"), &id)
	fmt.Sscan(request.PostFormValue("department"), &deptid)
	fmt.Sscan(request.PostFormValue("provider"), &providerid)
	fmt.Sscan(request.PostFormValue("usefor"), &usefor)
	datebegin := request.PostFormValue("datebegin")
	dateend := request.PostFormValue("dateend")
	cartypeStr := request.PostFormValue("cartype")
	carnum := request.PostFormValue("carnum")

	order, _ := data.OrderByID(id)
	fmt.Println(order)
	order.DepartmentId = deptid
	order.DateBegin = datebegin
	order.DateEnd = dateend
	order.ProviderId = providerid
	order.CarTypeId = carTypeIDByString(cartypeStr)
	order.UseFor = usefor
	order.CarNum = carnum
	fmt.Println(order)
	order.Update()

	http.Redirect(writer, request, "/orders", 302)
}
func createOrder(writer http.ResponseWriter, request *http.Request, user data.User) {

	request.ParseForm()
	var deptid, providerid, usefor int
	fmt.Sscan(request.PostFormValue("department"), &deptid)
	fmt.Sscan(request.PostFormValue("provider"), &providerid)
	fmt.Sscan(request.PostFormValue("usefor"), &usefor)
	datebegin := request.PostFormValue("datebegin")
	dateend := request.PostFormValue("dateend")
	cartypeStr := request.PostFormValue("cartype")
	carnum := request.PostFormValue("carnum")

	order := data.Order{DepartmentId: deptid, DateBegin: datebegin,
		DateEnd: dateend, ProviderId: providerid,
		CarTypeId: carTypeIDByString(cartypeStr),
		UseFor:    usefor, CarNum: carnum}
	order.Create()

	http.Redirect(writer, request, "/orders", 302)
}

func cartypeMap() (cartypemap map[int]string) {
	cartypemap = map[int]string{}
	for _, cartype := range cartypes {
		s := ""
		if cartype.Weight == 0 {
			s = cartype.TypeName
		} else {
			s = fmt.Sprintf("%dT%s", cartype.Weight, cartype.TypeName)
		}
		cartypemap[cartype.Id] = s
	}
	return
}
func deptMap() (deptmap map[int]string) {
	deptmap = map[int]string{}
	for _, dept := range depts {
		deptmap[dept.Id] = dept.Name
	}
	return
}
func providerMap() (providermap map[int]string) {
	providermap = map[int]string{}
	for _, prov := range provs {
		providermap[prov.Id] = prov.Name
	}
	return
}
func carTypeIDByString(s string) (cartypeid int) {
	for _, ct := range cartypes {
		if fmt.Sprintf("%dT%s", ct.Weight, ct.TypeName) == s {
			cartypeid = ct.Id
		}
	}
	if cartypeid == 0 {
		var w int
		var t string

		fmt.Sscanf(s, "%dT%s", &w, &t)
		if t == "" {
			t = s
		}
		ct := data.CarType{Weight: w, TypeName: t}
		ct.Create()
		cartypeid = ct.Id
	}
	return
}
func unlockOrder(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	id := 0
	fmt.Sscan(vals.Get("id"), &id)
	order, err := data.OrderByID(id)
	if err != nil {
		danger(err)
		return
	}
	order.Locked = false
	order.Update()
	http.Redirect(writer, request, "/orders", 302)
}
func lockOrder(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	id := 0
	fmt.Sscan(vals.Get("id"), &id)
	order, err := data.OrderByID(id)
	if err != nil {
		danger(err)
		return
	}
	if user.DepartmentId != order.DepartmentId {
		return
	}
	order.Locked = true
	order.Update()
	http.Redirect(writer, request, "/orders", 302)
}
func workitems(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	id := 0
	fmt.Sscan(vals.Get("pid"), &id)
	order, err := data.OrderByID(id)

	temp := "readonly_workitems"
	if (user.DepartmentId == order.DepartmentId) && (!order.Locked) &&
		(user.Privileges&data.CanEdit == data.CanEdit) {
		temp = "workitems"
	}
	workItems, err := order.WorkItems()
	if err != nil {
		http.Redirect(writer, request,
			fmt.Sprintf("/err?msg=%s", fmt.Sprint(err)), 302)
	}
	dwis := []displayWorkItem{}
	for _, wi := range workItems {
		dwi := displayWorkItem{wi, unitstrs[wi.Unit]}
		dwis = append(dwis, dwi)
	}
	info := struct {
		OrderID   int
		WorkItems []displayWorkItem
	}{
		id,
		dwis,
	}
	generateHTML(writer, info, temp)
}
func deleteWorkitem(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	id := 0
	fmt.Sscan(vals.Get("id"), &id)
	workItem, err := data.WorkItemByID(id)
	if err != nil {
		danger(err)
	}
	workItem.Delete()
	http.Redirect(writer, request, "/orders", 302)
}
func updateWorkitem(writer http.ResponseWriter, request *http.Request, user data.User) {
	request.ParseForm()
	var id, unit int
	var quantity float32
	fmt.Sscan(request.PostFormValue("workitemid"), &id)
	fmt.Sscan(request.PostFormValue("unit"), &unit)
	fmt.Sscan(request.PostFormValue("quantity"), &quantity)

	workItem, err := data.WorkItemByID(id)
	if err != nil {
		danger(err)
	}
	workItem.Work = request.PostFormValue("work")
	workItem.Place = request.PostFormValue("place")
	workItem.Unit = unit
	workItem.Quantity = quantity
	workItem.Update()
	http.Redirect(writer, request, "/orders", 302)
}
func createWorkitem(writer http.ResponseWriter, request *http.Request, user data.User) {
	request.ParseForm()
	var orderid, unit int
	var quantity float32
	fmt.Sscan(request.PostFormValue("orderid"), &orderid)
	fmt.Sscan(request.PostFormValue("unit"), &unit)
	fmt.Sscan(request.PostFormValue("quantity"), &quantity)
	work := request.PostFormValue("work")
	place := request.PostFormValue("place")
	order, err := data.OrderByID(orderid)
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}
	if user.DepartmentId != order.DepartmentId {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户无操作权限！"), 302)
		return
	}
	order.CreateWorkItem(work, place, unit, quantity)
	http.Redirect(writer, request, "/orders", 302)
}
func newWorkitem(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	pid := 0
	fmt.Sscan(vals.Get("pid"), &pid)
	order, err := data.OrderByID(pid)
	if err != nil {
		danger(err)
	}
	if user.DepartmentId != order.DepartmentId {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户无操作权限！"), 302)
		return
	}
	info := struct {
		Order      data.Order
		OrderNum   string
		Department string
		Provider   string
		Cartype    string
		WorkItem   data.WorkItem
	}{
		order,
		fmt.Sprintf("GXJCWL%07d", 20000+order.Id),
		deptMap()[order.DepartmentId],
		providerMap()[order.ProviderId],
		cartypeMap()[order.CarTypeId],
		data.WorkItem{},
	}
	generateHTML(writer, info, "login.layout", "public.navbar", "editworkitem")
}

func editWorkitem(writer http.ResponseWriter, request *http.Request, user data.User) {
	vals := request.URL.Query()
	id := 0
	fmt.Sscan(vals.Get("id"), &id)
	p(id)
	workItem, err := data.WorkItemByID(id)
	if err != nil {
		danger(err)
	}
	order, err := data.OrderByID(workItem.OrderId)
	if err != nil {
		danger(err)
	}
	if user.DepartmentId != order.DepartmentId {
		http.Redirect(writer, request, fmt.Sprintf("/err?msg=%s", "当前用户无操作权限！"), 302)
		return
	}
	info := struct {
		Order      data.Order
		OrderNum   string
		Department string
		Provider   string
		Cartype    string
		WorkItem   data.WorkItem
	}{
		order,
		fmt.Sprintf("GXJCWL%07d", 20000+order.Id),
		deptMap()[order.DepartmentId],
		providerMap()[order.ProviderId],
		cartypeMap()[order.CarTypeId],
		workItem,
	}
	generateHTML(writer, info, "login.layout", "public.navbar", "editworkitem")
}
func printOrders(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	sids := vals.Get("ids")
	ids := strings.Split(sids[2:len(sids)-2], ",")
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 648, H: 360}}) //595.28, 841.89 = A4
	err := pdf.AddTTFFont("kai", "fireflysung.ttf")
	if err != nil {
		danger(err)
	}
	file, err := os.Open("template.json")
	if err != nil {
		danger(err)
	}
	decoder := json.NewDecoder(file)
	template := PdfTemplate{}
	err = decoder.Decode(&template)
	if err != nil {
		danger(err)
	}
	for _, sid := range ids {
		id := 0
		fmt.Sscan(sid, &id)
		order, _ := data.OrderByID(id)
		if order.Locked {
			wis, _ := order.WorkItems()
			pdf.AddPage()

			err = pdf.SetFont("kai", "", 14)
			if err != nil {
				danger(err)
			}

			for _, cell := range template.Cells {
				DrawCellText(&pdf, cell)
			}
			if order.UseFor == 0 {
				template.Titles[0].Text = "安全生产部确认签章\n年    月    日"
			} else {
				template.Titles[0].Text = "设备管理中心确认签章\n年    月    日"
			}
			template.Titles[1].Text = "用车单位（签章）：" + deptMap()[order.DepartmentId]
			template.Titles[2].Text = order.DateBegin + " - " + order.DateEnd
			template.Titles[3].Text = "编号：" + fmt.Sprintf("GXJCWL%07d", 20000+order.Id)
			template.Titles[4].Text = providerMap()[order.ProviderId]
			template.Titles[5].Text = cartypeMap()[order.CarTypeId]
			template.Titles[6].Text = order.CarNum
			for _, cell := range template.Titles {
				DrawCellText(&pdf, cell)
			}
			detailcells := []CellText{}
			for i, wi := range wis {
				cell1 := template.Details[0]
				cell1.Text = wi.Work
				cell1.Top = cell1.Top + float64(i)*25
				detailcells = append(detailcells, cell1)
				cell2 := template.Details[1]
				cell2.Text = wi.Place
				cell2.Top = cell2.Top + float64(i)*25
				detailcells = append(detailcells, cell2)
				cell3 := template.Details[2]
				cell3.Text = fmt.Sprintf("%.2f%s", wi.Quantity, unitstrs[wi.Unit])
				cell3.Top = cell3.Top + float64(i)*25
				detailcells = append(detailcells, cell3)
			}
			for _, cell := range detailcells {
				DrawCellText(&pdf, cell)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/pdf")
	writer.Write(pdf.GetBytesPdf())

}
