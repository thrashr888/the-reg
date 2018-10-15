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

// DBGetAccountByConfirmToken gets a single Account
func DBGetAccountByConfirmToken(confirmToken string) Account {
	db := connect()
	row := db.QueryRow("SELECT id, email, email_confirm_token, email_confirmed, username, ip, authtoken, created_at, updated_at FROM accounts WHERE email_confirm_token = $1", confirmToken)
	al := AccountList{}
	return al.FromDBRow(row)
}

// DBInsertAccount inserts an Account into the DB
func DBInsertAccount(a Account) Account {
	db := connect()
	row := db.QueryRow(`INSERT INTO accounts
		(id, email, email_confirm_token, email_confirmed, username, ip, authtoken, created_at, updated_at)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`,
		a.ID, a.Email, a.EmailConfirmToken, a.EmailConfirmed, a.Username, a.IP, a.Authtoken)
	al := AccountList{}
	return al.FromDBRow(row)
}

// DBUpdateAccount updates an Account in the DB
func DBUpdateAccount(n Account) Account {
	db := connect()
	row := db.QueryRow(`UPDATE accounts SET
		username = $2, email = $3, email_confirm_token = $4, email_confirmed = $5, updated_at = NOW()
		WHERE id = $1`,
		n.ID, n.Username, n.Email, n.EmailConfirmToken, n.EmailConfirmed)
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

// DBUpdateNode updates a Node in the DB
func DBUpdateNode(n Node) Node {
	db := connect()
	row := db.QueryRow(`UPDATE nodes SET
		name = $2, url = $3, hostname = $4, port = $5, status = $6, public = $7, updated_at = NOW()
		WHERE id = $1`,
		n.ID, n.Name, n.URL, n.Hostname, n.Port, n.Status, n.Public)
	nl := NodeList{}
	return nl.FromDBRow(row)
}
