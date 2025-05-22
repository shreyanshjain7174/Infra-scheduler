package models

// VMRequest specifies user VM requirements
type VMRequest struct {
	UserID string `json:"user_id"` // Tenant identifier
	Cores  int    `json:"cores"`   // Number of vCPUs
	Memory int    `json:"memory"`  // RAM in MB
	Disk   int    `json:"disk"`    // Disk size in GB
}

// ProvisionResult returns the created VM identifier
type ProvisionResult struct {
	VMID string `json:"vm_id"` // Unique VM ID from libvirt
}

// ScheduleResult holds the selected host info
type ScheduleResult struct {
	HostID  string `json:"host_id"`
	Address string `json:"address"` // gRPC endpoint
}
