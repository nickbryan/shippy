package main

import (
	"context"
	"log"
	"os"

	"github.com/micro/go-micro/v2"

	pb "github.com/nickbryan/shippy/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	service := micro.NewService(micro.Name("shippy.consignment.service"))
	service.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Panic(err)
		}
	}()

	consignmentCollection := client.Database("shippy").Collection("consignments")
	repository := &MongoRepository{collection: consignmentCollection}

	vesselClient := vesselProto.NewVesselService("shippy.vessel.service", service.Client())

	h := &handler{repository: repository, vesselClient: vesselClient}

	if err := pb.RegisterShippingServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
