package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
)

func doCurrent(c pb.WeatherServiceClient) {

	res, err := c.Current(context.Background(), &pb.RequestCurrent{
		LocationType: pb.LocationType_LOCATION_TYPE_ZIP_CODE,
		Units:        pb.Units_UNITS_METRIC,
		Location: &pb.OneOfLocation{
			LocationId: &pb.OneOfLocation_ZipCode{
				ZipCode: "97330",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Payload)
}
