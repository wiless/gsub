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
	tlogger.Start("/monitor") // The /monitor url will print the stream of text coming from logger server or use a facncy html page ( /www/stream.html JS based event listener)
	go gsub.StartServer()
	queue := tlogger.Queue()

	var lmsg gsub.LogMessage
	for {
		msg := fmt.Sprintf("I am a random number %d ", rand.Intn(100))
		lmsg.Level = 1     // any number , for future use
		lmsg.Message = msg // A string message to be broadcast to client, the server may append timestamps !!
		queue <- lmsg
		time.Sleep(2 * time.Second)
	}

	// GStop()
}
