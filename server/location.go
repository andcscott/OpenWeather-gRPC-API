package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	url := "http://api.openweathermap.org/geo/1.0/direct?q="
	city := in.City
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + city + token

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
		log.Printf("Error decoding JSON: %v\n", err)
	}

	return &pb.SendLocation{
		Latitude:  coords[0].Latitude,
		Longitude: coords[0].Longitude,
	}, nil
}

// Used internally to fetch precise locations for Current and Extended
// Receives gRPC requests as an interface
// Returns the latitude (float32) and longitude (float32) for a given location
func getLocation(msg interface{}) (float32, float32) {

	url := "http://api.openweathermap.org/geo/1.0/direct?q="
	city := in.City
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + city + token

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
		log.Printf("Error decoding JSON: %v\n", err)
	}
	return coords[0].Latitude, coords[0].Longitude
}
