package server

import (
	"context"
	"environment/dump"
	"environment/logger"
	"net"
	"single/app"
	"single/proto/pbsingle"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	pbsingle.UnimplementedSingleServerServer
	sync.Mutex
}

// NewServer new
func NewServer() *Server {
	s := &Server{}
	return s
}

// Run server
func (s *Server) Run(addr string) error {
	logger.Info("start listen... addr:", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen, err:", err)
		return err
	}

	srv := grpc.NewServer()
	pbsingle.RegisterSingleServerServer(srv, s)

	if err := srv.Serve(lis); err != nil {
		logger.Error("failed to serve, err:", err)
	}
	return err
}

// SendMessage implements proto.
func (s *Server) SendMessage(ctx context.Context, req *pbsingle.Message) (*pbsingle.MessageReply, error) {
	logger.Debug("SendMessage transid:", req.Transid)

	// 网络事件处理计数器，dump会通过配置将当前服务的网络事件吞吐量提交给监控服务
	dump.NetEventRecvIncr(0)
	defer dump.NetEventRecvDecr(0)

	// Gen Message Record
	msgRec := &pbsingle.MessageRecord{
		ClientId: req.ClientId,
		SrvId:    app.MakeSrvMsgID(req.FromUid, req.ToUid),
		Create:   uint64(time.Now().Unix()),
		FromUid:  req.FromUid,
		ToUid:    req.ToUid,
		Msg:      req.Msg,
		Type:     req.Type,
		Status:   pbsingle.MessageStatus_ORIGIN,
	}

	// Store
	s.Lock()
	mmapCache := app.GetSingleMsgCache(req.FromUid, req.ToUid)
	data, err := proto.Marshal(msgRec)
	if nil != err {
		logger.Error("SendMessage transid:", req.Transid, " proto.Marshal err:", err)
		s.Unlock()
		return nil, err
	}
	n, err := mmapCache.WriteData(0, data, []byte(strconv.FormatUint(msgRec.SrvId, 10)), msgRec)
	if nil != err {
		logger.Error("SendMessage transid:", req.Transid, " mmapCache.WriteData n:", n, " err:", err)
		s.Unlock()
		return nil, err
	}
	s.Unlock()

	return &pbsingle.MessageReply{
		ClientId: msgRec.ClientId,
		SrvId:    msgRec.SrvId,
	}, nil
}

// ModifyMessage implements proto.
func (s *Server) ModifyMessage(ctx context.Context, req *pbsingle.Message) (*pbsingle.MessageReply, error) {
	// 网络事件处理计数器，dump会通过配置将当前服务的网络事件吞吐量提交给监控服务
	dump.NetEventRecvIncr(0)
	defer dump.NetEventRecvDecr(0)
	return nil, nil
}

// RecallMessage implements proto.
func (s *Server) RecallMessage(ctx context.Context, req *pbsingle.Message) (*pbsingle.MessageReply, error) {
	// 网络事件处理计数器，dump会通过配置将当前服务的网络事件吞吐量提交给监控服务
	dump.NetEventRecvIncr(0)
	defer dump.NetEventRecvDecr(0)
	return nil, nil
}
