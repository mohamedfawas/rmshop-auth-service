package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// log grpc req ans responses
func LoggingInterceptor(ctx context.Context,
	req interface{}, // data sent to the grpc server
	info *grpc.UnaryServerInfo, // Metadata about the gRPC method being called, such as its name
	handler grpc.UnaryHandler, // actual function that handles the request and generates the response
) (interface{}, error) {

	start := time.Now()

	// Log the method being called. This helps in tracking which gRPC method is being executed.
	log.Printf("Request - Method: %s", info.FullMethod)

	// Pass the request to the actual handler to process it
	// `handler` is the function that executes the business logic for this gRPC call
	resp, err := handler(ctx, req)

	// Measure the time taken to process the request
	duration := time.Since(start)

	if err != nil {

		// Extract the error details (message and code) using gRPC's status package
		st, _ := status.FromError(err)

		// Log the method name, duration, error message, and error code
		log.Printf("Response - Method: %s, Duration: %v, Error: %s, Code: %s",
			info.FullMethod, duration, st.Message(), st.Code())
	} else {
		// If no error occurred, log the method name, duration, and status code
		log.Printf("Response - Method: %s, Duration: %v, Status: %s",
			info.FullMethod, duration, codes.OK)
	}

	return resp, err
}
