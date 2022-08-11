# Grpc Server

This package provides a convenient way to configure and initialize a grpc server.

# Usage

To initialize and run grpc server:

```
registerHelloWorldServiceHandler := func(s *grpc.Server) {
    pb.RegisterHelloWorldServiceServer(s, helloworld.NewHelloWorldServer())
}

serverConfig := grpcserver.GetDefaultServerConfigs("8080", false, registerHelloWorldServiceHandler)

grpcserver.InitAndRunGrpcServer(context.Background(), serverConfig)
```

**Note**: The default `ServerConfigs` is configured for the following:

- Enable Prometheus monitoring
- Time to keepalive ping **1hr**
- Timeout to close connection after keepalive ping **10s**
- Max idle connection time **5mins**

**Note**: The following middlewares will always be enabled:

- grpc_ctxtags: adds a Tag map to context, with data populated from request body
- grpc_zap: Integrates a zap logger into grpc handlers
- grpc_recovery: Turns panic into grpc errors
