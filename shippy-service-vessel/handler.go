package main

import (
	"context"
	"fmt"

	pb "github.com/nickbryan/shippy/shippy-service-vessel/proto/vessel"
)

type handler struct {
	repository repository
}

func (h *handler) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	v, err := h.repository.FindAvailable(ctx, MarshalSpecification(req))
	if err != nil {
		return fmt.Errorf("unable to fetch vessel: %w", err)
	}

	res.Vessel = UnmarshalVessel(v)

	return nil
}

func (h *handler) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	if err := h.repository.Create(ctx, MarshalVessel(req)); err != nil {
		return fmt.Errorf("unable to create vessel: %w", err)
	}

	res.Vessel = req
	res.Created = true

	return nil
}
