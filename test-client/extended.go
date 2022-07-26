package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func doExtended(c pb.WeatherServiceClient) {

	res, err := c.Extended(context.Background(), &pb.RequestExtended{
		City: "Corvallis",
		Days: 7,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Payload)
}
