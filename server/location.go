package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

type Index struct {
	Index Coordinate `json:"0"`
}

type Coordinate struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

func getLocation(in *pb.RequestCurrent) (string, string) {
	log.Println("'getLocation' function called...")

	url := "http://api.openweathermap.org/geo/1.0/direct?q="
	city := in.City
	token := "&appid=" + os.Getenv("API_KEY")

	url = url + city + token
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching location: %v\n", err)
	}
	defer res.Body.Close()

	log.Println(res)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading location: %v\n", err)
	}

	var coords []Index
	err = json.Unmarshal(body, &coords)
	if err != nil {
		log.Println("Error reading location JSON")
		log.Printf("JSON: %v\n", body)
		log.Printf("Error: %v\n", err)
	}

	return "-123", "44"

}
