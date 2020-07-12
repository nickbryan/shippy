module github.com/nickbryan/shippy/shippy-cli-consignment

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/micro/go-micro/v2 v2.9.1
	github.com/nickbryan/shippy/shippy-service-consignment v0.0.0-20200712102434-0d36293d7492
)
