package thereg

import (
	"fmt"
)

// account - `reg account new :username :email` sign up for an account
// add - `reg add :name [hostname] :port` add a node
// create - get a user token
// get - `reg get :name` Get a service url
// help - show this list
// ip - get your public ip address
// list - list your nodes
// login - save your auth token
// me - your username
// name - `reg name :id :name` name a node
// start - attempt to reset status to "UP"
// server - run the web service

func account(username string, email string) {
	params := Account{Email: email, Username: username}
	UpdateAccount(params)
	fmt.Println(`Account created. Check your email to log in at https://www.the-reg.link/`)
}
func add(name string, hostnameOrPort string, port string) {
	params := Node{Name: name, Hostname: hostnameOrPort, Port: port}
	node := CreateNode(params)
	fmt.Println(node.ID, node.Name)
}
func create() {
	account := CreateAccount()
	fmt.Printf(`# echo "authtoken: %s" > ~/.thereg.yml
# export THE_REG_TOKEN=%s`, account.Authtoken, account.Authtoken)
}
func get(id string) {
	node := GetNode(id)
	fmt.Println(node.URL)
}
func ip() {
	fmt.Println(`76.87.249.25`)
}
func list() {
	nodeList := GetNodes()
	fmt.Println(`ID             NAME               HOST          PORT      STATUS   AGE   PUBLIC   TAGS`)

	for _, n := range nodeList.Nodes {
		fmt.Printf(`%s %s %s %s UP 1h Y\n`, n.ID, n.Name, n.Hostname, n.Port)
	}
}
func login(authtoken string) {
	fmt.Printf(`# echo "authtoken: %s" > ~/.thereg.yml
# export THE_REG_TOKEN=%s`, authtoken, authtoken)
}
func me() {
	account := GetAccount()
	fmt.Println(account.Username)
}
func name(idOrName string) {
	node := GetNode(idOrName)
	fmt.Println(node.URL)
}
func start(idOrName string) {
	params := Node{Status: "UP"}
	node := UpdateNode(idOrName, params)

	if node.Status == "UP" {
		fmt.Println(`Local port 8081 found.`)
	} else {
		fmt.Println(`Local port 8081 not found. Try restarting your server.`)
	}
}

// Register sets up the CLI parsing
func Register(args string) {

	action := "create"

	parsed := make(map[string]string)
	parsed["id"] = "c65e2d0eb499"
	parsed["username"] = "thrashr888"
	parsed["email"] = "thrashr888@gmail.com"
	parsed["name"] = "redis"
	parsed["hostnameOrPort"] = "example.com"
	parsed["port"] = "8080"
	parsed["authToken"] = "Sc1VvxLceT5MrMaAjoio_2uLEttzm4com5xT1zh7D7"

	switch action {
	case "account":
		account(parsed["username"], parsed["email"])
		break
	case "add":
		add(parsed["name"], parsed["hostnameOrPort"], parsed["port"])
		break
	case "create":
		create()
		break
	case "get":
		get(parsed["name"])
		break
	case "ip":
		ip()
		break
	case "list":
		list()
		break
	case "login":
		login(parsed["authToken"])
		break
	case "me":
		me()
		break
	case "name":
		name(parsed["id"])
		break
	case "start":
		start(parsed["id"])
		break
	default:
		// Give an error message.
	}
}
