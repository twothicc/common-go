package grpcserver

import (
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
