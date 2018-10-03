package thereg

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
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
	createdAt string
}

// NodeList is a collection of Nodes
type NodeList struct {
	Nodes []Node
}

func (p *Node) save() error {
	filename := "data/" + p.ID + ".txt"
	return ioutil.WriteFile(filename, []byte(p.URL), 0600)
}

func loadNode(id string) (*Node, error) {
	filename := "data/" + id + ".txt"

	url, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Node{ID: id, URL: string(url)}, nil
}

func loadNodes() ([]Node, error) {
	files, err := filepath.Glob("./data/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	nodes := []Node{}
	for _, f := range files {
		_, file := filepath.Split(f)
		id := strings.TrimSuffix(file, path.Ext(file))
		if err != nil {
			continue
		}
		nodes = append(nodes, Node{ID: id})
	}

	return nodes, nil
}

// GET /node
// POST /node
// GET /node/:id
// PATCH /node/:id
// DELETE /node/:id

// dff8522fe5dc   first-deployment   76.87.249.25  8080:80   UP       6h    N
// cf3f7336b1e0   http               76.87.249.25  80        UP       3h    Y
// d39dd625947b   https              76.87.249.25  443       UP       3h    Y
// bc2740d30a5f   httpexposed        76.87.249.25  8081:80   DOWN     2h    Y
// c65e2d0eb499   redis              76.87.249.25  6379      UP       2h    Y

// Index returns the user's nodes
// GET /node
func (nodes *NodeList) Index() NodeList {
	n := DBGetNodes("nyft708say7f")
	// TODO
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
