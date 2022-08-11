package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldResponse struct {
	Message string
}

func main() {
	//Profile
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	port := 8081
	http.HandleFunc("/helloworld", helloWorldHandler)
	log.Printf("Server starting on port %v\n", port)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%v", port), "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {

	url := *r.URL
	if url.Scheme != "http" {
		url.Scheme = "https"
		url.Host = r.Host
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
	}

	response := helloWorldResponse{Message: "hello world"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Ooops!")
	}

	fmt.Fprint(w, string(data))
}
