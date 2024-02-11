package policysetcontroller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfigV1 = "/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2 = "/mgmtconfig/v2/admin/customers/"
)

type PolicySet struct {
	CreationTime    string       `json:"creationTime,omitempty"`
	Description     string       `json:"description,omitempty"`
	Enabled         bool         `json:"enabled"`
	ID              string       `json:"id,omitempty"`
	ModifiedBy      string       `json:"modifiedBy,omitempty"`
	ModifiedTime    string       `json:"modifiedTime,omitempty"`
	Name            string       `json:"name,omitempty"`
	Sorted          bool         `json:"sorted"`
	PolicyType      string       `json:"policyType,omitempty"`
	MicroTenantID   string       `json:"microtenantId,omitempty"`
	MicroTenantName string       `json:"microtenantName,omitempty"`
	Rules           []PolicyRule `json:"rules"`
}

type PolicyRule struct {
	Action                   string               `json:"action,omitempty"`
	ActionID                 string               `json:"actionId,omitempty"`
	BypassDefaultRule        bool                 `json:"bypassDefaultRule"`
	CreationTime             string               `json:"creationTime,omitempty"`
	CustomMsg                string               `json:"customMsg,omitempty"`
	DefaultRule              bool                 `json:"defaultRule,omitempty"`
	DefaultRuleName          string               `json:"defaultRuleName,omitempty"`
	Description              string               `json:"description,omitempty"`
	ID                       string               `json:"id,omitempty"`
	IsolationDefaultRule     bool                 `json:"isolationDefaultRule"`
	ModifiedBy               string               `json:"modifiedBy,omitempty"`
	ModifiedTime             string               `json:"modifiedTime,omitempty"`
	Name                     string               `json:"name,omitempty"`
	Operator                 string               `json:"operator,omitempty"`
	PolicySetID              string               `json:"policySetId"`
	PolicyType               string               `json:"policyType,omitempty"`
	Priority                 string               `json:"priority,omitempty"`
	ReauthDefaultRule        bool                 `json:"reauthDefaultRule"`
	ReauthIdleTimeout        string               `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string               `json:"reauthTimeout,omitempty"`
	RuleOrder                string               `json:"ruleOrder"`
	LSSDefaultRule           bool                 `json:"lssDefaultRule"`
	ZpnCbiProfileID          string               `json:"zpnCbiProfileId,omitempty"`
	ZpnIsolationProfileID    string               `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID   string               `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string               `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string               `json:"microtenantId,omitempty"`
	MicroTenantName          string               `json:"microtenantName,omitempty"`
	Conditions               []Conditions         `json:"conditions"`
	AppServerGroups          []AppServerGroups    `json:"appServerGroups"`
	AppConnectorGroups       []AppConnectorGroups `json:"appConnectorGroups"`
}

type Conditions struct {
	CreationTime  string     `json:"creationTime,omitempty"`
	ID            string     `json:"id,omitempty"`
	ModifiedBy    string     `json:"modifiedBy,omitempty"`
	ModifiedTime  string     `json:"modifiedTime,omitempty"`
	Negated       bool       `json:"negated"`
	Operands      []Operands `json:"operands"`
	Operator      string     `json:"operator,omitempty"`
	MicroTenantID string     `json:"microtenantId,omitempty"`
}

type Operands struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	CreationTime  string `json:"creationTime,omitempty"`
	ModifiedBy    string `json:"modifiedBy,omitempty"`
	ModifiedTime  string `json:"modifiedTime,omitempty"`
	IdpID         string `json:"idpId,omitempty"`
	LHS           string `json:"lhs,omitempty"`
	RHS           string `json:"rhs,omitempty"`
	ObjectType    string `json:"objectType,omitempty"`
	MicroTenantID string `json:"microtenantId,omitempty"`
}

type AppServerGroups struct {
	ID string `json:"id,omitempty"`
}

type AppConnectorGroups struct {
	ID string `json:"id,omitempty"`
}

type Count struct {
	Count string `json:"count"`
}

