all:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative weather/weather.proto
amd64:
	GOOS=linux GOARCH=amd64 go build -o weather-service main.go
	aws s3 cp weather-service s3://grpc-example-app/