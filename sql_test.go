package belajar_golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	//INSERT UPDATE DELETE BISA
	script := "INSERT INTO customer(id, name) VALUES ('fauzi', 'Achmad Sidik Fauzi')"
	_, err := db.ExecContext(ctx, script)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert New Customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "SELECT id, name FROM customer "
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("id", id)
		fmt.Println("name", name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name, email sql.NullString
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool
		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("====================")
		fmt.Println("Id", id)

		if name.Valid {
			fmt.Println("Name", name.String)
		}

		if email.Valid {
			fmt.Println("Email", email.String)
		}
		fmt.Println("Balance", balance)
		fmt.Println("Rating", rating)

		if birthDate.Valid {
			fmt.Println("Birth Date", birthDate.Time)
		}
		fmt.Println("Married", married)
		fmt.Println("Created At", createdAt)

	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	username := "admin"
	password := "admin"
	script := "SELECT username FROM users WHERE username = ? AND password = ? LIMIT 1 "
	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses Login", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestAutoIncrements(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "fauzi@gmail.com"
	comment := "fauzi ou"
	//INSERT UPDATE DELETE SAFE
	script := "INSERT INTO comment(email, comment) VALUES (?,?)"
	result, err := db.ExecContext(ctx, script, email, comment)

	if err != nil {
		panic(err)
	}
	insertId, err := result.LastInsertId()

	if err != nil {
		panic(err)
	}
	fmt.Println("Success Insert New Comment with Id ", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()
	script := "INSERT INTO comment(email, comment) VALUES(?,?)"
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "fauzi" + strconv.Itoa(i) + "@gmail.com"
		comment := "Ini komen ke " + strconv.Itoa(i)
		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id", id)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comment(email, comment) VALUES(?,?)"
	//do transaction
	for i := 0; i < 10; i++ {
		email := "fauzi" + strconv.Itoa(i) + "@gmail.com"
		comment := "Ini komen ke " + strconv.Itoa(i)
		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id", id)
	}

	//tx.Commit() execute
	//tx.Rollback akan gagal
	err = tx.Rollback() // tidak akan tersimpan ke database kecuali pake tx.Commit()
	if err != nil {
		panic(err)
	}
}
