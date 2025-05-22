package provisioner

import (
	"context"
	"fmt"
	"os/exec"

	"infra-scheduler/pkg/models"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// TerraformExecutor wraps terraform-exec
type TerraformExecutor struct {
	tf      *tfexec.Terraform // executor instance
	workDir string            // Terraform directory
}

// NewTerraformExecutor initializes Terraform in the given directory
func NewTerraformExecutor(workDir string) (*TerraformExecutor, error) {
	// Locate terraform binary
	terraformPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, fmt.Errorf("terraform binary not found: %w", err)
	}
	// Create executor bound to workDir
	tf, err := tfexec.NewTerraform(workDir, terraformPath)
	if err != nil {
		return nil, err
	}
	return &TerraformExecutor{tf: tf, workDir: workDir}, nil
}

// ProvisionVM provisions a VM using Terraform
func (e *TerraformExecutor) ProvisionVM(ctx context.Context, req *models.VMRequest) (string, error) {
	// Initialize and install providers/modules
	if err := e.tf.Init(ctx, tfexec.Upgrade(true)); err != nil {
		return "", fmt.Errorf("terraform init failed: %w", err)
	}

	// Build var overrides
	planVars := []tfexec.PlanOption{
		tfexec.Var(fmt.Sprintf("user_id=%s", req.UserID)),
		tfexec.Var(fmt.Sprintf("cores=%d", req.Cores)),
		tfexec.Var(fmt.Sprintf("memory=%d", req.Memory)),
		tfexec.Var(fmt.Sprintf("disk=%d", req.Disk)),
	}

	// Plan
	if _, err := e.tf.Plan(ctx, planVars...); err != nil {
		return "", fmt.Errorf("terraform plan failed: %w", err)
	}

	// Apply
	applyVars := []tfexec.ApplyOption{
		tfexec.Lock(false),
		tfexec.Var(fmt.Sprintf("user_id=%s", req.UserID)),
		tfexec.Var(fmt.Sprintf("cores=%d", req.Cores)),
		tfexec.Var(fmt.Sprintf("memory=%d", req.Memory)),
		tfexec.Var(fmt.Sprintf("disk=%d", req.Disk)),
	}
	if err := e.tf.Apply(ctx, applyVars...); err != nil {
		return "", fmt.Errorf("terraform apply failed: %w", err)
	}

	// Output 'vm_id' from Terraform
	out, err := e.tf.Output(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get outputs: %w", err)
	}
	vmOutput, exists := out["vm_id"]
	if !exists {
		return "", fmt.Errorf("vm_id output not found in terraform state")
	}
	vmID := string(vmOutput.Value)
	if vmID == "" {
		return "", fmt.Errorf("vm_id is empty")
	}
	return vmID, nil
}

// Destroy removes all resources managed by this executor
func (e *TerraformExecutor) Destroy(ctx context.Context) error {
	if err := e.tf.Destroy(ctx); err != nil {
		return fmt.Errorf("terraform destroy failed: %w", err)
	}
	return nil
}
