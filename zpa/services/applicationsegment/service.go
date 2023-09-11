package applicationsegment

import (
	zpa "github.com/SecurityGeekIO/zscaler-sdk-go/v2/zpa"
)

type Service struct {
	Client *zpa.Client
}

func New(c *zpa.Client) *Service {
	return &Service{Client: c}
}
