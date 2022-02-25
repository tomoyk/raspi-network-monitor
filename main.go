package main

import "database/sql"
import "log"
import "time"
import _ "github.com/mattn/go-sqlite3"

func main() {
	db, err := sql.Open("sqlite3", "./metrics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table metrics (timestamp INTEGER NOT NULL PRIMARY KEY, value REAL);
	delete from metrics;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	stmt, err := db.Prepare("insert into metrics(timestamp, value) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	dt := time.Now()
	unix := dt.Unix()
	stmt.Exec(unix, 10)
}