func (service *Service) GetByPolicyType(policyType string) (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfigV1 + service.Client.Config.CustomerID + "/policySet/policyType/" + policyType)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// GET --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule/{ruleId}
func (service *Service) GetPolicyRule(policySetID, ruleId string) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	url := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("GET", url, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// POST --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule
func (service *Service) CreateRuleV1(rule *PolicyRule) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	path := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/%s/rule", rule.PolicySetID)
	resp, err := service.Client.NewRequestDo("POST", path, common.Filter{MicroTenantID: service.microTenantID}, &rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) UpdateRuleV1(policySetID, ruleId string, policySetRule *PolicyRule) (*http.Response, error) {
	if policySetRule != nil && len(policySetRule.Conditions) == 0 {
		policySetRule.Conditions = []Conditions{}
	} else {
		for i, condtion := range policySetRule.Conditions {
			if len(condtion.Operands) == 0 {
				policySetRule.Conditions[i].Operands = []Operands{}
			} else {
				for i, operand := range condtion.Operands {
					if operand.Name != "" {
						condtion.Operands[i].Name = ""
					}
				}
			}
		}
	}
	path := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, policySetRule, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DELETE --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) Delete(policySetID, ruleId string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetByNameAndType(policyType, ruleName string) (*PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](service.Client, relativeURL, common.Filter{Search: ruleName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}

	for _, p := range list {
		if strings.EqualFold(ruleName, p.Name) {
			return &p, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no policy rule named :%s found", ruleName)
}

func (service *Service) GetByNameAndTypes(policyTypes []string, ruleName string) (p *PolicyRule, resp *http.Response, err error) {
	for _, policyType := range policyTypes {
		p, resp, err = service.GetByNameAndType(policyType, ruleName)
		if err != nil {
			continue
		} else {
			return
		}
	}
	return
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySetId}/rule/{ruleId}/reorder/{newOrder}
func (service *Service) Reorder(policySetID, ruleId string, order int) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/%s/rule/%s/reorder/%d", policySetID, ruleId, order)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySet}/reorder
// ruleIdOrders is a map[ruleID]Order
func (service *Service) BulkReorder(policySetType string, ruleIdToOrder map[string]int) (*http.Response, error) {
	policySet, resp, err := service.GetByPolicyType(policySetType)
	if err != nil {
		return resp, err
	}
	all, resp, err := service.GetAllByType(policySetType)
	if err != nil {
		return resp, err
	}
	sort.SliceStable(all, func(i, j int) bool {
		ruleIDi := all[i].ID
		ruleIDj := all[j].ID

		// Check if ruleIDi and ruleIDj exist in the ruleIdToOrder map
		orderi, existsi := ruleIdToOrder[ruleIDi]
		orderj, existsj := ruleIdToOrder[ruleIDj]

		// If both rules exist in the map, compare their orders
		if existsi && existsj {
			return orderi <= orderj
		}

		// If only one of the rules exists in the map, prioritize it
		if existsi {
			return true
		} else if existsj {
			return false
		}

		// If neither rule exists in the map, maintain their relative order
		return i <= j
	})
	// Construct the URL path
	path := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/%s/reorder", policySet.ID)
	ruleIdsOrdered := []string{}
	for _, r := range all {
		ruleIdsOrdered = append(ruleIdsOrdered, r.ID)
	}

	// Create a new PUT request
	resp, err = service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, ruleIdsOrdered, nil)
	if err != nil {
		return nil, err
	}

	// Check for non-2xx status code and log response body for debugging
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close() // Ensure the body is always closed
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// Handle the error of reading the body (optional)
			log.Printf("Error reading response body: %s\n", err.Error())
		}
		log.Printf("Error response from API: %s\n", string(bodyBytes))
		return resp, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func (service *Service) RulesCount() (int, *http.Response, error) {
	v := new(Count)
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/rules/policyType/GLOBAL_POLICY/count", service.Client.Config.CustomerID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return 0, nil, err
	}
	count, err := strconv.Atoi(v.Count)
	return count, resp, err
}

func (service *Service) GetAllByType(policyType string) ([]PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
