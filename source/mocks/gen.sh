#!/bin/bash

# Logger
mockgen -destination=./mock-logger.go -package=mocks -mock_names=Logger=MockLogger doe/source/logger Logger

# Repository
mockgen -destination=./mock-repository.go -package=mocks -mock_names=Repository=MockRepository doe/source/data Repository

# gRPC
mockgen -destination=./mock-port-service-client.go -package=mocks -mock_names=PortServiceClient=MockPortServiceClient doe/source/models PortServiceClient
mockgen -destination=./mock-port-service-set-client.go -package=mocks -mock_names=PortService_SetClient=MockPortService_SetClient doe/source/models PortService_SetClient
 
# io.Reader
mockgen -destination=./mock-reader.go -package=mocks -mock_names=Reader=MockReader io Reader