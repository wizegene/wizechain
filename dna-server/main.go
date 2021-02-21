package main

import (
	"context"
	pb "github.com/wizegene/wizechain/dna-server/proto"
	"google.golang.org/grpc"
	"net"
)

type DNAServer struct { pb.DnaGeneratorServiceServer }

func NewServer() *DNAServer {
	return new(DNAServer)
}

func (s *DNAServer) GetDNA(ctx context.Context, req *pb.DNARequest) (*pb.DNAResponse, error) {

	return &pb.DNAResponse{}, nil

}

func (s *DNAServer) ValidateDNA(ctx context.Context, req *pb.ValidateDNARequest) (*pb.ValidateDNAResponse, error) {
	return &pb.ValidateDNAResponse{}, nil
}

func NewDNAServer() *DNAServer {
	return &DNAServer{}
}

var host = ":10000"

func main() {

	lis, err := net.Listen("tcp", host)
	if err != nil {
		panic(err);
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterDnaGeneratorServiceServer(grpcServer, NewServer())
	_ = grpcServer.Serve(lis)

}
