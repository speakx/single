package genserver

import (
	"encoding/binary"
	"environment/dump"
	"environment/logger"
	"log"
	"single/app"
	"single/proto/pbsingle"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool"
)

type codecServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *pool.WorkerPool
	sync.Mutex
}

func (cs *codecServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumLoops)
	return
}

func (cs *codecServer) React(c gnet.Conn) (out []byte, action gnet.Action) {
	// if cs.async {
	// 	data := append([]byte{}, c.ReadFrame()...)
	// 	_ = cs.workerPool.Submit(func() {
	// 		c.AsyncWrite(data)
	// 	})
	// 	return
	// }
	out = c.ReadFrame()
	req := &pbsingle.Message{}
	proto.Unmarshal(out, req)

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

	cs.Lock()
	// Store Chunk
	mmapCache, err := app.GetSingleMsgCache(req.FromUid, req.ToUid)
	if nil != err {
		logger.Error("SendMessage transid:", req.Transid, " GetSingleMsgCache err:", err)
		cs.Unlock()
		return
	}

	data, err := proto.Marshal(msgRec)
	if nil != err {
		logger.Error("SendMessage transid:", req.Transid, " proto.Marshal err:", err)
		cs.Unlock()
		return
	}
	n, err := mmapCache.WriteData(0, data, []byte(strconv.FormatUint(msgRec.SrvId, 10)), msgRec)
	if nil != err {
		logger.Error("SendMessage transid:", req.Transid, " mmapCache.WriteData n:", n, " err:", err)
	}

	cs.Unlock()

	reply := &pbsingle.MessageReply{
		ClientId: msgRec.ClientId,
		SrvId:    msgRec.SrvId,
	}
	data, _ = proto.Marshal(reply)
	return data, gnet.None
}

func TestCodecServe(addr string, multicore, async bool, codec gnet.ICodec) {
	var err error
	if codec == nil {
		encoderConfig := gnet.EncoderConfig{
			ByteOrder:                       binary.BigEndian,
			LengthFieldLength:               4,
			LengthAdjustment:                0,
			LengthIncludesLengthFieldLength: false,
		}
		decoderConfig := gnet.DecoderConfig{
			ByteOrder:           binary.BigEndian,
			LengthFieldOffset:   0,
			LengthFieldLength:   4,
			LengthAdjustment:    0,
			InitialBytesToStrip: 4,
		}
		codec = gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
	}
	cs := &codecServer{addr: addr, multicore: multicore, async: async, codec: codec, workerPool: pool.NewWorkerPool()}
	err = gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}
