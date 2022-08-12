package grpcserver

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

func insertIntoUnaryServerInterceptors(
	interceptors []grpc.UnaryServerInterceptor,
	interceptor grpc.UnaryServerInterceptor,
	idx int,
) []grpc.UnaryServerInterceptor {
	interceptors = append(interceptors[:idx+1], interceptors[idx:]...)
	interceptors[idx] = interceptor

	return interceptors
}

func insertIntoStreamServerInterceptors(
	interceptors []grpc.StreamServerInterceptor,
	interceptor grpc.StreamServerInterceptor,
	idx int,
) []grpc.StreamServerInterceptor {
	interceptors = append(interceptors[:idx+1], interceptors[idx:]...)
	interceptors[idx] = interceptor

	return interceptors
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
