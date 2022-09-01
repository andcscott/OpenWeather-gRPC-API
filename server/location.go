package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
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

	return coords[0].Latitude, coords[0].Longitude, err
}
