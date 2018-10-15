package thereg

import (
	"database/sql"
	"fmt"
	"log"
)

func connect() *sql.DB {
	connStr := "postgres://localhost/thereg?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// DBGetAccount gets a single Account
func DBGetAccount(id string) Account {
	db := connect()
	row := db.QueryRow("SELECT id, email, email_confirm_token, email_confirmed, username, ip, authtoken, created_at, updated_at FROM accounts WHERE id = $1", id)
	al := AccountList{}
	return al.FromDBRow(row)
}

// DBGetAccountByToken gets a single Account
func DBGetAccountByToken(authToken string) Account {
	db := connect()
	row := db.QueryRow("SELECT id, email, email_confirm_token, email_confirmed, username, ip, authtoken, created_at, updated_at FROM accounts WHERE authtoken = $1", authToken)
	al := AccountList{}
	return al.FromDBRow(row)
}

// DBInsertAccount inserts an Account into the DB
func DBInsertAccount(a Account) Account {
	db := connect()
	row := db.QueryRow(`INSERT INTO accounts
		(id, email, email_confirm_token, email_confirmed, username, ip, authtoken, created_at, updated_at)
		VALUES
		(%s, %s, %s, %s, %s, %s, %s, NOW(), NOW())`,
		a.ID, a.Email, a.EmailConfirmToken, a.EmailConfirmed, a.Username, a.IP, a.Authtoken)
	al := AccountList{}
	return al.FromDBRow(row)
}

// DBGetNode gets a single Node
func DBGetNode(id string) Node {
	db := connect()
	row := db.QueryRow("SELECT id, account_id, name, url, hostname, port, status, public, created_at, updated_at FROM nodes WHERE id = $1 OR name = $1", id)
	nl := NodeList{}
	return nl.FromDBRow(row)
}

// DBGetNodes gets a NodeList
func DBGetNodes(accountID string) NodeList {
	db := connect()
	rows, err := db.Query("SELECT id, account_id, name, url, hostname, port, status, public, created_at, updated_at FROM nodes WHERE account_id = $1", accountID)
	if err != nil {
		fmt.Println("Error loading Nodes:", err)
		return NodeList{}
	}
	nl := NodeList{}
	return nl.FromDBRows(rows)
}

// DBInsertNode inserts a Node into the DB
func DBInsertNode(n Node) Node {
	db := connect()
	row := db.QueryRow(`INSERT INTO nodes
		(id, account_id, name, url, hostname, port, status, public, created_at, updated_at)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`,
		n.ID, n.AccountID, n.Name, n.URL, n.Hostname, n.Port, n.Status, n.Public)
	nl := NodeList{}
	return nl.FromDBRow(row)
}
