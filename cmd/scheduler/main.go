package main

import (
	"log"
	"net"

	"infra-scheduler/pkg/scheduler"
	pb "infra-scheduler/pkg/schedulerpb"

	"google.golang.org/grpc"
)

func main() {
	// Listen for gRPC on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize the scheduler with Terraform work directory
	sched, err := scheduler.NewScheduler("/opt/terraform/environments/poc")
	if err != nil {
		log.Fatalf("failed to init scheduler: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register our Scheduler service implementation
	pb.RegisterSchedulerServer(grpcServer, sched)

	log.Println("Scheduler gRPC server listening on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
