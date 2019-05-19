package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/getData", GetData)
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Go!"))
}
