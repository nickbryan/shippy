package main

import (
	"log"

	"github.com/micro/go-micro/v2"
	pb "github.com/nickbryan/shippy/shippy-service-user/proto/user"
)

const schema = `
	create table if not exists users (
		id varchar(36) not null,
		name varchar(125) not null,
		email varchar(225) not null unique,
		password varchar(225) not null,
		company varchar(125),
		primary key (id)
	);
`

func main() {
	db, err := NewConnection()
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Panic(err)
		}
	}()

	db.MustExec(schema)

	repo := &PostgresRepository{db: db}
	tokenService := &TokenService{repo}

	service := micro.NewService(micro.Name("shippy.service.user"))
	service.Init()

	if err := pb.RegisterUserServiceHandler(service.Server(), &handler{repo, tokenService}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
