package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

// Receives a gRPC request for the current forecast
// Returns a SendCurrent message containing the forecast in JSON
func (s *Server) Current(ctx context.Context, in *pb.RequestCurrent) (*pb.SendCurrent, error) {
	log.Println("'Current' function called...")

	url := "https://pro.openweathermap.org/data/2.5/weather?"
	lat, lon := getLocation(in, in.City)
	units := "&units=imperial"
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + fmt.Sprintf("lat=%f", lat) + fmt.Sprintf("&lon=%f", lon) + units + token

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching weather: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading weather response: %v", err)
	}

	return &pb.SendCurrent{
		Payload: string(body),
	}, nil
}
