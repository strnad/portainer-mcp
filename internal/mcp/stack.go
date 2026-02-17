package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/portainer/portainer-mcp/pkg/toolgen"
)

func (s *PortainerMCPServer) AddStackFeatures() {
	s.addToolIfExists(ToolListStacks, s.HandleGetStacks())
	s.addToolIfExists(ToolGetStackFile, s.HandleGetStackFile())

	if !s.readOnly {
		s.addToolIfExists(ToolCreateStack, s.HandleCreateStack())
		s.addToolIfExists(ToolUpdateStack, s.HandleUpdateStack())
		s.addToolIfExists(ToolStartStack, s.HandleStartStack())
		s.addToolIfExists(ToolStopStack, s.HandleStopStack())
		s.addToolIfExists(ToolDeleteStack, s.HandleDeleteStack())
	}
}

func (s *PortainerMCPServer) HandleGetStacks() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stacks, err := s.cli.GetStacks()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to get stacks", err), nil
		}

		data, err := json.Marshal(stacks)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to marshal stacks", err), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func (s *PortainerMCPServer) HandleGetStackFile() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		stackFile, err := s.cli.GetStackFile(id)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to get stack file", err), nil
		}

		return mcp.NewToolResultText(stackFile), nil
	}
}

func (s *PortainerMCPServer) HandleCreateStack() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		name, err := parser.GetString("name", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid name parameter", err), nil
		}

		file, err := parser.GetString("file", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid file parameter", err), nil
		}

		endpointId, err := parser.GetInt("endpointId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid endpointId parameter", err), nil
		}

		id, err := s.cli.CreateStack(name, file, endpointId)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("error creating stack", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Stack created successfully with ID: %d", id)), nil
	}
}

func (s *PortainerMCPServer) HandleUpdateStack() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		file, err := parser.GetString("file", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid file parameter", err), nil
		}

		endpointId, err := parser.GetInt("endpointId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid endpointId parameter", err), nil
		}

		// pullImage defaults to true if not specified
		pullImage := true
		pullImageParam, pullErr := parser.GetString("pullImage", false)
		if pullErr == nil && pullImageParam == "false" {
			pullImage = false
		}

		err = s.cli.UpdateStack(id, file, endpointId, pullImage)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to update stack", err), nil
		}

		return mcp.NewToolResultText("Stack updated successfully"), nil
	}
}

func (s *PortainerMCPServer) HandleStartStack() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		endpointId, err := parser.GetInt("endpointId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid endpointId parameter", err), nil
		}

		err = s.cli.StartStack(id, endpointId)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to start stack", err), nil
		}

		return mcp.NewToolResultText("Stack started successfully"), nil
	}
}

func (s *PortainerMCPServer) HandleStopStack() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		endpointId, err := parser.GetInt("endpointId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid endpointId parameter", err), nil
		}

		err = s.cli.StopStack(id, endpointId)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to stop stack", err), nil
		}

		return mcp.NewToolResultText("Stack stopped successfully"), nil
	}
}

func (s *PortainerMCPServer) HandleDeleteStack() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		endpointId, err := parser.GetInt("endpointId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid endpointId parameter", err), nil
		}

		err = s.cli.DeleteStack(id, endpointId)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to delete stack", err), nil
		}

		return mcp.NewToolResultText("Stack deleted successfully"), nil
	}
}
