package main

import (
	"fmt"
	"math/rand"
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
	gsub.StartWebService("./web")
	// tlogger := gusb.NewHttpLogger("Temperature")
	var tlogger gsub.Logger
	tlogger.Start("/monitor")
	go gsub.StartServer()
	queue := tlogger.Queue()

	var lmsg gsub.LogMessage
	for {
		msg := fmt.Sprintf("I am a random number %d ", rand.Intn(100))
		lmsg = gsub.LogMessage{Level: 1, Message: msg}
		queue <- lmsg
		time.Sleep(1 * time.Second)
	}

	// GStop()
}
