package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	api "github.com/thrashr888/the-reg/thereg"
	cli "github.com/thrashr888/the-reg/thereg"
)

func help() {
	log.Println(`Usage reg command [options...]
A global service registry. Free public forwarding. $6.99/mo for unlimited private.

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
		server - run the web service`)
}

func main() {
	helpMe := flag.Bool("h", false, "help")
	serve := flag.Bool("s", false, "server")

	flag.Parse()

	if *helpMe {
		help()
		os.Exit(0)
	}

	cli.Register("account")

	if *serve {
		api.Serve()
		os.Exit(0)
	}
}

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())
}
