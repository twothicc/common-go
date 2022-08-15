package grpcclient

import "google.golang.org/grpc"

func insertIntoUnaryClientInterceptors(
	interceptors []grpc.UnaryClientInterceptor,
	interceptor grpc.UnaryClientInterceptor,
	idx int,
) []grpc.UnaryClientInterceptor {
	interceptors = append(interceptors[:idx+1], interceptors[idx:]...)
	interceptors[idx] = interceptor

	return interceptors
}

func insertIntoStreamClientInterceptors(
	interceptors []grpc.StreamClientInterceptor,
	interceptor grpc.StreamClientInterceptor,
	idx int,
) []grpc.StreamClientInterceptor {
	interceptors = append(interceptors[:idx+1], interceptors[idx:]...)
	interceptors[idx] = interceptor

	return interceptors
}
