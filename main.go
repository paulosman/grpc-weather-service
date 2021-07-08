package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/paulosman/grpc-weather-service/weather"
	"google.golang.org/grpc"

	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/config"
	"github.com/honeycombio/beeline-go/wrappers/hnygrpc"
)

const (
	baseURL = "http://api.openweathermap.org/data/2.5/weather"
)

var (
	weatherAPIKey    = os.Getenv("WEATHER_API_KEY")
	honeycombAPIKey  = os.Getenv("HONEYCOMB_API_KEY")
	honeycombDataset = os.Getenv("HONEYCOMB_DATASET")
)

type WeatherData struct {
	Humidity    int32   `json:"humidity"`
	Temperature float32 `json:"temp"`
}

type WeatherDescription struct {
	Description string `json:"description"`
}
type WeatherResponse struct {
	Data    WeatherData          `json:"main"`
	Weather []WeatherDescription `json:"weather"`
}

type WeatherServer struct{}

func getWeatherServiceURL(zipCode string, country string) string {
	return fmt.Sprintf(baseURL+"?q=%s,%s&appId=%s",
		url.QueryEscape(zipCode), url.QueryEscape(country), url.QueryEscape(weatherAPIKey))
}

func (w *WeatherServer) GetWeatherByZipCode(ctx context.Context, in *weather.WeatherRequest) (*weather.Weather, error) {
	ctx, span := beeline.StartSpan(ctx, "GetWeatherByZipCode")
	span.AddField("app.zipcode", in.Zip)
	span.AddField("app.countrycode", in.Country)
	defer span.Send()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}
	resp, err := http.Get(getWeatherServiceURL(in.Zip, in.Country))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	response := &WeatherResponse{}
	if err := json.Unmarshal(data, &response); err != nil {
		panic(err)
	}

	return &weather.Weather{
		Temperature: response.Data.Temperature,
		Humidity:    response.Data.Humidity,
		Description: response.Weather[0].Description,
		Hostname:    hostname,
	}, nil
}

func main() {
	beeline.Init(beeline.Config{
		WriteKey:    honeycombAPIKey,
		Dataset:     honeycombDataset,
		ServiceName: "grpc-weather-service",
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
