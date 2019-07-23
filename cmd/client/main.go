package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	api "mlnagents/grpccompanyserv"
	"time"
)

const address = "localhost:7777"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewCompanyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetCUsersByIDs(ctx, &api.GetCUsersByIDsRequest{CompanyUserIds: []int64{1}})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %+v", r)
}
