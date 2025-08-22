package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func newCon() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:6306)/demo_12")
	if err != nil {
		log.Panic(err)
	}
	return db
}

func reset() error {
	db := newCon()

	txn, _ := db.Begin()
	defer db.Close()

	_, err := txn.Exec("UPDATE assign_seat SET user_id = NULL")

	if err != nil {
		txn.Rollback()
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Id       int    `"json:id"`
	UserName string `"json:user_name"`
}

type Seat struct {
	Id int `json:id`
}

func GetUser() []User {
	db := newCon()
	rows, err := db.Query("SELECT id, user_name from users")
	if err != nil {
		log.Panic(err)
	}

	defer rows.Close() // Ensure rows are closed
	defer db.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.UserName); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return users
}

func AssignTicket(user User, wg *sync.WaitGroup) {
	db := newCon()
	defer wg.Done()
	defer db.Close()
	txn, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	row := txn.QueryRow(`select id 
	from assign_seat
	where user_id is NULL
	order by id  limit 1`)
	if row.Err() != nil {
		txn.Rollback()
		log.Fatal(row.Err())
		//return nil, row.Err()
	}
	var seat Seat
	err = row.Scan(&seat.Id)

	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}

	_, err = txn.Exec("UPDATE assign_seat SET  user_id = ? where id = ?", user.Id, seat.Id)

	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	reset()
	startTime := time.Now()
	users := GetUser()

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go AssignTicket(user, &wg)
	}
	wg.Wait()
	fmt.Println("Time since ", time.Since(startTime))
}
