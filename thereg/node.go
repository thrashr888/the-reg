package thereg

import (
	"database/sql"
	"log"
)

// Node represents a single networked resource
type Node struct {
	ID        string
	AccountID string
	Name      string
	URL       string
	Hostname  string
	Port      string
	Status    string
	Public    bool
	CreatedAt string
	UpdatedAt string
}

// NodeList is a collection of Nodes
type NodeList struct {
	Nodes []Node
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

// Index returns the user's nodes
// GET /node
func (nodes *NodeList) Index() NodeList {
	n := DBGetNodes(accountID)
	return n
}

// Create adds a Node
// POST /node
func (nodes *NodeList) Create(params Node) Node {
	url := createURL("dff8522fe5dc", "thrashr888", params.Port)
	return Node{
		ID:       "dff8522fe5dc",
		URL:      url,
		Hostname: params.Hostname,
		Port:     params.Port,
		Status:   "UP",
		Public:   true,
	}
}

// Read returns a node
// GET /node/:id
func (nodes *NodeList) Read(id string) Node {
	n := DBGetNode(id)
	return n
}

// Update changes a node
// PATCH /node/:id
func (nodes *NodeList) Update(id string, params Node) Node {
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

// Delete changes a node
// DELETE /node/:id
func (nodes *NodeList) Delete(id string) string {
	return id
}

func createURL(name string, username string, port string) string {
	// redis.full-buffallo-hotness.the-reg.link:6379
	return name + "." + username + "the-reg.link:" + port
}
