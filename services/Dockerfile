# Use the official Golang image based on Go 1.22
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the proto files
COPY proto/ proto/

# Copy the entire services directory
COPY . .

# Build arguments for specifying which service to build
ARG SERVICE
ENV SERVICE=${SERVICE}

# Build the respective service
RUN cd ${SERVICE} && go build -o main .

# Command to run the respective service
CMD ["sh", "-c", "cd ./${SERVICE} && ./main"]
