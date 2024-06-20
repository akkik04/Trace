package main

import (
	"fmt"
	"io"
	"net/http"
)

// function to handle the incoming logs.
func logsHandler(w http.ResponseWriter, r *http.Request) {

	// if the request method is not POST, return an error.
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// read the request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// print the log message.
	fmt.Printf("Received log message: %s\n", body)
	w.WriteHeader(http.StatusOK)
}

// main.
func main() {

	// run the /log_collector route on port 8080.
	http.HandleFunc("/log_collector", logsHandler)
	fmt.Println("Go microservice running on port 8080")
	http.ListenAndServe(":8080", nil)
}
