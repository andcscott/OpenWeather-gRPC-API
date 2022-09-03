package main

import (
	"context"
	"fmt"
	"log"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Receives a gRPC request for an extended forecast
// Returns a SendFiveDay message with the forecast as JSON
func (s *Server) FiveDay(ctx context.Context, in *pb.RequestFiveDay) (*pb.SendFiveDay, error) {
	log.Printf("'FiveDay' called: %v\n", in)

	token := "&appid=" + s.ApiKey
	url, err := s.createFiveDayUrl(in)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf(
				"Invalid location or location type: %v, %v\n",
				in.Location,
				in.LocationType,
			))
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

	return &pb.SendFiveDay{
		Payload: fcst,
	}, nil
}

// Assembles the URL for five day requests to Open Weather
func (s *Server) createFiveDayUrl(in *pb.RequestFiveDay) (string, error) {

	url := "https://api.openweathermap.org/data/2.5/forecast?"
	units := updateUnits(in.Units)
	lat, lon, err := s.fetchLocation(in.LocationType, in.Location)
	if err != nil {
		return "", err
	}

	url += fmt.Sprintf("lat=%f", lat) + fmt.Sprintf("&lon=%f", lon) + units
	return url, nil
}
