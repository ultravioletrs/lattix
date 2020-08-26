package main

import (
	"fhe"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Config struct {
	Port     string `default:":50051"`
	FilesDir string `default:"/tmp/"`
	Token    string `default:"123"`
}

func main() {
	var c Config
	err := envconfig.Process("server", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Println("server started")
	fhe.RegisterFhesrvServer(s, fhe.NewServer(c.FilesDir, "gob", c.Token))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
