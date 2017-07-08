// Package gsub provides functions to start  web & subscription services
package gsub

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

//StartWebService start a http server at "root" location
func StartWebService(root string) {

	adrs, _ := net.InterfaceAddrs()
	for _, adr := range adrs {
		fmt.Printf("\n Open http://%v:%v", strings.Split(adr.String(), "/")[0], ":8080")
	}

	// err = http.ListenAndServe(portid, http.FileServer(http.Dir(servedir)))

	log.Println("Starting Web Server ")
	// http.HandleFunc("/", landingPage)
	// http.HandleFunc("/www", http.FileServer(http.Dir(root)))

	http.ListenAndServe(":8080", http.FileServer(http.Dir(root)))
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	log.Print("Some HTTP Client connected ", r.Host)
	log.Println(" But I dont know what to tell him back")
	fmt.Fprintf(w, "Welcome Mr. %s", r.Host)
	// w.WriteHeader(http.StatusAccepted)
}
