package client

import (
	"environment/srvinstance"
	"singledb/proto/pbsingledb"
)

// SingleDBGrpcClient simpeclient
type SingleDBGrpcClient struct {
	srvinstance.GrpcClient
	pbsingledb.SingleDBServerClient
}

// Connect connect
func (s *SingleDBGrpcClient) Connect(addr string) error {
	err := s.GrpcClient.Connect(addr)
	if nil != err {
		return err
	}

	s.SingleDBServerClient = pbsingledb.NewSingleDBServerClient(s.GetConn())
	return nil
}
