package env

import "os"

func LogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

func ServiceName() string {
	return os.Getenv("SERVICE_NAME")
}

func GrpcAddress() string {
	return os.Getenv("GRPC_ADDRESS")
}

func ClickHouseURI() string {
	return os.Getenv("CLICKHOUSE_URI")
}
