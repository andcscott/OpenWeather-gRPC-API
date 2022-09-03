package main

import (
	"context"
	"log"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
)

func doCurrent(c pb.WeatherServiceClient) {

	res, err := c.Current(context.Background(), &pb.RequestCurrent{
		LocationType: pb.LocationType_LOCATION_TYPE_COORDS,
		Units:        pb.Units_UNITS_UNSPECIFIED,
		Location: &pb.OneOfLocation{
			LocationId: &pb.OneOfLocation_Coords{
				Coords: &pb.Coordinates{
					Latitude:  41,
					Longitude: -123,
				}}}})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Payload)
}
