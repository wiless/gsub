package main

import (
	"fmt"
	"time"
)

var stime time.Time

func init() {

}

func Start() {
	stime = time.Now()
	fmt.Printf("\n======== Started %v ===========\n", stime)
}

func Stop() {
	fmt.Printf("\n======== Started %v ===========\n", time.Since(stime))
}

func main() {
	defer Stop()
	fmt.Printf("Hello How are you..")
}
