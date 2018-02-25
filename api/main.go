package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"pilot.go.grpc/proto"
)

var client proto.BlockchainClient

func main() {
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	client = proto.NewBlockchainClient(conn)

	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	r.Handle("/health", HealthHandler).Methods("GET")
	r.Handle("/addblock", AddBlockHandler).Methods("GET")
	r.Handle("/blocks", GetBlockchainHandler).Methods("GET")

	// Our application will run on port 3000. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

var HealthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

var AddBlockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	block, addErr := client.AddBlock(context.Background(), &proto.AddBlockRequest{
		Data: time.Now().String(),
	})
	if addErr != nil {
		log.Fatalf("unable to add block: %v", addErr)
	}
	log.Printf("new block hash: %s\n", block.Hash)

	response, _ := json.Marshal(block)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})

var GetBlockchainHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	blockchain, getErr := client.GetBlockchain(context.Background(), &proto.GetBlockchainRequest{})
	if getErr != nil {
		log.Fatalf("unable to get blockchain: %v", getErr)
	}

	log.Println("blocks:")
	for _, b := range blockchain.Blocks {
		log.Printf("hash %s, prev hash: %s, data: %s\n", b.Hash, b.PrevBlockHash, b.Data)
	}

	response, _ := json.Marshal(blockchain)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})
