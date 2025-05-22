// pkg/hostagent/server.go
package hostagent

import (
	"context"
	"infra-scheduler/pkg/models"
	"infra-scheduler/pkg/provisioner"
)

type Server struct{}

func (s *Server) ProvisionVM(ctx context.Context, req *VMRequest) (*VMResponse, error) {
	// Convert VMRequest to models.VMRequest
	modelReq := &models.VMRequest{
		UserID: req.UserId,
		Cores:  int(req.Cores),
		Memory: int(req.Memory),
		Disk:   int(req.Disk),
	}

	tf, err := provisioner.NewTerraformExecutor("terraform/environments/poc")
	if err != nil {
		return &VMResponse{Success: false, Error: err.Error()}, nil
	}

	vmID, err := tf.ProvisionVM(ctx, modelReq)
	if err != nil {
		return &VMResponse{Success: false, Error: err.Error()}, nil
	}
	return &VMResponse{Success: true, VmId: vmID}, nil
}
