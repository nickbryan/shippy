package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	pb "github.com/nickbryan/shippy/shippy-service-user/proto/user"
)

type authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

type handler struct {
	repository   repository
	tokenService authable
}

func (h *handler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	result, err := h.repository.Get(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("unable to get user: %w", err)
	}

	res.User = UnmarshalUser(result)

	return nil
}

func (h *handler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	results, err := h.repository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("unable to get all users: %w", err)
	}

	users := UnmarshalUserCollection(results)
	res.Users = users

	return nil
}

func (h *handler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	user, err := h.repository.GetByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("unable to get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fmt.Errorf("password comparison failed: %w", err)
	}

	token, err := h.tokenService.Encode(req)
	if err != nil {
		return fmt.Errorf("unable to encode token: %w", err)
	}

	res.Token = token

	return nil
}

func (h *handler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)
	if err := h.repository.Create(ctx, MarshalUser(req)); err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	// Strip the password back out, so's we're not returning it
	req.Password = ""
	res.User = req

	return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}
