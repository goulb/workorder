package data

var depts = []Department{
	{
		Name: "冶炼厂",
	},
	{
		Name: "物流中心",
	},
}

var providers = []Provider{
	{
		Name: "海纳物流",
	},
	{
		Name: "中盛物流",
	},
}

var users = []User{
	{
		Name:         "Tom",
		Password:     "123",
		DepartmentID: 1,
	},
	{
		Name: "Mary",
	},
}

func setup() {

	PrviderDeleteAll()
	SessionDeleteAll()
	UserDeleteAll()
	DepartmentDeleteAll()
}
