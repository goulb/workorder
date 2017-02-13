package data

type Provider struct {
	Id   int
	Name string
}

func (prov *Provider) Create() (err error) {

	_, err = Db.Exec("insert into providers values($1)", prov.Name)
	if err != nil {
		return
	}

	err = Db.QueryRow(`select rowid,name from providers where 
		rowid in (select last_insert_rowid())`).Scan(&prov.Id, &prov.Name)
	if err != nil {
		return
	}

	return
}
func ProviderByID(id int) (prov Provider, err error) {
	err = Db.QueryRow("select rowid,name from providers where rowid=$1", id).Scan(&prov.Id, &prov.Name)
	return
}
func (prov *Provider) Delete() (err error) {
	_, err = Db.Exec("delete from providers where rowid=$1", prov.Id)
	return
}
func (prov *Provider) Update() (err error) {
	_, err = Db.Exec("update providers set name=$1 where rowid=$2", prov.Name, prov.Id)
	return
}
func Providers() (provs []Provider, err error) {
	rows, err := Db.Query("select rowid,name from providers order by rowid")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		prov := Provider{}
		if err = rows.Scan(&prov.Id, &prov.Name); err != nil {
			return
		}
		provs = append(provs, prov)
	}
	return
}
func PrviderDeleteAll() (err error) {
	_, err = Db.Exec("delete from providers")
	return
}
