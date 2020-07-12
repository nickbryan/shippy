package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/micro/go-micro/v2"

	pb "github.com/nickbryan/shippy/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"
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
	repo         repository
	vesselClient vesselProto.VesselServiceClient
}

func (cs *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := cs.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return fmt.Errorf("unable to find vessel: %w", err)
	}

	req.VesselId = vesselResponse.Vessel.Id

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

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", service.Client())

	service.Init()

	if err := pb.RegisterShippingServiceHandler(service.Server(), &consignmentService{repo: repo, vesselClient: vesselClient}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
