package grpcserver

import "google.golang.org/grpc"

func insertIntoUnaryServerInterceptors(
	interceptors []grpc.UnaryServerInterceptor,
	interceptor grpc.UnaryServerInterceptor,
	idx int,
) {
	interceptors = append(interceptors[:idx+1], interceptors[idx:]...)
	interceptors[idx] = interceptor
}

func insertIntoStreamServerInterceptors(
	interceptors []grpc.StreamServerInterceptor,
	interceptor grpc.StreamServerInterceptor,
	idx int,
) {
	interceptors = append(interceptors[:idx+1], interceptors[idx:]...)
	interceptors[idx] = interceptor
}
