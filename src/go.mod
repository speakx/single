module single

go 1.13

replace environment => ../../environment/src

replace idgenerator => ../../idgenerator/src

replace mmapcache => ../../mmapcache/src

require (
	environment v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	google.golang.org/grpc v1.25.1
	mmapcache v0.0.0-00010101000000-000000000000
)
