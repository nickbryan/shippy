package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/micro/go-micro/v2"
	pb "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"
)

type repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type Repository struct {
	vessels []*pb.Vessel
}

// FindAvailable - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
func (r *Repository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range r.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("no vessel found by that spec")
}

// Our grpc service handler
type vesselService struct {
	repo repository
}

func (s *vesselService) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	v, err := s.repo.FindAvailable(req)
	if err != nil {
		return fmt.Errorf("unable to find vessel: %w", err)
	}

	// Set the vessel as part of the response message type
	res.Vessel = v
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &Repository{vessels}

	service := micro.NewService(
		micro.Name("shippy.service.vessel"),
	)

	service.Init()

	if err := pb.RegisterVesselServiceHandler(service.Server(), &vesselService{repo}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
