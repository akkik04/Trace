package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/akkik04/Trace/services/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.LogServiceClient

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

	// set CORS headers -- allow all origins, methods, and headers for the sake of simplicity.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// handle preflight requests.
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	// get the IP address of the task (both the collector and ingestor are running as tasks in the same ECS service).
	// so we can use the task IP to communicate with the ingestor service.
	ecsIP, err := getTaskPublicIP()
	if err != nil {
		http.Error(w, "Error retrieving ECS Task Public-IP", http.StatusInternalServerError)
		return
	}

	// construct the ingestor address (assuming port 8082), and create a gRPC connection to the ingestor service.
	ingestorAddr := fmt.Sprintf("%s:8082", ecsIP)
	conn, err := grpc.NewClient(ingestorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error creating connection to ingestor: ", err)
		panic(err)
	}
	client = proto.NewLogServiceClient(conn)

	// store the log entry in the in-memory data store.
	logStore = append(logStore, logEntry)
	res, err := client.SendLog(context.Background(), &proto.LogEntry{
		Ip:        logEntry.IP,
		User:      logEntry.User,
		Timestamp: logEntry.Timestamp.String(),
		Method:    logEntry.Method,
		Endpoint:  logEntry.Endpoint,
		Protocol:  logEntry.Protocol,
		Status:    int32(logEntry.Status),
		Size:      int32(logEntry.Size),
	})
	if err != nil {
		http.Error(w, "Error sending log to ingestor", http.StatusInternalServerError)
		fmt.Println("Error sending log to ingestor: ", err)
		return
	}
	fmt.Printf("Response from ingestor: %v\n", res.GetMessage())
	w.WriteHeader(http.StatusOK)
}

// helper function to extract meaningful metrics from the log message.
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

// helper function to get the public IP address of the ECS task.
func getTaskPublicIP() (string, error) {

	// load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	// retrieve necessary environment variables.
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	clusterName := os.Getenv("AWS_ECS_CLUSTER_NAME")
	serviceName := os.Getenv("AWS_ECS_SERVICE_NAME")

	// initialize AWS session.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	}))

	// create ECS service client.
	ecsSvc := ecs.New(sess)

	// list the tasks for the specified service.
	listTasksOutput, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:     aws.String(clusterName),
		ServiceName: aws.String(serviceName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list tasks: %v", err)
	}

	// ensure there is at least one task running.
	if len(listTasksOutput.TaskArns) == 0 {
		return "", fmt.Errorf("no tasks found for service %s in cluster %s", serviceName, clusterName)
	}

	// describe the task to get the ENI (Elastic Network Interface)
	taskDesc, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(clusterName),
		Tasks:   listTasksOutput.TaskArns,
	})
	if err != nil {
		return "", fmt.Errorf("failed to describe tasks: %v", err)
	}

	// create EC2 service client.
	ec2Svc := ec2.New(sess)
	var publicIP string

	// extracting the ENI ID and retrieveing the public IP so that we can use it during gRPC comms w/ the ingestor service.
	for _, task := range taskDesc.Tasks {
		for _, attachment := range task.Attachments {
			if *attachment.Type == "ElasticNetworkInterface" {
				for _, detail := range attachment.Details {
					if *detail.Name == "networkInterfaceId" {
						eniID := *detail.Value

						// describe the ENI to get the public IP address, if it exists.
						eniDesc, err := ec2Svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
							NetworkInterfaceIds: []*string{aws.String(eniID)},
						})
						if err != nil {
							return "", fmt.Errorf("failed to describe network interface: %v", err)
						}

						if len(eniDesc.NetworkInterfaces) > 0 && eniDesc.NetworkInterfaces[0].Association != nil {
							publicIP = *eniDesc.NetworkInterfaces[0].Association.PublicIp

							// if the public IP is found, return it.
							return publicIP, nil
						}
					}
				}
			}
		}
	}
	return "", fmt.Errorf("public IP not found for any task in service %s", serviceName)
}

// main.
func main() {

	// run the /log_collector route on port 8080 to handle incoming requests from the web (frontend) simulator.
	http.HandleFunc("/log_collector", logsHandler)
	fmt.Println("Go microservice running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
