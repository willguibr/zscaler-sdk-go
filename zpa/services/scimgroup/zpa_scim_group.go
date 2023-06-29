package scimgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/zpa/services/common"
)

const (
	userConfig        = "/userconfig/v1/customers/"
	scimGroupEndpoint = "/scimgroup"
	idpId             = "/idpId"
)

type ScimGroup struct {
	CreationTime int64  `json:"creationTime,omitempty"`
	ID           int64  `json:"id,omitempty"`
	IdpGroupID   string `json:"idpGroupId,omitempty"`
	IdpID        int64  `json:"idpId,omitempty"`
	IdpName      string `json:"idpName,omitempty"`
	ModifiedTime int64  `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
}

func (service *Service) Get(scimGroupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint, scimGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(scimName, IdpId string) (*ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint+idpId, IdpId)
	list, resp, err := common.GetAllPagesGeneric[ScimGroup](service.Client, relativeURL, scimName)
	if err != nil {
		return nil, nil, err
	}
	for _, scim := range list {
		if strings.EqualFold(scim.Name, scimName) {
			return &scim, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no scim named '%s' was found", scimName)
}

func (service *Service) GetAllByIdpId(IdpId string) ([]ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint+idpId, IdpId)
	list, resp, err := common.GetAllPagesGeneric[ScimGroup](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
