package main

import (
	"context"
	pb "github.com/ahmadkarlam/shippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*consignment.Consignment) (*consignment.Consignment, error)
}

type Repository struct {
	mu sync.RWMutex
	consignments []*consignment.Consignment
}

func (repo *Repository) Create(consignment *consignment.Consignment) (*consignment.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *consignment.Consignment) (*consignment.Response, error) {
	consignment, err := s.repo.Create(req)

	if err != nil {
		return nil, err
	}

	return &consignment.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	consignment.RegisterShippingServiceServer(s, &service{repo})

	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
