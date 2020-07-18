package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	pb "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"
)

type repository interface {
	FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error)
	Create(ctx context.Context, vessel *Vessel) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

type Specification struct {
	Capacity  int32
	MaxWeight int32
}

func MarshalSpecification(spec *pb.Specification) *Specification {
	return &Specification{
		Capacity:  spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

func UnmarshalSpecification(spec *Specification) *pb.Specification {
	return &pb.Specification{
		Capacity:  spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

type Vessel struct {
	ID        string `bson:"id"`
	Capacity  int32  `bson:"capacity"`
	Name      string `bson:"name"`
	Available bool   `bson:"available"`
	OwnerID   string `bson:"owner_id"`
	MaxWeight int32  `bson:"max_weight"`
}

func MarshalVessel(vessel *pb.Vessel) *Vessel {
	return &Vessel{
		ID:        vessel.Id,
		Capacity:  vessel.Capacity,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerID:   vessel.OwnerId,
		MaxWeight: vessel.MaxWeight,
	}
}

func UnmarshalVessel(vessel *Vessel) *pb.Vessel {
	return &pb.Vessel{
		Id:        vessel.ID,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerId:   vessel.OwnerID,
	}
}

func (r *MongoRepository) FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error) {
	log.Println("spec", spec)
	filter := bson.D{
		{"max_weight", bson.D{{"$gte", spec.MaxWeight}}},
		{"capacity", bson.D{{"$gte", spec.Capacity}}},
	}

	v := &Vessel{}
	if err := r.collection.FindOne(ctx, filter).Decode(v); err != nil {
		return nil, fmt.Errorf("unable to find vessel: %w", err)
	}

	return v, nil
}

func (r *MongoRepository) Create(ctx context.Context, vessel *Vessel) error {
	_, err := r.collection.InsertOne(ctx, vessel)
	return err
}
