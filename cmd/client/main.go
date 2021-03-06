package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	pb "github.com/ilya-zz/grpc-probe/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var log = logrus.New().WithField("type", "CLIENT")

func defaultCallOptions() []grpc.CallOption {
	return []grpc.CallOption{
		grpc.Header(&metadata.MD{
			"foo": []string{"bar"},
		}),
	}
}

const certPath = "/tmp/certs/cert.pem"

func insecure() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
	}
}

func grpcOpts() []grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile(certPath, "www.fuck.off")
	if err != nil {
		log.Fatal(err)
	}

	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(defaultCallOptions()...),
	}
}

var (
	url   = flag.String("url", "", "Server url")
	count = flag.Int("count", 10000, "Count")
	bs    = flag.Int("bs", 1024, "Block size")
)

func main() {
	flag.Parse()

	if *url == "" {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "URL required\n")
		os.Exit(1)
	}

	c, err := grpc.Dial(fmt.Sprintf(*url), insecure()...)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	client := pb.NewWelcomeClient(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := client.Store(ctx)
	if err != nil {
		panic(err)
	}

	buf := []byte(strings.Repeat("Z", *bs))
	t0 := time.Now()

	log.Printf("Send %d packages (%d each) to %s\n", *count, humanize.Bytes(uint64(*bs)), *url)

	send := 0
	for i := 0; i < *count; i++ {
		st.Send(&pb.RecordMessage{
			Message: buf,
		})
		send += len(buf)
	}
	r, err := st.CloseAndRecv()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	fmt.Printf("Sent %s (%s) bytes in %f secs\n",
		humanize.Bytes(uint64(send)),
		humanize.Bytes(uint64(r.Written)),
		time.Since(t0).Seconds())
}

func main2() {
	port := flag.Int("p", 7777, "local port to connect")
	flag.Parse()

	c, err := grpc.Dial(fmt.Sprintf(":%d", *port), grpcOpts()...)
	if err != nil {
		logrus.Printf("%v, retry.. ", err)
		time.Sleep(200 * time.Millisecond)
		c, err = grpc.Dial(fmt.Sprintf(":%d", *port), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
	}

	defer c.Close()

	client := pb.NewWelcomeClient(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Translate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			r, err := stream.Recv()
			if err != nil {
				logrus.Warn("Client error: ", err)
				if err == io.EOF {
					return
				}
				return
			}
			logrus.Printf("%d -> %v", r.Id, r.Results)
		}
	}()

	var id int64

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Err() != nil {
			return
		}
		stream.Send(&pb.ToTranslate{
			Id:   id,
			Text: s.Text(),
		})
		id++
	}

}
