package models

import (
	"time"

	apimodels "github.com/portainer/client-api-go/v2/pkg/models"
	"github.com/portainer/portainer-mcp/pkg/portainer/utils"
)

type Stack struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Status              string `json:"status"`
	CreatedAt           string `json:"created_at"`
	EndpointID          int    `json:"endpoint_id,omitempty"`
	EnvironmentGroupIds []int  `json:"group_ids,omitempty"`
}

func ConvertEdgeStackToStack(rawEdgeStack *apimodels.PortainereeEdgeStack) Stack {
	createdAt := time.Unix(rawEdgeStack.CreationDate, 0).Format(time.RFC3339)

	return Stack{
		ID:                  int(rawEdgeStack.ID),
		Name:                rawEdgeStack.Name,
		CreatedAt:           createdAt,
		EnvironmentGroupIds: utils.Int64ToIntSlice(rawEdgeStack.EdgeGroups),
	}
}

func ConvertRegularStackToStack(rawStack *apimodels.PortainereeStack) Stack {
	createdAt := time.Unix(rawStack.CreationDate, 0).Format(time.RFC3339)

	status := "inactive"
	if rawStack.Status == 1 {
		status = "active"
	}

	return Stack{
		ID:         int(rawStack.ID),
		Name:       rawStack.Name,
		Status:     status,
		CreatedAt:  createdAt,
		EndpointID: int(rawStack.EndpointID),
	}
}
