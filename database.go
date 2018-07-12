package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

//Initialsing database
func initDB(state *RuntimeState) (err error) {

	state.dbType = "sqlite3"
	state.db, err = sql.Open("sqlite3", state.Config.Base.StorageURL)
	if err != nil {
		return err
	}
	if true {
		sqlStmt := `create table if not exists users_data (id INTEGER PRIMARY KEY AUTOINCREMENT, useremail text not null,username text not null, password text not null, time_stamp int not null);`
		_, err = state.db.Exec(sqlStmt)
		if err != nil {
			log.Printf("init sqlite3 err: %s: %q\n", err, sqlStmt)
			return err
		}
	}

	return nil
}

//insert a Userinfo into DB
func insertUserInDB(useremail string, username string, password string, state *RuntimeState) error {

	stmtText := "insert into users_data(useremail,username, password, time_stamp) values (?,?,?,?);"
	stmt, err := state.db.Prepare(stmtText)
	if err != nil {
		log.Print("Error Preparing statement")
		log.Fatal(err)
	}
	defer stmt.Close()
	if useremailExistsorNot(useremail, state) {
		log.Println("UserEmail already exists")
		return errors.New("UserEmail already exists")
	}
	if usernameExistsorNot(username, state) {
		log.Println("UserName already exists")
		return errors.New("UserName already exists")
	}
	_, err = stmt.Exec(useremail, username, password, time.Now().Unix())

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//deleting the username info in DB.
func deleteEntryInDB(username string, state *RuntimeState) error {

	stmtText := "delete from users_data where username= ?;"
	stmt, err := state.db.Prepare(stmtText)
	if err != nil {
		log.Print("Error Preparing statement")
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

//look for password of a username
func findpasswdofUserinDB(username string, state *RuntimeState) (string, bool, error) {
	stmtText := "select password from users_data where username=?;"
	stmt, err := state.db.Prepare(stmtText)
	if err != nil {
		log.Print("Error Preparing statement")
		log.Fatal(err)
		return "", false, err
	}
	defer stmt.Close()
	var password string
	rows, err := stmt.Query(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("err='%s'", err)
			return "", false, nil
		} else {
			log.Printf("Problem with db ='%s'", err)
			return "", false, err
		}
	}
	defer rows.Close()
	//there should only be one entry with a username.
	if rows.Next() {
		err = rows.Scan(&password)
	}

	return password, true, nil

}

//check if user email already exists or not
func useremailExistsorNot(useremail string, state *RuntimeState) bool {
	stmtText := "select * from users_data where useremail=?"
	stmt, err := state.db.Prepare(stmtText)
	if err != nil {
		log.Print("Error Preparing statement")
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(useremail)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("err='%s'", err)
			return false
		} else {
			log.Printf("Problem with db ='%s'", err)
			return false
		}
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

//check if the username already exists or not
func usernameExistsorNot(username string, state *RuntimeState) bool {
	stmtText := "select * from users_data where useremail=?"
	stmt, err := state.db.Prepare(stmtText)
	if err != nil {
		log.Print("Error Preparing statement")
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("err='%s'", err)
			return false
		} else {
			log.Printf("Problem with db ='%s'", err)
			return false
		}
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}
