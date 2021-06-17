package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/paulosman/grpc-weather-service/weather"
	"google.golang.org/grpc"

	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/config"
	"github.com/honeycombio/beeline-go/wrappers/hnygrpc"
)

type WeatherServer struct {
}

func (w *WeatherServer) GetWeatherByZipCode(ctx context.Context, in *weather.ZipCode) (*weather.Weather, error) {
	ctx, span := beeline.StartSpan(ctx, "GetWeatherByZipCode")
	span.AddField("app.zipcode", in.Value)
	defer span.Send()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}
	return &weather.Weather{
		Temperature: 56,
		Humidity:    99,
		Description: "Cloudy",
		Hostname:    hostname,
	}, nil
}

func main() {
	beeline.Init(beeline.Config{
		WriteKey:    "honeycomb-write-key",
		Dataset:     "test-grpc-beeline",
		ServiceName: "grpc-weather-service",
		//APIHost: "http://localhost:8081",
	})
	lis, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", "9000"))
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
	}
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(hnygrpc.UnaryServerInterceptorWithConfig(config.GRPCIncomingConfig{})),
	}
	grpcServer := grpc.NewServer(serverOpts...)
	s := WeatherServer{}
	weather.RegisterWeatherServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
