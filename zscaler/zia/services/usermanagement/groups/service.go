package groups

import (
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

type Service struct {
	Client    *zscaler.Client
	sortOrder common.SortOrder
	sortBy    common.SortField
}

func New(c *zscaler.Client) *Service {
	return &Service{
		Client:    c,
		sortOrder: common.ASCSortOrder,
		sortBy:    common.NameSortField,
	}
}

func (service *Service) WithSort(sortBy common.SortField, sortOrder common.SortOrder) *Service {
	c := Service{
		Client:    service.Client,
		sortOrder: service.sortOrder,
		sortBy:    service.sortBy,
	}
	if sortBy == common.IDSortField || sortBy == common.NameSortField || sortBy == common.CreationTimeSortField || sortBy == common.ModifiedTimeSortField {
		c.sortBy = sortBy
	}

	if sortOrder == common.ASCSortOrder || sortOrder == common.DESCSortOrder {
		c.sortOrder = sortOrder
	}
	return &c
}
