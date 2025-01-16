FROM golang:1.23.4-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/server/main.go

# Run the application
CMD ["/auth-service"] 