package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "data.db")
	checkErr(err)
	defer db.Close()
	stmt, err := db.Prepare("create table if not exists data (key BLOB PRIMARY KEY, value BLOB)")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	stmt.Close()

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rows, err := db.Query("select value from data where key = ?", r.URL.Path)
			checkErr(err)
			defer rows.Close()
			if !rows.Next() {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			var value string
			rows.Scan(&value)
			fmt.Fprintf(w, value)
		case http.MethodPut:
			body, err := ioutil.ReadAll(r.Body)
			checkErr(err)
			stmt, err = db.Prepare("insert into data values (?,?)")
			checkErr(err)
			defer stmt.Close()
			_, err = stmt.Exec(r.URL.Path, body)
			checkErr(err)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	})

	err = http.ListenAndServe(":8765", nil)
	log.Fatal(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
