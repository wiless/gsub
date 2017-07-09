// Package gsub provides functions to start  web & subscription services
package gsub

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router *mux.Router
var webrootdir = "./"

//StartWebService start a http server at "root" location
func StartWebService(root string) {
	webrootdir = root
	adrs, _ := net.InterfaceAddrs()
	for _, adr := range adrs {
		fmt.Printf("\n Open http://%v%v", strings.Split(adr.String(), "/")[0], ":8888")
	}

	log.Println("Starting Web Server...")
	// http.ListenAndServe(":8888", http.FileServer(http.Dir(root)))

	router = mux.NewRouter().StrictSlash(true)

	// Set API paths
	// router.HandleFunc("/api/apps", handleApps).Methods("GET")
	// router.HandleFunc("/api/apps/{appid:[0-9]+}", getApp).Methods("GET")

	//Set Statick service
	//Approach 1
	router.PathPrefix("/www").HandlerFunc(serverootfiles)
	//Approah 2
	//router.PathPrefix("/www").Handler(http.StripPrefix("/www", http.FileServer(http.Dir("."))))

	// Set API paths
	router.HandleFunc("/", welcomePage)

}

func StartServer() {
	log.Fatal(http.ListenAndServe(":8888", router))
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	log.Print("Some HTTP Client connected ", r.Host)
	log.Println(" But I dont know what to tell him back")
	fmt.Fprintf(w, "Welcome Mr. %s", r.Host)
	// w.WriteHeader(http.StatusAccepted)
}

// func Attach(l *Logger){
// 	router
// }
