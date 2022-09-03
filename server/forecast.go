package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Receives a gRPC request for the current forecast
// Returns a SendCurrent message containing the forecast in JSON
func (s *Server) Current(ctx context.Context, in *pb.RequestCurrent) (*pb.SendCurrent, error) {
	log.Printf("'Current' called: %v\n", in)

	token := "&appid=" + s.ApiKey
	url, err := s.createCurrentUrl(in)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid location or location type: %s, %s\n",
				in.Location.String(),
				in.LocationType.String()),
		)
	}

	fcst, err := fetchForecast(url + token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf(
				"A server error occurred while fetching the forecast: %v\n",
				err,
			))
	}

	return &pb.SendCurrent{
		Payload: fcst,
	}, nil
}

// Assembles the URL for current weather requests to OpenWeather
func (s *Server) createCurrentUrl(in *pb.RequestCurrent) (string, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?"
	units := updateUnits(in.Units)
	lat, lon, err := s.fetchLocation(in.LocationType, in.Location)
	if err != nil {
		return "", err
	}
	url += fmt.Sprintf("lat=%f", lat) + fmt.Sprintf("&lon=%f", lon) + units
	return url, nil
}

// Creates the units portion of the request to OpenWeater
func updateUnits(units pb.Units) string {
	unitStr := "&units="
	switch units {
	case pb.Units_UNITS_IMPERIAL:
		unitStr += "imperial"
	case pb.Units_UNITS_METRIC:
		unitStr += "metric"
	default:
		unitStr += "standard"
	}
	return unitStr
}

// Obtains a locations exact coordinates for the request to OpenWeather
func (s *Server) fetchLocation(locType pb.LocationType, loc *pb.OneOfLocation) (float32, float32, error) {

	var lat, lon float32
	var err error

	switch locType {
	case pb.LocationType_LOCATION_TYPE_CITY:
		lat, lon, err = fetchCityCoords(loc.GetCity(), s.ApiKey)
	case pb.LocationType_LOCATION_TYPE_ZIP_CODE:
		lat, lon, err = fetchZipCoords(loc.GetZipCode(), s.ApiKey)
	case pb.LocationType_LOCATION_TYPE_COORDS:
		lat = loc.GetCoords().Latitude
		lon = loc.GetCoords().Longitude
	default:
		_, err = strconv.Atoi(loc.GetZipCode())
		if err != nil {
			lat, lon, err = fetchCityCoords(loc.GetCity(), s.ApiKey)
		} else {
			lat, lon, err = fetchZipCoords(loc.GetZipCode(), s.ApiKey)
		}
	}
	if err != nil {
		return 0, 0, err
	}
	return lat, lon, nil
}

// Obtains forecasts from OpenWeather
// Receives the URL as a string
// Returns the JSON response from OpenWeather
func fetchForecast(url string) (string, error) {

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching extended weather: %v\n", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading extending weather: %v\n", err)
		return "", err
	}

	return string(body), nil
}
