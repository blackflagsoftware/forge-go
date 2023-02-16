package main

import (
	"net" // --- replace migration header os once text - do not remove ---

	"github.com/blackflagsoftware/forge-go/base/config"
	m "github.com/blackflagsoftware/forge-go/base/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// --- replace migration header once text - do not remove ---
	// --- replace grpc import once - do not remove ---
	// --- replace grpc import - do not remove ---
)

func main() { // --- replace migration once text - do not remove ---

	tcpListener, err := net.Listen("tcp", ":"+config.GrpcPort)
	if err != nil {
		m.Default.Panic("Unable to start GRPC port:", err)
	}
	defer tcpListener.Close()
	s := grpc.NewServer()

	// --- replace grpc text - do not remove ---

	reflection.Register(s)
	m.Default.Printf("Starting GRPC server on port: %s...\n", config.GrpcPort)
	s.Serve(tcpListener)
}
