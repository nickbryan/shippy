module github.com/nickbryan/shippy/shippy-service-consignment

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/nickbryan/shippy/shippy-service-vessel v0.0.0-20200712110101-c04650868374
	google.golang.org/protobuf v1.25.0
)
