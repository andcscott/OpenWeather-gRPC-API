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

func getExtended(in *pb.RequestExtended) string {

	city := in.City
	days := "&cnt=" + fmt.Sprint(in.Days)
	url := "https://api.openweathermap.org/data/2.5/forecast/daily?q="
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + city + "&units=imperial" + days + token

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching extended weather: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading extending weather: %v\n", err)
	}
	return string(body)
}

func (s *Server) Extended(ctx context.Context, in *pb.RequestExtended) (*pb.SendExtended, error) {
	log.Println("'Extended' function called...")

	return &pb.SendExtended{
		Payload: getExtended(in),
	}, nil
}
