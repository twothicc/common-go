# Grpc Server

This package provides a convenient way to configure and initialize a gRPC server, with debugging and monitoring tools.

The gRPC server will be initialized with these chained middlewares (top -> bottom):

- `grpc_ctxtags` (default): Extracts request information from incoming request payloads into a Tag. This Tag is then added to handler's context.
- `grpc_opentracing` (default): Configured with a jaeger tracer as global OpenTracing tracer. This middleware will extract parent span context from incoming requests, then creates a new span referencing the parent span context. The span context of the new span is then injected into Tag in handler's context.
- `grpc_prometheus` (optional): Creates and monitors server metrics
- `grpc_zap` (default): Configured with common-go logger to log completed gRPC calls. The logger is then populated into the handler's context.
- `grpc_recovery` (default): Configured with default settings to convert panics into gRPC error with `code.Internal`.

The server is configured to listen for interrupt, terminate, quit os signals and will gracefully shutdown the http server running prometheus (if exists) and then finally the gRPC server.

# Usage

## Initialize and run the gRPC server:

```
// Create lambda functions to register services
registerHelloWorldServiceHandler := func(s *grpc.Server) {
    pb.RegisterHelloWorldServiceServer(s, helloworld.NewHelloWorldServer())
}

serverConfig := grpcserver.GetDefaultServerConfigs("myService", "localhost", "8080", false, registerHelloWorldServiceHandler)

grpcserver.InitAndRunGrpcServer(context.Background(), serverConfig)
```

**Note**: The default `ServerConfigs` is configured for the following:

- Enable Prometheus monitoring
- Time to keepalive ping **1hr**
- Timeout to close connection after keepalive ping **10s**
- Max idle connection time **5mins**

## Prometheus metrics

Prometheus is configured to monitor default server metrics such as method handling counter and can be accessed at `<domain>:9090`

```
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 11
# HELP grpc_server_handled_total Total number of RPCs completed on the server, regardless of success or failure.
# TYPE grpc_server_handled_total counter
grpc_server_handled_total{grpc_code="Aborted",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 3
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 0
# HELP grpc_server_msg_received_total Total number of RPC stream messages received on the server.
# TYPE grpc_server_msg_received_total counter
grpc_server_msg_received_total{grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 3
# HELP grpc_server_msg_sent_total Total number of gRPC stream messages sent by the server.
# TYPE grpc_server_msg_sent_total counter
grpc_server_msg_sent_total{grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 3
# HELP grpc_server_started_total Total number of RPCs started on the server.
# TYPE grpc_server_started_total counter
grpc_server_started_total{grpc_method="HelloWorld",grpc_service="datasync.v1.HelloWorldService",grpc_type="unary"} 3
```

## Jaeger UI

Jaeger UI tracks services and visualizes spans within each trace. This provides us a platform to pinpoint failures and identify sources of poor performance by monitoring the spans.

Make sure you have docker installed, then run this:

```
docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.9
```

Then visit `http://localhost:16686` to access the UI
