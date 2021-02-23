package env

import "os"

func GrpcAddress() string {
	return os.Getenv("GRPC_ADDRESS")
}
