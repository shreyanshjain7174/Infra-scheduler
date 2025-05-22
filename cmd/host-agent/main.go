package main

import (
	"context"
	"log"
	"net"

	pb "infra-scheduler/pkg/hostagent"
	"infra-scheduler/pkg/models"
	"infra-scheduler/pkg/provisioner"

	"google.golang.org/grpc"
)

// server implements the HostAgent gRPC service
type server struct {
	pb.UnimplementedHostAgentServer
	tf *provisioner.TerraformExecutor
}

// ProvisionVM handles incoming gRPC requests to provision a VM
func (s *server) ProvisionVM(ctx context.Context, req *pb.VMRequest) (*pb.VMResponse, error) {
	// Convert gRPC request to model
	modelReq := &models.VMRequest{
		UserID: req.UserId,
		Cores:  int(req.Cores),
		Memory: int(req.Memory),
		Disk:   int(req.Disk),
	}

	// Delegate to our Terraform-based provisioner
	vmID, err := s.tf.ProvisionVM(ctx, modelReq)
	if err != nil {
		log.Printf("ProvisionVM error: %v", err)
		return &pb.VMResponse{Success: false, Error: err.Error()}, nil
	}
	return &pb.VMResponse{Success: true, VmId: vmID}, nil
}

func main() {
	// Listen for gRPC on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize the Terraform executor
	tf, err := provisioner.NewTerraformExecutor("/opt/terraform/environments/poc")
	if err != nil {
		log.Fatalf("failed to init terraform: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register our HostAgent service implementation
	pb.RegisterHostAgentServer(grpcServer, &server{tf: tf})

	log.Println("Host-agent gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
