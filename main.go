package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/glebarez/go-sqlite"
	"github.com/go-chi/chi"
)

func mainx() {
	fmt.Println(net.ParseIP("127.0.0.1"))
	fmt.Println(net.ParseIP("127.0.0.11/12"))
}

func main() {
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
