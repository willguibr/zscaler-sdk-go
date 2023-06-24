package eventlogentryreport

import (
	"github.com/willguibr/zscaler-sdk-go/zia"
)

type Service struct {
	Client *zia.Client
}

func New(c *zia.Client) *Service {
	return &Service{Client: c}
}
