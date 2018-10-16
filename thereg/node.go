package thereg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	shortid "github.com/teris-io/shortid"
)

// Node represents a single networked resource
type Node struct {
	ID        string   `json:"id"`
	AccountID string   `json:"account_id"`
	Name      string   `json:"name"`
	URL       string   `json:"url"`
	Hostname  string   `json:"hostname"`
	Port      string   `json:"port"`
	Status    string   `json:"status"`
	Public    bool     `json:"public"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// NodeList is a collection of Nodes
type NodeList struct {
	Nodes []Node `json:"nodes"`
}

// INSERT INTO nodes
// (id, account_id, name, port, url, hostname, status)
// VALUES
// ('c65e2d0eb499', 'nyft708say7f', 'redis', 6973, 'redis.c65e2d0eb499.the-reg.link', '76.87.249.25', 'UP');

// FromDBRows returns a NodeList from sql.Rows
func (nodes *NodeList) FromDBRows(rows *sql.Rows) NodeList {
	n := []Node{}
	for rows.Next() {
		var node Node
		if err := rows.Scan(
			&node.ID,
			&node.AccountID,
			&node.Name,
			&node.URL,
			&node.Hostname,
			&node.Port,
			&node.Status,
			&node.Public,
			&node.CreatedAt,
			&node.UpdatedAt); err != nil {
			log.Println(err.Error())
		}
		n = append(n, node)
	}
	return NodeList{Nodes: n}
}

// FromDBRow returns a NodeList from sql.Row
func (nodes *NodeList) FromDBRow(row *sql.Row) Node {
	var node Node
	if err := row.Scan(
		&node.ID,
		&node.AccountID,
		&node.Name,
		&node.URL,
		&node.Hostname,
		&node.Port,
		&node.Status,
		&node.Public,
		&node.CreatedAt,
		&node.UpdatedAt); err != nil {
		log.Println(err.Error())
	}
	return node
}

// ListFromJSONBody returns a NodeList from a JSON Body
func (nodes *NodeList) ListFromJSONBody(res *http.Response) NodeList {
	decoder := json.NewDecoder(res.Body)
	var n NodeList
	decoder.Decode(&n)
	return n
}

// FromJSONBody returns a Node from a JSON Body
func (nodes *NodeList) FromJSONBody(res *http.Response) Node {
	decoder := json.NewDecoder(res.Body)
	var n Node
	decoder.Decode(&n)
	return n
}

// Index returns the user's nodes
// GET /api/node
func (nodes *NodeList) Index(account Account) NodeList {
	n := DBGetNodes(account.ID)
	return n
}

// Create adds a Node
// POST /api/node
func (nodes *NodeList) Create(account Account, params Node) Node {
	// check for existing Node
	old := DBGetNode(params.Name)
	if old.ID != "" {
		return old
	}

	// make a Node
	id := createID()
	params.ID = id
	params.AccountID = account.ID
	params.Public = true
	params.Status = "UP"
	if params.Name == "" {
		params.Name = id
	}
	params.URL = createURL(params.Name, account.Username, params.Port)
	DBInsertNode(params)

	// return Node from DB
	n := DBGetNode(id)
	return n
}

// Read returns a node
// GET /api/node/:id
func (nodes *NodeList) Read(id string) Node {
	n := DBGetNode(id)
	return n
}

// Update changes a node
// PATCH /api/node/:id
func (nodes *NodeList) Update(account Account, id string, params Node) Node {
	// check for existing Node
	node := DBGetNode(id)
	if node.ID == "" {
		return Node{}
	}

	// update the Node
	if params.Name != "" {
		node.Name = params.Name
	}
	node.URL = createURL(node.Name, account.Username, node.Port)
	DBUpdateNode(node)

	// return Node from DB
	n := DBGetNode(id)
	return n
}

// Delete changes a node
// DELETE /api/node/:id
func (nodes *NodeList) Delete(id string) string {
	return id
}

func createURL(name string, username string, port string) string {
	// redis.full-buffallo-hotness.the-reg.link:6379
	return fmt.Sprintf("%v.%v.the-reg.link:%v", name, username, port)
}

func createID() string {
	ret, _ := shortid.Generate()
	return strings.TrimSpace(ret)
}
