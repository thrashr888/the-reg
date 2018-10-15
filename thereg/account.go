package thereg

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	petname "github.com/dustinkirkland/golang-petname"
)

var accountID = "nyft708say7f"
var mySigningKey = []byte("nyv0ny790grewty9er0wyn8t390r2y5t347902b67tqw9067bwqt80c6bdtw78q0")

// Account represents... the user
type Account struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	EmailConfirmToken string `json:"email_confirm_token"`
	EmailConfirmed    bool   `json:"email_confirmed"`
	Username          string `json:"username"`
	IP                string `json:"ip"`
	Authtoken         string `json:"auth_token"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// AccountList is a collection of Accounts
type AccountList struct {
	Accounts []Account `json:"accounts"`
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

// FromJSONBody returns a Account from a JSON Body
func (accounts *AccountList) FromJSONBody(res *http.Response) Account {
	decoder := json.NewDecoder(res.Body)
	var a Account
	decoder.Decode(&a)
	return a
}

// Create makes a new Account
// POST /account
func (accounts *AccountList) Create() Account {
	id := createID()
	username := petname.Generate(3, "-")
	authToken := createToken(id, username)

	params := Account{
		ID:             id,
		Username:       username,
		Authtoken:      authToken,
		EmailConfirmed: false,
	}

	DBInsertAccount(params)

	// return Account from DB
	n := DBGetAccount(id)
	return n
}

// Read returns the current User
// GET /account
func (accounts *AccountList) Read(account Account) Account {
	n := DBGetAccount(account.ID)
	return n
}

// Confirm validates the Users's email
// GET /account/confirm/:confirmToken
func (accounts *AccountList) Confirm(confirmToken string) Account {
	account := DBGetAccountByConfirmToken(confirmToken)

	// update the Account
	if account.ID != "" {
		account.EmailConfirmToken = ""
		account.EmailConfirmed = true
	}
	DBUpdateAccount(account)

	// return Account from DB
	n := DBGetAccount(account.ID)
	return n
}

// Update changes the account info
// PATCH /account
func (accounts *AccountList) Update(account Account, params Account) Account {
	// update the Account
	if params.Email != "" {
		account.Email = params.Email
		account.EmailConfirmToken = createID()
		account.EmailConfirmed = false
	}
	if params.Username != "" {
		account.Username = params.Username
	}
	DBUpdateAccount(account)

	// TODO send email with confirm token link

	// return Account from DB
	n := DBGetAccount(account.ID)
	return n
}

// Delete removes the User
// DELETE /account
func (accounts *AccountList) Delete(account Account) string {
	return "yb7fd0as"
}

func createToken(id string, username string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = false
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}
