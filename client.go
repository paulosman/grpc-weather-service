package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paulosman/grpc-weather-service/weather"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	zipCode := flag.String("zipcode", "70115", "The zip code as a string")
	flag.Parse()

	conn, err := grpc.Dial("grpc.paulosman.me:443", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("dit not connect: %s", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := weather.NewWeatherServiceClient(conn)

	// exit cleanly on SIGTERM or SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		log.Printf("Received", s, "Exiting")
	}()

	for {
		response, err := client.GetWeatherByZipCode(ctx, &weather.ZipCode{Value: *zipCode})
		if err != nil {
			log.Fatalf("Error when calling GetWeatherByZipCode: %s", err)
		}
		log.Printf("Response from server: %+v", response)
		time.Sleep(1000 * time.Millisecond)
	}
}
