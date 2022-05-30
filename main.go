package main

import (
	"compress/gzip"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/go-chi/chi"
)

//go:embed web/public/*
var myfs embed.FS

func main() {
	sfs, err := fs.Sub(myfs, "web/public")
	if err != nil {
		log.Fatal(err)
	}

	hfs := http.FileServer(http.FS(sfs))
	http.Handle("/", gzipped(hfs))

	log.Print("Listening on :3000...")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func gzipped(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, req)
			return
		}

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		for k, v := range rec.Header() {
			if k == "Content-Length" {
				continue
			}
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		if rec.Result().StatusCode == http.StatusOK {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(rec.Result().StatusCode)
			gw := gzip.NewWriter(w)
			gw.Write(rec.Body.Bytes())
			gw.Flush()
		} else {
			w.WriteHeader(rec.Result().StatusCode)
			w.Write(rec.Body.Bytes())
		}
		// fmt.Println(rec.Body.Len())
		// w.Write(rec.Body.Bytes())

		// gzip.NewWriter(w).Write(rec.Body.Bytes())

	})
}

func mainx() {
	fmt.Println(net.ParseIP("127.0.0.1"))
	fmt.Println(net.ParseIP("127.0.0.11/12"))
}

func maindb() {
	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM certificates")
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for rows.Next() {
		i++
	}
	rows.Close()

	fmt.Println(i, "rows found")
}

func main2() {
	// connect
	db, err := sql.Open("sqlite", "./sscm.sqlite")
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

var createTableQuery = `CREATE TABLE IF NOT EXISTS certificates (
  id INT PRIMARY KEY,
  parent INT,
  cert BLOB,
  key BLOB,
  info TEXT
);
CREATE INDEX IF NOT EXISTS cert_parent ON certificates (parent);
`

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./sscm.sqlite")
	if err != nil {
		return db, err
	}

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return db, err
	}

	return db, err
}
