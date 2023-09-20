// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policydefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policysetdefinitions"
)

func getPolicyDefinitionByDisplayName(ctx context.Context, client *policydefinitions.PolicyDefinitionsClient, subscriptionID, displayName, managementGroupName string,
	builtInOnly bool) (policydefinitions.PolicyDefinition, error) {
	var policyDefinitions []policydefinitions.PolicyDefinition
	var err error

	if managementGroupName != "" {
		id := commonids.NewManagementGroupID(managementGroupName)
		var result policydefinitions.ListByManagementGroupCompleteResult
		result, err = client.ListByManagementGroupComplete(ctx, id, policydefinitions.ListByManagementGroupOperationOptions{})
		policyDefinitions = result.Items
	} else {
		if builtInOnly {
			var result policydefinitions.ListBuiltInCompleteResult
			result, err = client.ListBuiltInComplete(ctx, policydefinitions.ListBuiltInOperationOptions{})
			policyDefinitions = result.Items
		} else {
			var result policydefinitions.ListCompleteResult
			result, err = client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), policydefinitions.ListOperationOptions{})
			policyDefinitions = result.Items
		}
	}
	if err != nil {
		return policydefinitions.PolicyDefinition{}, fmt.Errorf("loading Policy Definition List: %+v", err)
	}

	// var results []policy.Definition
	// for policyDefinitions.NotDone() {
	// 	def := policyDefinitions.Value()
	// 	if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
	// 		results = append(results, def)
	// 	}
	//
	// 	if err := policyDefinitions.NextWithContext(ctx); err != nil {
	// 		return policy.Definition{}, fmt.Errorf("loading Policy Definition List: %s", err)
	// 	}
	// }

	// we found none
	if len(policyDefinitions) == 0 {
		return policydefinitions.PolicyDefinition{}, fmt.Errorf("loading Policy Definition List: could not find policy '%s'. has the policies name changed? list available with `az policy definition list`", displayName)
	}

	// we found more than one
	if len(policyDefinitions) > 1 {
		return policydefinitions.PolicyDefinition{}, fmt.Errorf("loading Policy Definition List: found more than one (%d) policy '%s'", len(policyDefinitions), displayName)
	}

	return policyDefinitions[0], nil
}
func getBuiltInPolicyDefinitionByName(ctx context.Context, client *policydefinitions.PolicyDefinitionsClient, name string) (res policydefinitions.PolicyDefinition, err error) {
	builtIn, err := client.GetBuiltIn(ctx, policydefinitions.NewPolicyDefinitionID(name))
	if err == nil && builtIn.Model != nil {
		return *builtIn.Model, nil
	}
	return policydefinitions.PolicyDefinition{}, err
}

func getPolicyDefinitionByName(ctx context.Context, client *policydefinitions.PolicyDefinitionsClient,
	name, managementGroupName, subscriptionID string) (res policydefinitions.GetOperationResponse, err error) {

	if managementGroupName == "" {
		var builtIn policydefinitions.GetBuiltInOperationResponse
		builtIn, err = client.GetBuiltIn(ctx, policydefinitions.NewPolicyDefinitionID(name))
		if err == nil && builtIn.Model != nil {
			res.HttpResponse = builtIn.HttpResponse
			res.OData = builtIn.OData
			res.Model = builtIn.Model
			return res, nil
		}

		if response.WasNotFound(builtIn.HttpResponse) {
			var getResult policydefinitions.GetOperationResponse
			getResult, err = client.Get(ctx, policydefinitions.NewProviderPolicyDefinitionID(subscriptionID, name))
			if err == nil && getResult.Model != nil {
				return getResult, nil
			}
		}
	} else {
		var result policydefinitions.GetAtManagementGroupOperationResponse
		result, err = client.GetAtManagementGroup(ctx, policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, name))
		if err == nil && result.Model != nil {
			res.HttpResponse = result.HttpResponse
			res.OData = result.OData
			res.Model = result.Model
			return res, nil
		}
	}

	return res, err
}

