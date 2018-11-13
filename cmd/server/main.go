package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/stats"

	pb "github.com/ilya-zz/foo/api"
)

var log = logrus.New().WithField("type", "SERVER")

type server struct {
	storage io.Writer
}

func newServer(path string) *server {
	_, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return &server{
		//		storage: bufio.NewWriterSize(f, 4096),
		storage: ioutil.Discard,
	}
}

func (*server) Hi(m *pb.Hello, r pb.Welcome_HiServer) error {

	for i := 42; i < 52; i++ {
		err := r.Send(&pb.Status{
			Tsid: int64(i),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func calc(data []string) map[string]string {
	m := make(map[string]int)
	for _, s := range data {
		m[s]++
	}
	rt := make(map[string]string)
	for k, v := range m {
		rt[k] = strconv.Itoa(v)
	}
	return rt
}

func calc2(data []string) map[string]string {
	m := make(map[string]string)
	for _, s := range data {
		m[s] = strings.ToUpper(s)
	}
	return m
}

func (s *server) Store(req pb.Welcome_StoreServer) error {
	logrus.Println("Got store request")
	var total int
	for {
		msg, err := req.Recv()
		if err == io.EOF {
			req.SendAndClose(&pb.StoreSummary{
				Written: int64(total),
			})
			logrus.Println("Stored.")
			return nil
		}
		if err != nil {
			logrus.Fatal(err)
		}
		n, err := s.storage.Write(msg.Message)
		if err != nil {
			panic(err)
		}
		total += n
	}
}

func (*server) Translate(req pb.Welcome_TranslateServer) error {
	for {
		text, err := req.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Warn("Server error:  ", err)
			return err
		}
		ws := strings.Fields(text.Text)

		req.Send(&pb.TranslateResult{
			Id:      text.Id,
			Results: calc(ws),
		})
		req.Send(&pb.TranslateResult{
			Id:      text.Id,
			Results: calc2(ws),
		})
	}
}

type stHnd struct {
}

func (*stHnd) HandleRPC(ctx context.Context, st stats.RPCStats) {
	log.Info("HandlerRPC ", st, ctx)
}
func (*stHnd) TagRPC(ctx context.Context, i *stats.RPCTagInfo) context.Context {
	log.Info("TagRPC ", *i)
	return ctx
}

func (*stHnd) TagConn(ctx context.Context, i *stats.ConnTagInfo) context.Context {
	log.Info("TagConn ", *i)
	return ctx
}

func (*stHnd) HandleConn(ctx context.Context, st stats.ConnStats) {
	log.Info("HandleConn ", st, ctx)
}

/*
   // the returned context.
    TagRPC(context.Context, *RPCTagInfo) context.Context
    // HandleRPC processes the RPC stats.
    HandleRPC(context.Context, RPCStats)

    // TagConn can attach some information to the given context.
    // The returned context will be used for stats handling.
    // For conn stats handling, the context used in HandleConn for this
    // connection will be derived from the context returned.
    // For RPC stats handling,
    //  - On server side, the context used in HandleRPC for all RPCs on this
    // connection will be derived from the context returned.
    //  - On client side, the context is not derived from the context returned.
    TagConn(context.Context, *ConnTagInfo) context.Context
    // HandleConn processes the Conn stats.
    HandleConn(context.Context, ConnStats)
*/

func options() []grpc.ServerOption {
	_, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	return []grpc.ServerOption{}
}

const certPath = "/tmp/certs/cert.pem"
const keyPath = "/tmp/certs/key.pem"

func main() {

	p := flag.Int("p", 7777, "GRPC server port")
	tout := flag.Int("t", 120, "Server lifetime")

	flag.Parse()

	//net.JoinHostPort("", strconv.Itoa(*p))
	srvAddr := fmt.Sprintf("192.168.1.13:%d", *p)

	l, err := net.Listen("tcp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	gs := grpc.NewServer(options()...)

	pb.RegisterWelcomeServer(gs, newServer("grpc-storage.db"))

	log.Printf("Starting GRPC server on %s\n", srvAddr)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		log.Println(gs.Serve(l))
		wg.Done()
	}()

	time.Sleep(time.Duration((*tout)) * time.Second)

	log.Warn("Time's up")

	gs.Stop()
	wg.Wait()
}
