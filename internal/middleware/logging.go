package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Log request
	log.Printf("Request - Method: %s", info.FullMethod)

	// Handle request
	resp, err := handler(ctx, req)

	// Log response
	duration := time.Since(start)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Response - Method: %s, Duration: %v, Error: %s, Code: %s",
			info.FullMethod, duration, st.Message(), st.Code())
	} else {
		log.Printf("Response - Method: %s, Duration: %v, Status: %s",
			info.FullMethod, duration, codes.OK)
	}

	return resp, err
}
