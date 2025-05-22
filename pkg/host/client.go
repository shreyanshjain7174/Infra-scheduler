package host

import (
	"context"
	"fmt"
	"time"

	pb "infra-scheduler/pkg/hostagent"

	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	cli  pb.HostAgentClient
}

// NewClient dials the host-agent at the given address
func NewClient(addr string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to dial host-agent: %w", err)
	}
	return &Client{conn: conn, cli: pb.NewHostAgentClient(conn)}, nil
}

// ProvisionVM invokes gRPC to provision the VM on the host
func (c *Client) ProvisionVM(ctx context.Context, req *pb.VMRequest) (*pb.VMResponse, error) {
	return c.cli.ProvisionVM(ctx, req)
}
