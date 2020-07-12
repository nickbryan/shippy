package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/micro/go-micro/v2"

	pb "github.com/nickbryan/shippy/shippy-cli-vessel/proto/vessel"
)

type repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type Repository struct {
	vessels []*pb.Vessel
}

func (r *Repository) FindAvailable(s *pb.Specification) (*pb.Vessel, error) {
	for _, v := range r.vessels {
		if s.Capacity <= v.Capacity && s.MaxWeight <= v.MaxWeight {
			return v, nil
		}
	}

	return nil, errors.New("no vessel found by that specification")
}

type vesselService struct {
	repo repository
}

func (vs *vesselService) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	v, err := vs.repo.FindAvailable(req)
	if err != nil {
		return fmt.Errorf("unable to find vessel: %w", err)
	}

	res.Vessel = v

	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{
			Id:        "vessel001",
			Capacity:  500,
			MaxWeight: 200000,
			Name:      "Boaty McBoatface",
		},
	}

	repo := &Repository{vessels}

	service := micro.NewService(micro.Name("shippy.service.vessel"))
	service.Init()

	if err := pb.RegisterVesselServiceHandler(service.Server(), &vesselService{repo}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
