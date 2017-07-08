package main

import (
	"fmt"
	"time"

	"github.com/wiless/gsub"
)

var stime time.Time

func init() {
	Start()
}

func Start() {
	stime = time.Now()
	fmt.Printf("\n======== Started %v ===========\n", stime)
}

func GStop() {
	fmt.Printf("\n======== Code Run for  %v ===========\n\n", time.Since(stime))
}

func main() {
	defer GStop() // Ensure to call Stop() function when main() routine ends/terminates
	gsub.StartWebService("./")

	fmt.Printf("Hello How are you..")

	// GStop()
}
