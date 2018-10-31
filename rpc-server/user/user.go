package user

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var (
	saveStmt, editStmt, getStmt *sql.Stmt
	db                          *sql.DB
	err                         error
)

type User struct {
	Uuid  string `json:"uuid"`
	Login string `json:"login"`
	Date  string `json:"date"`
}

// used for testing
type Response struct {
	Result string `json:"result"`
}
type Message struct {
	Message string
}

// create all stmts that User will need
func init() {
	// for now presume we read config file and get all params form there
	db, err = makeConnection("postgres", "postgres://user:pass@localhost/users")
	if err != nil {
		log.Fatal(err)
	}
	saveStmt, err = db.Prepare("INSERT INTO users (uuid, login, date) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	// edit just one field
	editStmt, err = db.Prepare("UPDATE users SET login = ? WHERE uuid = ?")
	if err != nil {
		log.Fatal(err)
	}
	getStmt, err = db.Prepare("SELECT * FROM users WHERE uuid = ?")
	if err != nil {
		log.Fatal(err)
	}
}

func makeConnection(driver string, params string) (db *sql.DB, err error) {
	db, err = sql.Open(driver, params)
	return
}

func makeUuid(login string) string {
	data := sha1.Sum([]byte(fmt.Sprintf("%s%s", login, time.Now())))
	t := fmt.Sprintf("%x", data)
	return fmt.Sprintf("%s-%s-%s-%s-%s", t[:8], t[8:12], t[12:16], t[16:20], t[20:32])
}

func makeDate() string {
	t := time.Now()
	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
}

func (u *User) Test(msg Message, r *Response) error {
	*r = Response{Result: msg.Message}
	return nil
}

func (u *User) Save(args User, rows *int64) (err error) {
	login := args.Login
	if len(login) == 0 {
		return fmt.Errorf("Can not read required argumet login from argumets")
	}
	id := makeUuid(login)
	date := makeDate()
	res, err := saveStmt.Exec(id, login, date)
	*rows, err = res.RowsAffected()
	return
}

func (u *User) Edit(args User, rows *int64) (err error) {
	login := args.Login
	if len(login) == 0 {
		return fmt.Errorf("Can not read required argumet Login from argumets")
	}
	id := args.Uuid
	res, err := editStmt.Exec(login, id)
	*rows, err = res.RowsAffected()
	return
}

func (u *User) Get(args User, user *User) (err error) {
	id := args.Uuid
	res, err := getStmt.Query(id)
	if err != nil {
		return err
	}
	err = res.Scan(user.Uuid, user.Login, user.Date)
	return
}
