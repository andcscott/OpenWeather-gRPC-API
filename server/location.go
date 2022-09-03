package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Coordinates struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

// Receives a RequestLocation message
// Returns a SendLocation message with the Latitude and Longitude
func (s *Server) Location(ctx context.Context, in *pb.RequestLocation) (*pb.SendLocation, error) {
	log.Printf("'Location' called, location: %v\n", in.Location)

	var err error
	var lat, lon float32

	switch in.LocationType {
	case pb.LocationType_LOCATION_TYPE_CITY:
		lat, lon, err = fetchCityCoords(in.Location.GetCity(), s.ApiKey)
	case pb.LocationType_LOCATION_TYPE_ZIP_CODE:
		lat, lon, err = fetchZipCoords(in.Location.GetZipCode(), s.ApiKey)
	default:
		_, err = strconv.Atoi(in.Location.GetZipCode())
		if err != nil {
			lat, lon, err = fetchCityCoords(in.Location.GetCity(), s.ApiKey)
		} else {
			lat, lon, err = fetchZipCoords(in.Location.GetZipCode(), s.ApiKey)
		}
	}
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid location or location type: %s, %s\n",
				in.Location.String(),
				in.LocationType.String()),
		)
	}
	return &pb.SendLocation{
		Latitude:  lat,
		Longitude: lon,
	}, nil
}

// Used internally to fetch precise locations
// Receives a city and the OpenWeather API key
// Returns the latitude and longitude for the given location, otherwise an error
func fetchCityCoords(city string, key string) (float32, float32, error) {
	log.Printf("'fetchCityCoords' called, location: %v\n", city)

	url := "http://api.openweathermap.org/geo/1.0/direct?q="
	token := "&appid=" + key
	url += city + token

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

// Used internally to fetch precise locations
// Receives a zip code and the OpenWeather API key
// Returns the latitude and longitude for the given location, otherwise an error
func fetchZipCoords(zip string, key string) (float32, float32, error) {
	log.Printf("'fetchZipCoords' called, zip code: %v\n", zip)

	url := "https://api.openweathermap.org/geo/1.0/zip?zip="
	token := "&appid=" + key
	url += zip + token

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
