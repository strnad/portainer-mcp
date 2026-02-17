package client

import (
	"crypto/tls"
	"fmt"
	"net/http"

	goruntime "github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/portainer/client-api-go/v2/client"
	apimodels "github.com/portainer/client-api-go/v2/pkg/models"

	sdkstacks "github.com/portainer/client-api-go/v2/pkg/client/stacks"
)

// PortainerAPIClient defines the interface for the underlying Portainer API client
type PortainerAPIClient interface {
	ListEdgeGroups() ([]*apimodels.EdgegroupsDecoratedEdgeGroup, error)
	CreateEdgeGroup(name string, environmentIds []int64) (int64, error)
	UpdateEdgeGroup(id int64, name *string, environmentIds *[]int64, tagIds *[]int64) error
	ListEdgeStacks() ([]*apimodels.PortainereeEdgeStack, error)
	CreateEdgeStack(name string, file string, environmentGroupIds []int64) (int64, error)
	UpdateEdgeStack(id int64, file string, environmentGroupIds []int64) error
	GetEdgeStackFile(id int64) (string, error)
	ListEndpointGroups() ([]*apimodels.PortainerEndpointGroup, error)
	CreateEndpointGroup(name string, associatedEndpoints []int64) (int64, error)
	UpdateEndpointGroup(id int64, name *string, userAccesses *map[int64]string, teamAccesses *map[int64]string) error
	AddEnvironmentToEndpointGroup(groupId int64, environmentId int64) error
	RemoveEnvironmentFromEndpointGroup(groupId int64, environmentId int64) error
	ListEndpoints() ([]*apimodels.PortainereeEndpoint, error)
	GetEndpoint(id int64) (*apimodels.PortainereeEndpoint, error)
	UpdateEndpoint(id int64, tagIds *[]int64, userAccesses *map[int64]string, teamAccesses *map[int64]string) error
	GetSettings() (*apimodels.PortainereeSettings, error)
	ListTags() ([]*apimodels.PortainerTag, error)
	CreateTag(name string) (int64, error)
	ListTeams() ([]*apimodels.PortainerTeam, error)
	ListTeamMemberships() ([]*apimodels.PortainerTeamMembership, error)
	CreateTeam(name string) (int64, error)
	UpdateTeamName(id int, name string) error
	DeleteTeamMembership(id int) error
	CreateTeamMembership(teamId int, userId int) error
	ListUsers() ([]*apimodels.PortainereeUser, error)
	UpdateUserRole(id int, role int64) error
	GetVersion() (string, error)
	ProxyDockerRequest(environmentId int, opts client.ProxyRequestOptions) (*http.Response, error)
	ProxyKubernetesRequest(environmentId int, opts client.ProxyRequestOptions) (*http.Response, error)
}

// PortainerClient is a wrapper around the Portainer SDK client
// that provides simplified access to Portainer API functionality.
type PortainerClient struct {
	cli       PortainerAPIClient
	stacksSvc sdkstacks.ClientService
	authInfo  goruntime.ClientAuthInfoWriter
}

// ClientOption defines a function that configures a PortainerClient.
type ClientOption func(*clientOptions)

// clientOptions holds configuration options for the PortainerClient.
type clientOptions struct {
	skipTLSVerify bool
}

// WithSkipTLSVerify configures whether to skip TLS certificate verification.
// Setting this to true is not recommended for production environments.
func WithSkipTLSVerify(skip bool) ClientOption {
	return func(o *clientOptions) {
		o.skipTLSVerify = skip
	}
}

