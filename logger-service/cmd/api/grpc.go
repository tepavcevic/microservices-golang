package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tepavcevic/microservices-golang/logger/data"
	"github.com/tepavcevic/microservices-golang/logger/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	if err := l.Models.LogEntry.Insert(logEntry); err != nil {
		res := logs.LogResponse{
			Result: "failed",
		}
		return &res, err
	}

	res := logs.LogResponse{Result: "logged!"}

	return &res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen to gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC server started on port :%s", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen to gRPC: %v", err)
	}
}
