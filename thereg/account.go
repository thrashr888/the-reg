package thereg

// Account represents... the user
type Account struct {
	ID                string
	Email             string
	EmailConfirmToken string
	EmailConfirmed    bool
	Username          string
	IP                string
	Authtoken         string
}

// AccountList is a collection of Accounts
type AccountList struct {
	Accounts []Account
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
	return Account{
		ID:       "yb7fd0as",
		Email:    "thrashr888@gmail.com",
		Username: "thrashr888",
	}
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
