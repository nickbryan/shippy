package main

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/nickbryan/shippy/shippy-service-user/proto/user"
)

var (
	key = []byte("mySuperSecretKey")
)

type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

type TokenService struct {
	repository repository
}

func (ts *TokenService) Decode(token string) (*CustomClaims, error) {
	tokenType, err := jwt.ParseWithClaims(string(key), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	if claims, ok := tokenType.Claims.(*CustomClaims); ok && tokenType.Valid {
		return claims, nil
	}

	return nil, errors.New("unable to decode token, reason unknown")
}

func (ts *TokenService) Encode(user *pb.User) (string, error) {
	claims := CustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "shippy.service.user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
