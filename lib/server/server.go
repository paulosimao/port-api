package server

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/paulosimao/ports-api/lib/db"
	pb "github.com/paulosimao/ports-api/lib/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPortDbServer
}

//PutPort - upsert a port.
func (s *server) PutPort(ctx context.Context, in *pb.PortData) (*pb.PutPortRes, error) {
	log.Printf("Saving port: %#v", in)
	dbp := &db.Port{Code: in.Code, Data: in.Data}
	err := db.PutPort(dbp)
	return &pb.PutPortRes{}, err
}

//GetPorts - returns the ports, implements pb.
func (s *server) GetPorts(req *pb.GetRequest, svc pb.PortDb_GetPortsServer) error {
	log.Printf("Getting ports")
	chports := db.GetPorts()
	ok := true
	var data *db.Port
	for ok {
		data, ok = <-chports
		if !ok {
			return nil
		}
		err := svc.Send(&pb.PortData{Code: data.Code, Data: data.Data})
		if err != nil {
			return err
		}
	}

	return nil

}

//Run - Executes the service as a whole
func Run() error {
	err := db.Init()
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", os.Getenv("ADDR"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPortDbServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
