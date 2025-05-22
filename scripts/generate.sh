#!/bin/bash
set -e

# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/hostagent/host.proto pkg/schedulerpb/scheduler.proto

echo "Protobuf files generated successfully." 