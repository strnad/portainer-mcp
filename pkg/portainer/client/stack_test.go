package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: GetStacks, GetStackFile, CreateStack, UpdateStack, StartStack, StopStack, DeleteStack
// now use the SDK stacks service (stacksSvc) directly rather than the edge stacks API (cli).
// Unit testing these requires mocking the SDK transport, which is covered by integration tests.
// Handler-level tests in internal/mcp/stack_test.go cover the interface contract.

func TestClientStackMethodsRequireStacksSvc(t *testing.T) {
	// Client without stacksSvc should return errors for all stack operations
	client := &PortainerClient{cli: nil, stacksSvc: nil}

	t.Run("GetStacks without stacksSvc", func(t *testing.T) {
		_, err := client.GetStacks()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})

	t.Run("GetStackFile without stacksSvc", func(t *testing.T) {
		_, err := client.GetStackFile(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})

	t.Run("CreateStack without stacksSvc", func(t *testing.T) {
		_, err := client.CreateStack("test", "file", 8)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})

	t.Run("UpdateStack without stacksSvc", func(t *testing.T) {
		err := client.UpdateStack(1, "file", 8, true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})

	t.Run("StartStack without stacksSvc", func(t *testing.T) {
		err := client.StartStack(1, 8)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})

	t.Run("StopStack without stacksSvc", func(t *testing.T) {
		err := client.StopStack(1, 8)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})

	t.Run("DeleteStack without stacksSvc", func(t *testing.T) {
		err := client.DeleteStack(1, 8)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stacks service not initialized")
	})
}
