package main

import (
	api "companyserv/grpccompanyserv"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const address = "localhost:7070"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewCompanyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	CreateCompany(c, ctx)
	//	r, err := c.GetCUsersByIDs(ctx, &api.GetCUsersByIDsRequest{CompanyUserIds: []int64{1}})

}

func CreateCompany(c api.CompanyServiceClient, ctx context.Context) {
	r, err := c.CreateCompany(ctx, &api.CreateCompanyRequest{Name: "a3", Title: "aaМащА«2»"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %+v", r)

}
