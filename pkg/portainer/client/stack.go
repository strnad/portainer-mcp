package client

import (
	"fmt"

	"github.com/portainer/portainer-mcp/pkg/portainer/models"
)

// GetStacks retrieves all regular (non-edge) stacks from the Portainer server.
func (c *PortainerClient) GetStacks() ([]models.Stack, error) {
	regularStacks, err := c.ListRegularStacks()
	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}

	stacks := make([]models.Stack, len(regularStacks))
	for i, s := range regularStacks {
		stacks[i] = models.ConvertRegularStackToStack(s)
	}

	return stacks, nil
}

// GetStackFile retrieves the compose file content of a stack.
func (c *PortainerClient) GetStackFile(id int) (string, error) {
	file, err := c.GetRegularStackFile(int64(id))
	if err != nil {
		return "", fmt.Errorf("failed to get stack file: %w", err)
	}

	return file, nil
}

// CreateStack creates a new Docker Compose stack on the specified endpoint.
func (c *PortainerClient) CreateStack(name, file string, endpointId int) (int, error) {
	id, err := c.CreateRegularStack(name, file, int64(endpointId))
	if err != nil {
		return 0, fmt.Errorf("failed to create stack: %w", err)
	}

	return int(id), nil
}

// UpdateStack updates an existing stack with new compose file content.
func (c *PortainerClient) UpdateStack(id int, file string, endpointId int, pullImage bool) error {
	err := c.UpdateRegularStack(int64(id), int64(endpointId), file, pullImage)
	if err != nil {
		return fmt.Errorf("failed to update stack: %w", err)
	}

	return nil
}

// StartStack starts a stopped stack.
func (c *PortainerClient) StartStack(id int, endpointId int) error {
	err := c.StartRegularStack(int64(id), int64(endpointId))
	if err != nil {
		return fmt.Errorf("failed to start stack: %w", err)
	}

	return nil
}

// StopStack stops a running stack.
func (c *PortainerClient) StopStack(id int, endpointId int) error {
	err := c.StopRegularStack(int64(id), int64(endpointId))
	if err != nil {
		return fmt.Errorf("failed to stop stack: %w", err)
	}

	return nil
}

// DeleteStack removes a stack.
func (c *PortainerClient) DeleteStack(id int, endpointId int) error {
	err := c.DeleteRegularStack(int64(id), int64(endpointId))
	if err != nil {
		return fmt.Errorf("failed to delete stack: %w", err)
	}

	return nil
}
