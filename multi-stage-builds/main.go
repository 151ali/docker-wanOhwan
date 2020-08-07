package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/",roothandler)
	fmt.Println("server running on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err!=nil {
		log.Fatal(err)
	}
}

func roothandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Hello Gorilla :)")
}