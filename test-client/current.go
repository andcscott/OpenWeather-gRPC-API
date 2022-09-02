package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
)

func doCurrent(c pb.WeatherServiceClient) {

	res, err := c.Current(context.Background(), &pb.RequestCurrent{
		LocationType: pb.LocationType_LOCATION_TYPE_CITY,
		Units:        pb.Units_UNITS_METRIC,
		Location: &pb.OneOfLocation{
			LocationId: &pb.OneOfLocation_City{
				City: "Corvallis",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Payload)
}