// NewPortainerClient creates a new PortainerClient instance with the provided
// server URL and authentication token.
//
// Parameters:
//   - serverURL: The base URL of the Portainer server
//   - token: The authentication token for API access
//   - opts: Optional configuration options for the client
//
// Returns:
//   - A configured PortainerClient ready for API operations
func NewPortainerClient(serverURL string, token string, opts ...ClientOption) *PortainerClient {
	options := clientOptions{
		skipTLSVerify: false, // Default to secure TLS verification
	}

	for _, opt := range opts {
		opt(&options)
	}

	// Create a shared transport for regular stacks API access
	transport := httptransport.New(serverURL, "/api", []string{"https"})
	if options.skipTLSVerify {
		transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	apiKeyAuth := goruntime.ClientAuthInfoWriterFunc(func(r goruntime.ClientRequest, _ strfmt.Registry) error {
		return r.SetHeaderParam("x-api-key", token)
	})
	transport.DefaultAuthentication = apiKeyAuth

	stacksSvc := sdkstacks.New(transport, strfmt.Default)

	sdkCli := client.NewPortainerClient(serverURL, token, client.WithSkipTLSVerify(options.skipTLSVerify))

	return &PortainerClient{
		cli:       sdkCli,
		stacksSvc: stacksSvc,
		authInfo:  apiKeyAuth,
	}
}

// NewPortainerClientWithAPIClient creates a new PortainerClient with a custom API client.
// This is primarily used for testing.
func NewPortainerClientWithAPIClient(cli PortainerAPIClient) *PortainerClient {
	return &PortainerClient{
		cli: cli,
	}
}

// ListRegularStacks lists all regular (non-edge) stacks from the Portainer API.
func (c *PortainerClient) ListRegularStacks() ([]*apimodels.PortainereeStack, error) {
	if c.stacksSvc == nil {
		return nil, fmt.Errorf("stacks service not initialized")
	}

	params := sdkstacks.NewStackListParams()
	ok, _, err := c.stacksSvc.StackList(params, c.authInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}

	if ok == nil {
		return nil, nil
	}

	return ok.Payload, nil
}

// GetRegularStackFile retrieves the compose file content for a regular (non-edge) stack.
func (c *PortainerClient) GetRegularStackFile(id int64) (string, error) {
	if c.stacksSvc == nil {
		return "", fmt.Errorf("stacks service not initialized")
	}

	params := sdkstacks.NewStackFileInspectParams().WithID(id)
	resp, err := c.stacksSvc.StackFileInspect(params, c.authInfo)
	if err != nil {
		return "", fmt.Errorf("failed to get stack file: %w", err)
	}

	if resp.Payload == nil {
		return "", fmt.Errorf("empty stack file response")
	}

	return resp.Payload.StackFileContent, nil
}

// CreateRegularStack creates a new Docker Compose stack via the regular stacks API.
func (c *PortainerClient) CreateRegularStack(name, file string, endpointId int64) (int64, error) {
	if c.stacksSvc == nil {
		return 0, fmt.Errorf("stacks service not initialized")
	}

	body := &apimodels.StacksComposeStackFromFileContentPayload{
		Name:             &name,
		StackFileContent: &file,
	}

	params := sdkstacks.NewStackCreateDockerStandaloneStringParams().
		WithEndpointID(endpointId).
		WithBody(body)

	resp, err := c.stacksSvc.StackCreateDockerStandaloneString(params, c.authInfo)
	if err != nil {
		return 0, fmt.Errorf("failed to create stack: %w", err)
	}

	if resp.Payload == nil {
		return 0, fmt.Errorf("empty create stack response")
	}

	return resp.Payload.ID, nil
}

// UpdateRegularStack updates an existing regular stack with new compose content.
func (c *PortainerClient) UpdateRegularStack(id, endpointId int64, file string, pullImage bool) error {
	if c.stacksSvc == nil {
		return fmt.Errorf("stacks service not initialized")
	}

	body := &apimodels.StacksUpdateStackPayload{
		StackFileContent: file,
		PullImage:        pullImage,
	}

	params := sdkstacks.NewStackUpdateParams().
		WithID(id).
		WithEndpointID(endpointId).
		WithBody(body)

	_, err := c.stacksSvc.StackUpdate(params, c.authInfo)
	if err != nil {
		return fmt.Errorf("failed to update stack: %w", err)
	}

	return nil
}

// StartRegularStack starts a stopped stack.
func (c *PortainerClient) StartRegularStack(id, endpointId int64) error {
	if c.stacksSvc == nil {
		return fmt.Errorf("stacks service not initialized")
	}

	params := sdkstacks.NewStackStartParams().
		WithID(id).
		WithEndpointID(endpointId)

	_, err := c.stacksSvc.StackStart(params, c.authInfo)
	if err != nil {
		return fmt.Errorf("failed to start stack: %w", err)
	}

	return nil
}

// StopRegularStack stops a running stack.
func (c *PortainerClient) StopRegularStack(id, endpointId int64) error {
	if c.stacksSvc == nil {
		return fmt.Errorf("stacks service not initialized")
	}

	params := sdkstacks.NewStackStopParams().
		WithID(id).
		WithEndpointID(endpointId)

	_, err := c.stacksSvc.StackStop(params, c.authInfo)
	if err != nil {
		return fmt.Errorf("failed to stop stack: %w", err)
	}

	return nil
}

// DeleteRegularStack removes a stack.
func (c *PortainerClient) DeleteRegularStack(id, endpointId int64) error {
	if c.stacksSvc == nil {
		return fmt.Errorf("stacks service not initialized")
	}

	params := sdkstacks.NewStackDeleteParams().
		WithID(id).
		WithEndpointID(endpointId)

	_, err := c.stacksSvc.StackDelete(params, c.authInfo)
	if err != nil {
		return fmt.Errorf("failed to delete stack: %w", err)
	}

	return nil
}
