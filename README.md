# weather-service

An example gRPC service written in Golang.

### Running

By default, this service binds to port 9000 on 0.0.0.0. In order to use the [OpenWeatherMap](https://home.openweathermap.org/) API, you'll need to sign up and get an API key. You'll also want to set the Honeycomb API key and dataset as environment variables:

```
export HONEYCOMB_API_KEY=my-honeycomb-api-key
export HONEYCOMB_DATASET=my-honeycomb-dataset
export WEATHER_API_KEY=my-openweathermap-api-key

go run main.go
```

You can then run the client:

```
go run client.go
```

The client just loops and spits out whatever response it gets every second or so. This is incredibly boring!