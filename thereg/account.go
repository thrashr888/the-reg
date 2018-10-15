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

	return Account{
		ID:        id,
		Username:  username,
		Authtoken: authToken,
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
