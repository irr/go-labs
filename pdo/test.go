package main

// http://blog.prevoty.com/simple-tools-to-solve-simple-problems

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prevoty/pdo"
	"log"
)

var (
	sqlitedb    *pdo.Sqlite
	Sqlitebuilt = false
)

type User struct {
	_meta string `table:"user"`

	Id    int    `column:"id"`
	First string `column:"first_name"`
	Last  string `column:"last_name"`
}

func EnsureSqlite() {

	if Sqlitebuilt {
		return
	}

	database, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec(`DROP TABLE IF EXISTS "user"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(`CREATE TABLE "user" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,
	"first_name" VARCHAR, "last_name" VARCHAR)`)
	if err != nil {
		log.Fatal(err)
	}

	database.Exec(`INSERT INTO "user" ("first_name","last_name") 
	VALUES ("Alessandra","Santos"),("Ivan","Rocha")`)
	database.Close()

	sqlitedb, err = pdo.NewSqlite("test.db")
	if err != nil {
		log.Fatal(err)
	}

	Sqlitebuilt = true
}

func main() {

	EnsureSqlite()

	db, err := pdo.NewSqlite("test.db")
	if err != nil {
		log.Fatal(err)
	}

	id, err := db.Create(&User{
		First: "Artemis",
		Last:  "Prime",
	})

	if err != nil {
		log.Fatal(err)
	}

	user := new(User)

	err = db.Find(user, "WHERE `id` = ?", id)
	switch err {
	case sql.ErrNoRows:
		log.Println("no record found...")
	case nil:
		log.Println("found a user")
	default:
		log.Fatal(err)
	}

	users := make([]*User, 0, 0)

	err = db.FindAll(&users, "WHERE `id` < 10")
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		log.Printf("%#v\n", u)
	}
}
