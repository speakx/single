package app

import (
	"environment/srvinstance"
	"single/proto/pbsingle"
)

// SimpleGrpcClient simpeclient
type SimpleGrpcClient struct {
	srvinstance.GrpcClient
	pbsingle.SimpleServerClient
}

// Connect connect
func (s *SimpleGrpcClient) Connect(addr string) error {
	err := s.GrpcClient.Connect(addr)
	if nil != err {
		return err
	}

	s.SimpleServerClient = pbsingle.NewSimpleServerClient(s.GetConn())
	return nil
}
