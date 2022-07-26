package main

import (
	"log"

	pb "codeberg.org/andcscott/OpenWeatherMap-gRPC-API/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewWeatherServiceClient(conn)

	doCurrent(c)
}
