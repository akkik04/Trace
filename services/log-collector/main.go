package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// create a struct to represent the log entry.
type LogEntry struct {
	IP        string
	User      string
	Timestamp time.Time
	Method    string
	Endpoint  string
	Protocol  string
	Status    int
	Size      int
}

// create a struct to represent the incoming log message.
type LogMessage struct {
	Message string `json:"message"`
}

// create an in-memory data store to store the log entries.
var logStore []LogEntry

// pre-compile the regx to extract meaningful metrics from the log message.
// we know the incoming log message format is CLF, so we can use a regex to extract the metrics.
var logRegex = regexp.MustCompile(`(?P<IP>\d+\.\d+\.\d+\.\d+) - (?P<User>[^\s]+) \[(?P<Timestamp>[^\]]+)\] "(?P<Method>[A-Z]+) (?P<Endpoint>[^\s]+) (?P<Protocol>[^\s]+)" (?P<Status>\d+) (?P<Size>\d+)`)

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

	// parse the log message from the JSON body.
	var logMsg LogMessage
	err = json.Unmarshal(body, &logMsg)
	if err != nil {
		http.Error(w, "Error parsing JSON request body", http.StatusBadRequest)
		return
	}

	// extract the message field.
	message := logMsg.Message
	if message == "" {
		http.Error(w, "Missing or empty 'message' field in JSON body", http.StatusBadRequest)
		return
	}

	// extract meaningful metrics from the log message using regex.
	logEntry, err := extractLogMetrics(message)
	if err != nil {
		http.Error(w, "Error extracting log metrics", http.StatusBadRequest)
		return
	}

	// print the extracted log entry.
	fmt.Printf("Extracted Log Entry: %+v\n", logEntry)

	// store the log entry in the in-memory data store.
	logStore = append(logStore, logEntry)
	w.WriteHeader(http.StatusOK)
}

// function to extract meaningful metrics from the log message.
func extractLogMetrics(log string) (LogEntry, error) {

	// use the regex to extract the metrics.
	matches := logRegex.FindStringSubmatch(log)
	
	if matches == nil {
		return LogEntry{}, fmt.Errorf("log format not recognized")
	}

	// parse the timestamp.
	timestamp, err := time.Parse("2006-01-02T15:04:05.000Z", matches[3])
	if err != nil {
		return LogEntry{}, fmt.Errorf("timestamp format not recognized")
	}

	// parse the status and size.
	status, err := strconv.Atoi(matches[7])
	if err != nil {
		return LogEntry{}, fmt.Errorf("status format not recognized")
	}
	size, err := strconv.Atoi(matches[8])
	if err != nil {
		return LogEntry{}, fmt.Errorf("size format not recognized")
	}

	// create a log entry.
	logEntry := LogEntry{
		IP:        matches[1],
		User:      matches[2],
		Timestamp: timestamp,
		Method:    matches[4],
		Endpoint:  matches[5],
		Protocol:  matches[6],
		Status:    status,
		Size:      size,
	}

	return logEntry, nil
}

// main.
func main() {

	// run the /log_collector route on port 8080.
	http.HandleFunc("/log_collector", logsHandler)
	fmt.Println("Go microservice running on port 8080")
	http.ListenAndServe(":8080", nil)
}
