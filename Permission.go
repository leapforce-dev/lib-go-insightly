package insightly

import (
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

// Permission stores Permission from Service
//
type Permission struct {
	PermissionState   string             `json:"PERMISSIONS_STATE"`
	ObjectPermissions []ObjectPermission `json:"OBJECT_PERMISSIONS"`
}

type ObjectPermission struct {
	ObjectType       string `json:"OBJECT_TYPE"`
	PermissionsState string `json:"PERMISSIONS_STATE"`
	CanRead          bool   `json:"CAN_READ"`
	CanCreate        bool   `json:"CAN_CREATE"`
	CanEdit          bool   `json:"CAN_EDIT"`
	CanDelete        bool   `json:"CAN_DELETE"`
}

// GetPermissions returns all permissions
//

func (service *Service) GetPermissions() (*[]Permission, *errortools.Error) {
	permissions := []Permission{}
	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url("Permissions"),
		ResponseModel: &permissions,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &permissions, nil
}
