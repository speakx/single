package app

import (
	"environment/srvinstance"
	pb "svrdemo/proto"
)

// SimpleGrpcClient simpeclient
type SimpleGrpcClient struct {
	srvinstance.GrpcClient
	pb.SimpleServerClient
}

// Connect connect
func (s *SimpleGrpcClient) Connect(addr string) error {
	err := s.GrpcClient.Connect(addr)
	if nil != err {
		return err
	}

	s.SimpleServerClient = pb.NewSimpleServerClient(s.GetConn())
	return nil
}
