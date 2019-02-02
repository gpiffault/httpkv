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
	_, err = db.Exec("create table if not exists data (key BLOB PRIMARY KEY, value BLOB)")
	checkErr(err)

	handler := &KVHandler{db: db}
	http.Handle("/", handler)

	err = http.ListenAndServe(":8765", nil)
	log.Fatal(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type KVHandler struct {
	db *sql.DB
}

func (handler KVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.Get(w, r)
	case http.MethodPut:
		handler.Put(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler KVHandler) Get(w http.ResponseWriter, r *http.Request) {
	rows, err := handler.db.Query("select value from data where key = ?", r.URL.Path)
	checkErr(err)
	defer rows.Close()
	if !rows.Next() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var value string
	rows.Scan(&value)
	fmt.Fprintf(w, value)
}

func (handler KVHandler) Put(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	_, err = handler.db.Exec("insert or replace into data values (?,?)", r.URL.Path, body)
	checkErr(err)
}
