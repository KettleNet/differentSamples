package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

func main() {
	// remove database file for clean experiment
	os.Remove("./sqlite/example.db")
	// open database file
	db, err := sql.Open("sqlite3", "./sqlite/example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// create script TODO: add file.sql
	sqlStmt := `CREATE TABLE data (
  		id integer PRIMARY KEY AUTOINCREMENT,
  		did integer,
  		date datetime,
  		data varchar,
  		data_type varchar
	);

	CREATE TABLE device (
	  id integer PRIMARY KEY AUTOINCREMENT,
	  device_name varchar,
	  device_condition varchar
	);
	`
	// exec script
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// prepare executing
	stmt, err := tx.Prepare("insert into device(id, device_name, device_condition) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("deviceName" + strconv.Itoa(i)), fmt.Sprintf("OK"))
		if err != nil {
			log.Fatal(err)
		}
	}
	// commit transaction
	tx.Commit()

	// query example
	rows, err := db.Query("select id, device_name from device")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var device_name string
		err = rows.Scan(&id, &device_name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, device_name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// prepare query
	stmt, err = db.Prepare("select device_name from device where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var device_name string
	err = stmt.QueryRow("3").Scan(&device_name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(device_name)

	// delete all devices
	_, err = db.Exec("delete from device")
	if err != nil {
		log.Fatal(err)
	}

	// simple execution
	_, err = db.Exec("insert into device(id, device_name, device_condition) values(1, 'light', 'OK'), (2, 'temp', 'OVERHEAT'), (3, 'pressure', 'DEAD')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, device_name, device_condition from device")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var device_name, device_condition string
		err = rows.Scan(&id, &device_name, &device_condition)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, device_name, device_condition)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
