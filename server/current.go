package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func getCurrent(in *pb.RequestCurrent) string {

	city := in.City
	url := "https://pro.openweathermap.org/data/2.5/weather?q="
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + city + token

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching weather: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading weather response: %v", err)
	}

	return string(body)

}

func (s *server) Current(ctx context.Context, in *pb.RequestCurrent) (*pb.SendCurrent, error) {
	log.Println("'Current' function called...")

	return &pb.SendCurrent{
		Payload: getCurrent(in),
	}, nil
}
