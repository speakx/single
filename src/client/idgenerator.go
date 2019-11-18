package client

import (
	"environment/srvinstance"
	"idgenerator/proto/pbidgenerator"
)

// IDGeneratorGrpcClient simpeclient
type IDGeneratorGrpcClient struct {
	srvinstance.GrpcClient
	pbidgenerator.IdGeneratorServerClient
}

// Connect connect
func (s *IDGeneratorGrpcClient) Connect(addr string) error {
	err := s.GrpcClient.Connect(addr)
	if nil != err {
		return err
	}

	s.IdGeneratorServerClient = pbidgenerator.NewIdGeneratorServerClient(s.GetConn())
	return nil
}
