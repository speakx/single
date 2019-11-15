package server

import (
	"context"
	"environment/dump"
	"environment/logger"
	"fmt"
	"mmapcache/cache"
	"net"
	pb "svrdemo/proto"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	pb.UnimplementedSimpleServerServer
	mmapCache   *cache.MMapCache
	mmapCacheCh chan proto.Message
}

// NewServer new
func NewServer() *Server {
	s := &Server{
		mmapCacheCh: make(chan proto.Message, 0x1000),
	}
	s.writeMMapCacheLoop()
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
	pb.RegisterSimpleServerServer(srv, s)

	if err := srv.Serve(lis); err != nil {
		logger.Error("failed to serve, err:", err)
	}
	return err
}

// SayHello implements proto.
func (s *Server) SayHello(ctx context.Context, req *pb.SimpleHello) (*pb.SimpleHelloReply, error) {
	// 网络事件处理计数器，dump会通过配置将当前服务的网络事件吞吐量提交给监控服务
	dump.NetEventRecvIncr(0)
	defer dump.NetEventRecvDecr(0)

	// 构建回包 & 处理业务
	reply := &pb.SimpleHelloReply{
		Transid: req.Transid,
		Name:    req.Name,
		Ack:     fmt.Sprintf("%v - hehe", req.Name),
	}

	s.mmapCacheCh <- req
	return reply, nil
}

func (s *Server) writeMMapCacheLoop() {
	go func() {
		for {
			msg, ok := <-s.mmapCacheCh
			if false == ok {
				return
			}

			data, err := proto.Marshal(msg)
			if nil != err {
				logger.Error("proto marshal err:", err)
				continue
			}

			if nil == s.mmapCache {
				s.mmapCache = cache.DefPoolMMapCache.Alloc()
			}
			if n, _ := s.mmapCache.WriteData(0x0, data, []byte(msg.(*pb.SimpleHello).Transid), msg); -1 == n {
				cache.DefPoolMMapCache.Collect(s.mmapCache)
				s.mmapCache = cache.DefPoolMMapCache.Alloc()
				s.mmapCache.WriteData(0x0, data, []byte(msg.(*pb.SimpleHello).Transid), msg)
			}
		}
	}()
}
