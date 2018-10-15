package thereg

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
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
	var params Node
	if port != "" {
		params = Node{Name: name, Hostname: hostnameOrPort, Port: port}
	} else {
		res, _ := http.Get("https://api.ipify.org")
		ip, _ := ioutil.ReadAll(res.Body)
		params = Node{Name: name, Hostname: string(ip), Port: hostnameOrPort}
	}
	node := CreateNode(params)
	fmt.Println(node.ID)
}
func create() {
	if checkAuthToken() {
		_, err := readAuthToken()
		if err == nil {
			fmt.Println("Client already registered")
			return
		}
	}

	// create a new account
	account := CreateAccount(Account{})
	fileContent := fmt.Sprintf("authtoken: %s", account.Authtoken)
	writeAuthToken(fileContent)
	os.Setenv("THE_REG_TOKEN", account.Authtoken)
	fmt.Printf("echo \"authtoken: %s\" > ~/.thereg.yml\nTHE_REG_TOKEN=%s\n", account.Authtoken, account.Authtoken)
}
func get(id string) {
	node := GetNode(id)
	fmt.Println(node.URL)
}
func ip() {
	fmt.Println(GetIP())
}
func list() {
	nodeList := GetNodes()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Host", "Port", "Status", "Public"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, n := range nodeList.Nodes {
		table.Append([]string{n.ID, n.Name, n.Hostname, n.Port, n.Status, strconv.FormatBool(n.Public)})
	}

	table.Render()
}
func login(authtoken string) {
	if checkAuthToken() {
		_, err := readAuthToken()
		if err == nil {
			fmt.Println("Client already logged in")
			return
		}
	}

	fileContent := fmt.Sprintf("authtoken: %s", authtoken)
	writeAuthToken(fileContent)
	os.Setenv("THE_REG_TOKEN", authtoken)
	fmt.Printf("echo \"authtoken: %s\" > ~/.thereg.yml\nTHE_REG_TOKEN=%s\n", authtoken, authtoken)
}
func logout() {
	removeAuthToken()
	fmt.Println("Logged out")
}
func me() {
	account := GetAccount()
	fmt.Println(account.Username)
}
func nameNode(id string, name string) {
	params := Node{Name: name}
	node := UpdateNode(id, params)
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
func help() {
	log.Println(`A global service registry. Free public forwarding. $6.99/mo for unlimited private.

Usage:
    $ reg <command> [options...]

Commands:

    account - "reg account new :username :email" sign up for an account
    add - "reg add :name [hostname] :port" add a node
    create - get a user token
    get - "reg get :name" Get a service url
    help - show this list
    ip - get your public ip address
    list - list your nodes
    login - save your auth token
    me - your username
    name - "reg name :id :name" name a node
    start - attempt to reset status to "UP"
    server - run the web service

Examples:

    $ reg create
    $ reg account new <username> <email>
    $ reg me
    $ reg add redis 6379
    $ reg list
    $ reg get redis`)
}

// Run sets up the CLI parsing
func Run() {
	flag.Parse()

	action := flag.Arg(0)

	switch action {
	case "account":
		username := flag.Arg(1)
		email := flag.Arg(2)
		account(username, email)
		break
	case "add":
		name := flag.Arg(1)
		hostnameOrPort := flag.Arg(2)
		port := flag.Arg(3)
		add(name, hostnameOrPort, port)
		break
	case "create":
		create()
		break
	case "get":
		name := flag.Arg(1)
		get(name)
		break
	case "help":
		help()
		break
	case "ip":
		ip()
		break
	case "list":
		list()
		break
	case "login":
		authToken := flag.Arg(1)
		login(authToken)
		break
	case "logout":
		logout()
		break
	case "me":
		me()
		break
	case "name":
		id := flag.Arg(1)
		name := flag.Arg(2)
		nameNode(id, name)
		break
	case "serve":
		Serve()
		os.Exit(0)
		break
	case "start":
		id := flag.Arg(1)
		start(id)
		break
	default:
		fmt.Println("Command not found.")
	}
}

// GetIP returns the user's IP string
func GetIP() string {
	res, _ := http.Get("https://api.ipify.org")
	ip, _ := ioutil.ReadAll(res.Body)
	return string(ip)
}

func checkAuthToken() bool {
	fileName, _ := homedir.Expand("~/.thereg.yml")
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return true
	}
	return false
}
func writeAuthToken(token string) {
	fileName, _ := homedir.Expand("~/.thereg.yml")
	fileHandle, _ := os.Create(fileName)
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()

	fmt.Fprintln(writer, token)
	writer.Flush()
}
func readAuthToken() (string, error) {
	// get the file contents
	fileName, _ := homedir.Expand("~/.thereg.yml")
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	line := string(b)

	// find the token
	r, _ := regexp.Compile("authtoken: (.+)")
	match := r.FindStringSubmatch(line)
	if len(match) > 0 {
		return match[1], nil
	}

	return "", errors.New("authToken not found")
}
func removeAuthToken() {
	if checkAuthToken() {
		fileName, _ := homedir.Expand("~/.thereg.yml")
		var err = os.Remove(fileName)
		if err != nil {
			// fmt.Println(err.Error())
		}
	}
}
