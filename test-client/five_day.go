package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
)

func doFiveDay(c pb.WeatherServiceClient) {

	res, err := c.FiveDay(context.Background(), &pb.RequestFiveDay{
		LocationType: pb.LocationType_LOCATION_TYPE_COORDS,
		Units:        pb.Units_UNITS_IMPERIAL,
		Location: &pb.OneOfLocation{
			LocationId: &pb.OneOfLocation_Coords{
				Coords: &pb.Coordinates{
					Latitude:  41,
					Longitude: -123,
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Payload)
}
