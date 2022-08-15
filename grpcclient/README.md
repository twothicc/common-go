# grpc client

This package makes it convenient for users to initialize and use a gRPC client.

The gRPC client will have connection pooling.

The gRPC client will be initialized with these chained middlewares:

- `grpc_opentracing` (default): Configured with a jaeger tracer as global OpenTracing tracer. This middleware will extract parent span context from incoming requests, then creates a new span referencing the parent span context. The span context of the new span is then injected into Tag in handler's context.
- `grpc_zap` (default): Configured with common-go logger to log completed gRPC calls. The logger is then populated into the handler's context.

The client can also handle listening for interrupt, terminate, quit os signals to close all connections before shutting down the server.

# Usage

## Initialize and use as a dependency

It is intended for the gRPC client to be initialized as a dependency and passed down for usage.

```
import (
    pool "github.com/twothicc/common-go/grpcclient/pool"
    grpcclient "github.com/twothicc/common-go/grpcclient"
)

...
gRPCClient := grpcclient.NewClient(context.Background(),
    grpcclient.GetDefaultClientConfigs(
        "my_service",
        true,
        // Adds connection pools with specific configs and interceptors
        pool.PoolCreator(pool.GetDefaultConnPoolConfigs("localhost:8080"), nil, nil),
        pool.PoolCreator(pool.GetDefaultConnPoolConfigs("localhost:8081"), nil, nil),
    ),
)

// Listens for interrupt, terminate, quit os signals
go gRPCClient.ListenSignals(context.Background())

doSomethingHandler := DoSomethingHandler(context.Background(), gRPCClient)
...
```

## Call another service

```
resp := &pb.HelloWorldResponse{}

gRPCClient.Call(context.TODO(), "helloWorldServer", "SayHello", &pb.HelloWorldRequest{}, resp)
```

result of the call will be populated into resp.
