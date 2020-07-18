package main

import (
	"context"
	"log"
	"os"

	pb "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"

	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(micro.Name("shippy.vessel.service"))
	service.Init()

	uri := os.Getenv("DB_HOST")

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Panic(err)
		}
	}()

	vesselCollection := client.Database("shippy").Collection("vessels")
	repository := &MongoRepository{vesselCollection}

	h := &handler{repository: repository}

	if err := h.Create(
		context.Background(),
		&pb.Vessel{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
		&pb.Response{},
	); err != nil {
		log.Panic(err)
	}

	if err := pb.RegisterVesselServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
