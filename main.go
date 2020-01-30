package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", HelloServer)
    http.ListenAndServe(":80", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	if(r.URL.Path[1:] == ""){
    	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	}

}