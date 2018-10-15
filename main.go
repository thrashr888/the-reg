package main

import (
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	cli "github.com/thrashr888/the-reg/thereg"
)

func main() {
	cli.Run()
}

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())
}
