package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8000", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	current := t.Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "Hello world "+current)
}
