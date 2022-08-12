package grpcserver

import (
	"path"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

// BasicRequestFieldExtractor - extracts service and method name from incoming requests
func BasicRequestFieldExtractor() grpc_ctxtags.RequestFieldExtractorFunc {
	return func(fullMethod string, req interface{}) map[string]interface{} {
		service := path.Dir(fullMethod)[1:]
		method := path.Base(fullMethod)

		return map[string]interface{}{
			"grpc.service": service,
			"grpc.method":  method,
		}
	}
}
