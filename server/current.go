package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func (s *Server) Current(ctx context.Context, in *pb.RequestCurrent) (*pb.SendCurrent, error) {
	log.Println("'Current' function called...")

	url := "https://pro.openweathermap.org/data/2.5/weather?q="
	city := in.City
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + city + "&units=imperial" + token

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
