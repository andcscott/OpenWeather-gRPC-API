package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func doCurrent(c pb.WeatherServiceClient) {

	res, err := c.Current(context.Background(), &pb.RequestCurrent{
		City: "Corvallis",
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Payload)
}
