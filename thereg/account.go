package thereg

import (
	"database/sql"
	"log"
)

var accountID = "nyft708say7f"

// Account represents... the user
type Account struct {
	ID                string
	Email             string
	EmailConfirmToken string
	EmailConfirmed    bool
	Username          string
	IP                string
	Authtoken         string
	CreatedAt         string
	UpdatedAt         string
}

// AccountList is a collection of Accounts
type AccountList struct {
	Accounts []Account
}

// INSERT INTO accounts
// (id, username, email, email_confirmation_token, ip, auth_token)
// VALUES
// ('nyft708say7f', 'thrashr888', '', '', '76.87.249.25', 'Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7');

// FromDBRows returns a NodeList from sql.Rows
func (accounts *AccountList) FromDBRows(rows *sql.Rows) AccountList {
	n := []Account{}
	for rows.Next() {
		// var account =: n.FromDBRow()
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.Email,
			&account.EmailConfirmToken,
			&account.EmailConfirmed,
			&account.Username,
			&account.IP,
			&account.Authtoken,
			&account.CreatedAt,
			&account.UpdatedAt); err != nil {
			log.Println(err.Error())
		}
		n = append(n, account)
	}
	return AccountList{Accounts: n}
}

// FromDBRow returns a AccountList from sql.Row
func (accounts *AccountList) FromDBRow(row *sql.Row) Account {
	var account Account
	if err := row.Scan(
		&account.ID,
		&account.Email,
		&account.EmailConfirmToken,
		&account.EmailConfirmed,
		&account.Username,
		&account.IP,
		&account.Authtoken,
		&account.CreatedAt,
		&account.UpdatedAt); err != nil {
		log.Println(err.Error())
	}
	return account
}

// Create makes a new Account
// POST /account
func (accounts *AccountList) Create() Account {
	return Account{
		ID:       "yb7fd0as",
		Username: "full-buffallo-hotness",
	}
}

// Read returns the current User
// GET /account
func (accounts *AccountList) Read() Account {
	n := DBGetAccount(accountID)
	return n
}

// Confirm validates the Users's email
// GET /account/confirm/:token
func (accounts *AccountList) Confirm(token string) Account {
	return Account{
		ID:       "yb7fd0as",
		Email:    "thrashr888@gmail.com",
		Username: "thrashr888",
	}
}

// Update changes the account info
// PATCH /account
func (accounts *AccountList) Update(params Account) Account {
	return Account{
		ID:       "yb7fd0as",
		Email:    params.Email,
		Username: params.Username,
	}
}

// Delete removes the User
// DELETE /account
func (accounts *AccountList) Delete() string {
	return "yb7fd0as"
}
