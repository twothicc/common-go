# Logger

This package handles the initialization and configuration of a zap logger, with utility functions that handles adding and propagating values to be logged via the Go context.

It is meant to be used with the [grpcserver](https://github.com/twothicc/common-go/grpcserver) package, of which one of the uses is to provide important log fields via its chain of middlewares.

# Usage

## Initialize logger

At the start of your program, call `logger.InitLogger(<isTest>)` where `<isTest>` refers to whether the development environment is test. This will initialize and configure the global zap logger.

## Logging with context

Call e.g. `logger.WithContext(ctx).Info("hello world", zap.String("msg", "hello world"))` to log with context.

By default, fields specified in defaultLogFields (if available) will be prepended to the log message.

```
var defaultLogFields = []string{
	"trace.traceid",
	"trace.spanid",
	"grpc.request.service",
	"grpc.request.method",
}
```

## Usage with grpcserver package

Take note to initialize the logger as shown in an earlier example before starting the gRPC server.
