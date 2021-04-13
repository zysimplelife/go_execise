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
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := helloWorldResponse{Message: "hello world"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Ooops!")
	}

	fmt.Fprint(w, string(data))
}
