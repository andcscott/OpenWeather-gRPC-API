package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
)

func doLocation(c pb.WeatherServiceClient) {

	res, err := c.Location(context.Background(), &pb.RequestLocation{
		LocationType: pb.LocationType_LOCATION_TYPE_CITY,
		Location: &pb.OneOfLocation{
			LocationId: &pb.OneOfLocation_City{
				City: "Seattle",
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Latitude: %v, Longitude: %v\n", res.Latitude, res.Longitude)
}
