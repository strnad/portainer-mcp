package client

import (
	"fmt"

	"github.com/portainer/portainer-mcp/pkg/portainer/models"
)

// GetStacks retrieves all regular (non-edge) stacks from the Portainer server.
//
// Returns:
//   - A slice of Stack objects
//   - An error if the operation fails
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

// GetStackFile retrieves the file content of a stack from the Portainer server.
//
// Parameters:
//   - id: The ID of the stack to retrieve
//
// Returns:
//   - The file content of the stack (Compose file)
//   - An error if the operation fails
func (c *PortainerClient) GetStackFile(id int) (string, error) {
	file, err := c.GetRegularStackFile(int64(id))
	if err != nil {
		return "", fmt.Errorf("failed to get stack file: %w", err)
	}

	return file, nil
}

// CreateStack creates a new stack on the Portainer server.
// This function specifically creates a Docker Compose stack.
// Stacks are the equivalent of Edge Stacks in Portainer.
//
// Parameters:
//   - name: The name of the stack
//   - file: The file content of the stack (Compose file)
//   - environmentGroupIds: A slice of environment group IDs to include in the stack
//
// Returns:
//   - The ID of the created stack
//   - An error if the operation fails
func (c *PortainerClient) CreateStack(name, file string, environmentGroupIds []int) (int, error) {
	id, err := c.cli.CreateEdgeStack(name, file, intToInt64Slice(environmentGroupIds))
	if err != nil {
		return 0, fmt.Errorf("failed to create edge stack: %w", err)
	}

	return int(id), nil
}

// UpdateStack updates an existing stack on the Portainer server.
// This function specifically updates a Docker Compose stack.
// Stacks are the equivalent of Edge Stacks in Portainer.
//
// Parameters:
//   - id: The ID of the stack to update
//   - file: The file content of the stack (Compose file)
//   - environmentGroupIds: A slice of environment group IDs to include in the stack
//
// Returns:
//   - An error if the operation fails
func (c *PortainerClient) UpdateStack(id int, file string, environmentGroupIds []int) error {
	err := c.cli.UpdateEdgeStack(int64(id), file, intToInt64Slice(environmentGroupIds))
	if err != nil {
		return fmt.Errorf("failed to update edge stack: %w", err)
	}

	return nil
}

func intToInt64Slice(in []int) []int64 {
	out := make([]int64, len(in))
	for i, v := range in {
		out[i] = int64(v)
	}
	return out
}
