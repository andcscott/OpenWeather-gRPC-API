package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) createFiveDayUrl(in *pb.RequestFiveDay) (string, error) {

	var lat, lon float32
	var err error
	url := "https://api.openweathermap.org/data/2.5/forecast?"
	units := "&units="

	switch in.Units {
	case pb.Units_UNITS_IMPERIAL:
		units += "imperial"
	case pb.Units_UNITS_METRIC:
		units += "metric"
	default:
		units += "standard"
	}

	switch in.LocationType {
	case pb.LocationType_LOCATION_TYPE_CITY:
		lat, lon, err = getLocation(in.Location.GetCity(), s.ApiKey)
		if err != nil {
			return "", err
		}
	case pb.LocationType_LOCATION_TYPE_ZIP_CODE:
		lat, lon, err = getZipLocation(in.Location.GetZipCode(), s.ApiKey)
		if err != nil {
			return "", err
		}
	case pb.LocationType_LOCATION_TYPE_COORDS:
		lat = in.Location.GetCoords().Latitude
		lon = in.Location.GetCoords().Longitude
	default:
		lat, lon, err = getLocation(in.Location.String(), s.ApiKey)
		if err != nil {
			return "", err
		}
	}

	url = url + fmt.Sprintf("lat=%f", lat) + fmt.Sprintf("&lon=%f", lon) + units
	return url, nil

}

// Receives a gRPC request for an extended forecast
// Returns a SendExtended message with the forecast in JSON
func (s *Server) FiveDay(ctx context.Context, in *pb.RequestFiveDay) (*pb.SendFiveDay, error) {
	log.Println("'FiveDay' function called...")

	url, err := s.createFiveDayUrl(in)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid location: %s", in.Location.String()),
		)
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
