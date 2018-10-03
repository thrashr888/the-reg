package thereg

// Client represents the API client
type Client struct {
}

// GET /node
func GetNodes() NodeList {
	var n = []Node{
		Node{
			ID:       "dff8522fe5dc",
			Name:     "first-deployment",
			URL:      "string",
			Hostname: "76.87.249.25",
			Port:     "8080:80",
			Status:   "UP",
			Public:   true,
		},
		Node{
			ID:       "cf3f7336b1e0",
			Name:     "http",
			URL:      "string",
			Hostname: "76.87.249.25",
			Port:     "80",
			Status:   "UP",
			Public:   true,
		},
	}
	return NodeList{
		Nodes: n,
	}
}

// POST /node
func CreateNode(params Node) Node {
	url := createURL(params.Name, "thrashr888", params.Port)
	return Node{
		ID:       "dff8522fe5dc",
		Name:     params.Name,
		URL:      url,
		Hostname: params.Hostname,
		Port:     params.Port,
		Status:   "UP",
		Public:   true,
	}
}

// GET /node/:id
func GetNode(id string) Node {
	return Node{
		ID:       id,
		Name:     "first-deployment",
		URL:      "first-deployment.thrashr888.the-reg.link:80",
		Hostname: "76.87.249.25",
		Port:     "8080:80",
		Status:   "UP",
		Public:   true,
	}
}

// PATCH /node/:id
func UpdateNode(id string, params Node) Node {
	url := createURL(params.Name, "thrashr888", params.Port)
	return Node{
		ID:       id,
		Name:     params.Name,
		URL:      url,
		Hostname: params.Hostname,
		Port:     params.Port,
		Status:   params.Status,
		Public:   true,
	}
}

// DELETE /node/:id
func DeleteNode(id string) string {
	return id
}

// POST /account
func CreateAccount() Account {
	return Account{
		ID:       "yb7fd0as",
		Username: "full-buffallo-hotness",
	}
}

// GET /account
func GetAccount() Account {
	return Account{
		ID:       "yb7fd0as",
		Email:    "thrashr888@gmail.com",
		Username: "thrashr888",
	}
}

// PATCH /account
func UpdateAccount(params Account) Account {
	return Account{
		ID:       "yb7fd0as",
		Email:    params.Email,
		Username: params.Username,
	}
}

// DELETE /account
func DeleteAccount() string {
	return "yb7fd0as"
}

func getAuth() string {
	return "Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7"
}
