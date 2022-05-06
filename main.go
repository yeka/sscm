package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
	"github.com/go-chi/chi"
)

func main() {
	b := []byte{}
	bb := bytes.NewBuffer(b)
	bb.Write([]byte("Hello"))
	fmt.Println(b)
}

func main1() {
	// connect
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	// get SQLite version
	row := db.QueryRow("select sqlite_version()")
	var s string
	row.Scan(&s)
	fmt.Println(s)
	_ = chi.NewMux()
}
