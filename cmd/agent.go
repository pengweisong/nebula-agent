package main

import (
	"flag"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/vesoft-inc/nebula-agent/internal/clients"
	_ "github.com/vesoft-inc/nebula-agent/internal/log"
	"github.com/vesoft-inc/nebula-agent/internal/server"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

var (
	GitInfoSHA string
)

var (
	agent = flag.String("agent", "auto", "The agent server address")
	meta  = flag.String("meta", "", "The nebula metad service address, any metad address will be ok")
	hbs   = flag.Int("hbs", 60, "Agent heartbeat interval to nebula meta, in seconds")
	debug = flag.Bool("debug", false, "Open debug will output more detail info")
)

func main() {
	flag.Parse()
	log.WithField("version", GitInfoSHA).Info("Start agent server...")

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	lis, err := net.Listen("tcp", *agent)
	if err != nil {
		log.WithError(err).Fatalf("Failed to listen: %v.", *agent)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	metaCfg, err := clients.NewMetaConfig(*agent, *meta, GitInfoSHA, *hbs)
	if err != nil {
		log.WithError(err).Fatalf("Failed to create meta config.")
	}
	agentServer, err := server.NewAgent(metaCfg)
	if err != nil {
		log.WithError(err).Fatalf("Failed to create agent server.")
	}

	pb.RegisterAgentServiceServer(grpcServer, agentServer)
	pb.RegisterStorageServiceServer(grpcServer, server.NewStorage())
	grpcServer.Serve(lis)
}
