package main

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/nickbryan/shippy/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"
)

type handler struct {
	repository   repository
	vesselClient vesselProto.VesselService
}

func (h *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := h.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})

	if err != nil {
		return fmt.Errorf("unable to fetch vessel: %w", err)
	}

	if vesselResponse == nil {
		return errors.New("client returned nil")
	}

	req.VesselId = vesselResponse.Vessel.Id

	if err := h.repository.Create(ctx, MarshalConsignment(req)); err != nil {
		return fmt.Errorf("unable to create consignment: %w", err)
	}

	res.Created = true
	res.Consignment = req

	return nil
}

func (h *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := h.repository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch consignments: %w", err)
	}

	res.Consignments = UnmarshalConsignmentCollection(consignments)

	return nil
}
