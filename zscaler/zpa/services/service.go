package services

// type Service struct {
// 	Client        *zpa.Client
// 	microTenantID *string
// }

// func New(c *zpa.Client) *Service {
// 	return &Service{Client: c}
// }

// func (service *Service) WithMicroTenant(microTenantID string) *Service {
// 	var mid *string
// 	if microTenantID != "" {
// 		mid_ := microTenantID
// 		mid = &mid_
// 	}
// 	return &Service{
// 		Client:        service.Client,
// 		microTenantID: mid,
// 	}
// }

// func (service *Service) MicroTenantID() *string {
// 	return service.microTenantID
// }

// Helper function to get the full path including mgmtConfig, CustomerID, and endpoint