package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Log the incoming request
		log.Printf("[gRPC] --> Method: %s | Started at: %s", info.FullMethod, start.Format(time.RFC3339))

		// Call the handler to complete the RPC
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log the response
		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("[gRPC] <-- Method: %s | Duration: %v | Status: %s | Error: %v",
				info.FullMethod, duration, st.Code(), err)
		} else {
			log.Printf("[gRPC] <-- Method: %s | Duration: %v | Status: OK",
				info.FullMethod, duration)
		}

		return resp, err
	}
}
