package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"pilot.go.grpc/blockchainsvc/blockchain"
	"pilot.go.grpc/proto"
)

const (
	CONN_PORT = ":3000"
	CONN_TYPE = "tcp"
)

func main() {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Fatalf("unable to listen on 3000 port: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterBlockchainServer(srv, &Server{
		Blockchain: blockchain.NewBlockchain(),
	})
	srv.Serve(listener)
}

// Server implements proto.BlockchainServer interface
type Server struct {
	Blockchain *blockchain.Blockchain
}

// AddBlock : adds new block to blockchain
func (s *Server) AddBlock(ctx context.Context, in *proto.AddBlockRequest) (*proto.AddBlockResponse, error) {
	block := s.Blockchain.AddBlock(in.Data)
	return &proto.AddBlockResponse{
		Hash: block.Hash,
	}, nil
}

// GetBlockchain : returns blockchain
func (s *Server) GetBlockchain(ctx context.Context, in *proto.GetBlockchainRequest) (*proto.GetBlockchainResponse, error) {
	resp := new(proto.GetBlockchainResponse)
	for _, b := range s.Blockchain.Blocks {
		resp.Blocks = append(resp.Blocks, &proto.Block{
			PrevBlockHash: b.PrevBlockHash,
			Data:          b.Data,
			Hash:          b.Hash,
		})
	}

	return resp, nil
}
