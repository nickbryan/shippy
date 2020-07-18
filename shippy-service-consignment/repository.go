package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	pb "github.com/nickbryan/shippy/shippy-service-consignment/proto/consignment"
	"go.mongodb.org/mongo-driver/mongo"
)

type Consignment struct {
	ID          string     `bson:"id"`
	Weight      int32      `bson:"weight"`
	Description string     `bson:"description"`
	Containers  Containers `bson:"containers"`
	VesselID    string     `bson:"vessel_id"`
}

type Container struct {
	ID         string `bson:"id"`
	CustomerID string `bson:"customer_id"`
	UserID     string `bson:"user_id"`
}

type Containers []*Container

func MarshalContainerCollection(containers []*pb.Container) []*Container {
	collection := make([]*Container, 0)
	for _, container := range containers {
		collection = append(collection, MarshalContainer(container))
	}
	return collection
}

func UnmarshalContainerCollection(containers []*Container) []*pb.Container {
	collection := make([]*pb.Container, 0)
	for _, container := range containers {
		collection = append(collection, UnmarshalContainer(container))
	}
	return collection
}

func UnmarshalConsignmentCollection(consignments []*Consignment) []*pb.Consignment {
	collection := make([]*pb.Consignment, 0)
	for _, consignment := range consignments {
		collection = append(collection, UnmarshalConsignment(consignment))
	}
	return collection
}

func MarshalContainer(container *pb.Container) *Container {
	return &Container{
		ID:         container.Id,
		CustomerID: container.CustomerId,
		UserID:     container.UserId,
	}
}

func UnmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id:         container.ID,
		CustomerId: container.CustomerID,
		UserId:     container.UserID,
	}
}

func MarshalConsignment(consignment *pb.Consignment) *Consignment {
	containers := MarshalContainerCollection(consignment.Containers)
	return &Consignment{
		ID:          consignment.Id,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  containers,
		VesselID:    consignment.VesselId,
	}
}

func UnmarshalConsignment(consignment *Consignment) *pb.Consignment {
	return &pb.Consignment{
		Id:          consignment.ID,
		Description: consignment.Description,
		Weight:      consignment.Weight,
		Containers:  UnmarshalContainerCollection(consignment.Containers),
		VesselId:    consignment.VesselID,
	}
}

type repository interface {
	Create(ctx context.Context, consignment *Consignment) error
	GetAll(ctx context.Context) ([]*Consignment, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func (r *MongoRepository) Create(ctx context.Context, consignment *Consignment) error {
	_, err := r.collection.InsertOne(ctx, consignment)
	if err != nil {
		return fmt.Errorf("unable to insert into mongodb: %w", err)
	}
	return nil
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]*Consignment, error) {
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("unable to get all from mongodb: %w", err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Panic(err)
		}
	}()

	var consignments []*Consignment
	for cur.Next(ctx) {
		var consignment *Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, fmt.Errorf("unable to decode consignment: %w", err)
		}
		consignments = append(consignments, consignment)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return consignments, nil
}
