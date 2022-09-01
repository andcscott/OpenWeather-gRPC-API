package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func (s *Server) createFiveDayUrl(in *pb.RequestFiveDay) (string, error) {

	var lat, lon float32
	var units string
	var err error
	url := "https://api.openweathermap.org/data/2.5/forecast?"

	switch in.Units {
	case pb.Units_UNITS_IMPERIAL:
		units = "imperial"
	case pb.Units_UNITS_METRIC:
		units = "metric"
	default:
		units = "standard"
	}

	if in.LocationType == pb.LocationType_LOCATION_TYPE_COORDS {
		lat = in.Location.GetCoords().Latitude
		lon = in.Location.GetCoords().Longitude
	} else {
		lat, lon, err = getLocation(in.Location.String(), s.ApiKey)
		if err != nil {
			return "", fmt.Errorf("Error: %v\n", err)
		}
	}

	url = url + fmt.Sprintf("lat=%f", lat) + fmt.Sprintf("&lon=%f", lon) + units
	return url, err

}

// Receives a gRPC request for an extended forecast
// Returns a SendExtended message with the forecast in JSON
func (s *Server) FiveDay(ctx context.Context, in *pb.RequestFiveDay) (*pb.SendFiveDay, error) {
	log.Println("'Extended' function called...")

	url, err := s.createFiveDayUrl(in)
	if err != nil {
		return nil, err
	}
	token := "&appid=" + s.ApiKey

	res, err := http.Get(url + token)
	if err != nil {
		log.Printf("Error fetching extended weather: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading extending weather: %v\n", err)
	}

	return &pb.SendFiveDay{
		Payload: string(body),
	}, nil
}
