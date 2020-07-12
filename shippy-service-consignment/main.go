package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/micro/go-micro/v2"

	pb "github.com/nickbryan/shippy/shippy-service-consignment/proto/consignment"
)

type repository interface {
	Create(c *pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

func (r *Repository) Create(c *pb.Consignment) (*pb.Consignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.consignments = append(r.consignments, c)
	return c, nil
}

func (r *Repository) GetAll() []*pb.Consignment {
	return r.consignments
}

type consignmentService struct {
	repo repository
}

func (cs *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	c, err := cs.repo.Create(req)
	if err != nil {
		return fmt.Errorf("unable to create Consignment: %w", err)
	}

	res.Created = true
	res.Consignment = c

	return nil
}

func (cs *consignmentService) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	res.Consignments = cs.repo.GetAll()

	return nil
}

func main() {
	repo := &Repository{}

	service := micro.NewService(micro.Name("shippy.consignmentService.consignment"))
	service.Init()

	if err := pb.RegisterShippingServiceHandler(service.Server(), &consignmentService{repo}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
