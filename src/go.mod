module single

go 1.13

replace environment => ../../environment/src

replace idgenerator => ../../idgenerator/src

replace mmapcache => ../../mmapcache/src

replace singledb => ../../singledb/src

replace svrdemo => ../../svrdemo/src

replace github.com/panjf2000/gnet => ../../pkg/github.com/panjf2000/gnet

require (
	environment v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	github.com/panjf2000/gnet v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	google.golang.org/grpc v1.25.1
	mmapcache v0.0.0-00010101000000-000000000000
	singledb v0.0.0-00010101000000-000000000000
)