func getPolicySetDefinitionByName(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, name, managementGroupID, subscriptionID string) (res policysetdefinitions.PolicySetDefinition, err error) {
	if managementGroupID == "" {
		var builtIn policysetdefinitions.GetBuiltInOperationResponse
		builtIn, err = client.GetBuiltIn(ctx, policysetdefinitions.NewPolicySetDefinitionID(name))
		if err == nil && builtIn.Model != nil {
			return *builtIn.Model, nil
		}
		if response.WasNotFound(builtIn.HttpResponse) {
			var result policysetdefinitions.GetOperationResponse
			result, err = client.Get(ctx, policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionID, name))
			if err == nil && result.Model != nil {
				return *result.Model, nil
			}
		}
	} else {
		var result policysetdefinitions.GetAtManagementGroupOperationResponse
		result, err = client.GetAtManagementGroup(ctx, policysetdefinitions.NewProviders2PolicySetDefinitionID(managementGroupID, name))
		if err == nil && result.Model != nil {
			return *result.Model, nil
		}
	}

	return policysetdefinitions.PolicySetDefinition{}, err
}

func getPolicySetDefinitionByDisplayName(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, displayName, managementGroupID, subscriptionID string) (policysetdefinitions.PolicySetDefinition, error) {
	var setDefinitions []policysetdefinitions.PolicySetDefinition
	var err error

	if managementGroupID != "" {
		var result policysetdefinitions.ListByManagementGroupCompleteResult
		result, err = client.ListByManagementGroupComplete(ctx, commonids.NewManagementGroupID(managementGroupID), policysetdefinitions.ListByManagementGroupOperationOptions{})
		setDefinitions = result.Items
	} else {
		var result policysetdefinitions.ListCompleteResult
		result, err = client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), policysetdefinitions.ListOperationOptions{})
		setDefinitions = result.Items
	}
	if err != nil {
		return policysetdefinitions.PolicySetDefinition{}, fmt.Errorf("loading Policy Set Definition List: %+v", err)
	}

	// var results []policy.SetDefinition
	// for setDefinitions.NotDone() {
	// 	def := setDefinitions.Value()
	// 	if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
	// 		results = append(results, def)
	// 	}
	//
	// 	if err := setDefinitions.NextWithContext(ctx); err != nil {
	// 		return policy.SetDefinition{}, fmt.Errorf("loading Policy Set Definition List: %s", err)
	// 	}
	// }

	// throw error when we found none
	if len(setDefinitions) == 0 {
		return policysetdefinitions.PolicySetDefinition{}, fmt.Errorf("loading Policy Set Definition List: could not find policy '%s'", displayName)
	}

	// throw error when we found more than one
	if len(setDefinitions) > 1 {
		return policysetdefinitions.PolicySetDefinition{}, fmt.Errorf("loading Policy Set Definition List: found more than one policy set definition '%s'", displayName)
	}

	return setDefinitions[0], nil
}

func expandParameterDefinitionsValueFromString(jsonString string) (map[string]policydefinitions.ParameterDefinitionsValue, error) {
	var result map[string]policydefinitions.ParameterDefinitionsValue

	err := json.Unmarshal([]byte(jsonString), &result)

	return result, err
}

func flattenParameterDefinitionsValueToString(input *map[string]policydefinitions.ParameterDefinitionsValue) (string, error) {
	if input == nil || len(*input) == 0 {
		return "", nil
	}

	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	compactJson := bytes.Buffer{}
	if err := json.Compact(&compactJson, result); err != nil {
		return "", err
	}

	return compactJson.String(), nil
}

func expandParameterValuesValueFromString(jsonString string) (map[string]assignments.ParameterValuesValue, error) {
	var result map[string]assignments.ParameterValuesValue

	err := json.Unmarshal([]byte(jsonString), &result)

	return result, err
}

func flattenParameterValuesValueToString(input map[string]*policysetdefinitions.ParameterValuesValue) (string, error) {
	if input == nil {
		return "", nil
	}

	// no need to call `json.Compact` for the result of `json.Marshal`, it's compacted already
	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(result), err
}

func flattenParameterValuesValueToStringV2(input *map[string]assignments.ParameterValuesValue) (string, error) {
	if input == nil || *input == nil {
		return "", nil
	}
	bs, err := json.Marshal(input)
	return string(bs), err
}

func getPolicyRoleDefinitionIDs(ruleStr string) (res []string, err error) {
	type policyRule struct {
		Then struct {
			Details struct {
				RoleDefinitionIds []string `json:"roleDefinitionIds"`
			} `json:"details"`
		} `json:"then"`
	}
	var ins policyRule
	err = json.Unmarshal([]byte(ruleStr), &ins)
	res = ins.Then.Details.RoleDefinitionIds
	return
}
