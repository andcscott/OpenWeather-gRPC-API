package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func doFiveDay(c pb.WeatherServiceClient) {

	res, err := c.FiveDay(context.Background(), &pb.RequestFiveDay{
		LocationType: pb.LocationType_LOCATION_TYPE_CITY,
		Units:        pb.Units_UNITS_IMPERIAL,
		Location: &pb.OneOfLocation{
			LocationId: &pb.OneOfLocation_City{
				City: "Corvalls",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Payload)
}
