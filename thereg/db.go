package thereg

import (
	"database/sql"
	"fmt"
	"log"
)

func connect() *sql.DB {
	connStr := "dbname=thereg sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// DBGetAccount gets a single Account
func DBGetAccount(id string) Account {
	db := connect()
	rows, err := db.Query("SELECT id, email, email_confirm_token, email_confirmed, username, ip, authtoken FROM accounts WHERE id = $1", id)
	fmt.Println(rows)
	if err != nil {
		// nothing
	}
	// TODO make this return Account
	return rows
}

// DBGetNode gets a single Node
func DBGetNode(id string) Node {
	db := connect()
	rows, err := db.Query("SELECT id, name, url, hostname, port, status, public, created_at, updated_at FROM nodes WHERE id = $1", id)
	fmt.Println(rows)
	if err != nil {
		// nothing
	}
	// TODO make this return Node
	return rows
}

// DBGetNodes gets a NodeList
func DBGetNodes(id string) NodeList {
	db := connect()
	rows, err := db.Query("SELECT id, name, url, hostname, port, status, public, created_at, updated_at FROM nodes WHERE account_id = $1", id)
	fmt.Println(rows)
	if err != nil {
		// nothing
	}
	// TODO make this return NodeList
	return rows
}
