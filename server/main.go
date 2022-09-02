package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "codeberg.org/andcscott/OpenWeather-gRPC-API/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type Server struct {
	pb.WeatherServiceServer
	ApiKey string
}

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// Read PORT from .env
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Start server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	} else {
		log.Printf("Listening on port %d...\n", port)
	}

	// Initialize gRPC server
	s := grpc.NewServer()
	pb.RegisterWeatherServiceServer(s, &Server{ApiKey: os.Getenv("API_KEY")})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("gRPC Server error: %v\n", err)
	}
}
