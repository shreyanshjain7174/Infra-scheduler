package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"infra-scheduler/pkg/models"
	"infra-scheduler/pkg/provisioner"
	pb "infra-scheduler/pkg/schedulerpb"
)

// Scheduler handles VM scheduling and provisioning
type Scheduler struct {
	pb.UnimplementedSchedulerServer
	tf *provisioner.TerraformExecutor
	mu sync.Mutex
}

// NewScheduler creates a new scheduler instance
func NewScheduler(workDir string) (*Scheduler, error) {
	tf, err := provisioner.NewTerraformExecutor(workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to init terraform: %v", err)
	}

	return &Scheduler{
		tf: tf,
	}, nil
}

// ScheduleVM handles incoming scheduling requests
func (s *Scheduler) ScheduleVM(ctx context.Context, req *pb.ScheduleRequest) (*pb.ScheduleResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert scheduler request to VM request
	vmReq := &models.VMRequest{
		UserID: req.Name, // Using Name as UserID since that's what we have
		Cores:  int(req.Cpu),
		Memory: int(req.Memory),
		Disk:   int(req.Disk),
	}

	// Delegate to Terraform executor
	vmID, err := s.tf.ProvisionVM(ctx, vmReq)
	if err != nil {
		log.Printf("ScheduleVM error: %v", err)
		return &pb.ScheduleResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.ScheduleResponse{
		Success: true,
		VmId:    vmID,
	}, nil
}
