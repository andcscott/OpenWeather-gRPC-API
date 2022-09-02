package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Coordinates struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

// Receives a gRPC request for Location
// Returns a SendLocation message with the Latitude and Longitude
func (s *Server) Location(ctx context.Context, in *pb.RequestLocation) (*pb.SendLocation, error) {
	log.Println("'Location' function called...")

	lat, lon, err := getLocation(in.Location.String(), s.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("Error: %v\n", err)
	}

	return &pb.SendLocation{
		Latitude:  lat,
		Longitude: lon,
	}, nil
}

// Used internally to fetch precise locations
// Receives the city name and the server's API key
// Returns the latitude and longitude for the given location
func getLocation(location string, key string) (float32, float32, error) {
	log.Println("'getLocation' function called...")

	url := "http://api.openweathermap.org/geo/1.0/direct?q="
	token := "&appid=" + key

	url = url + location + token

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching location: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading location: %v\n", err)
	}

	coords := []Coordinates{}
	err = json.Unmarshal(body, &coords)
	if err != nil {
		log.Printf("Error decoding geolocation JSON: %v\n", err)
	}

	if len(coords) < 1 {
		return 0, 0, status.Error(codes.NotFound, "Location not found")
	}

	return coords[0].Latitude, coords[0].Longitude, nil
}

func getZipLocation(zip string, key string) (float32, float32, error) {
	log.Println("'getZipLocation' function called...")

	url := "https://api.openweathermap.org/geo/1.0/zip?zip="
	token := "&appid=" + key

	url = url + zip + token

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching location: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading location: %v\n", err)
	}

	coords := Coordinates{}
	err = json.Unmarshal(body, &coords)
	if err != nil {
		log.Printf("Error decoding geolocation JSON: %v\n", err)
	}

	if coords.Latitude == 0 && coords.Longitude == 0 {
		return 0, 0, status.Error(codes.NotFound, "Location not found")
	}

	return coords.Latitude, coords.Longitude, nil
}
