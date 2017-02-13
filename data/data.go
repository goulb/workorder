// data.go
package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type workorder struct {
	Num        string
	Department string
}

const timeformat = "2006-01-02 15:04:05.000000-0700"

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	return
}
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}
func GetWorkorder() []workorder {
	return []workorder{{"12345", "xx"}, {"67890", "yy"}}
}
