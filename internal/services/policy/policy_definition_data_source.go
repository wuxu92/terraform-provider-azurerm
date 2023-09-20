// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmPolicyDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: policyDefinitionReadFunc(false),

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: policyDefinitionDataSourceSchema(),
	}
}

func policyDefinitionDataSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"management_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"policy_rule": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parameters": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func policyDefinitionReadFunc(builtInOnly bool) func(d *pluginsdk.ResourceData, meta interface{}) error {
	return func(d *pluginsdk.ResourceData, meta interface{}) error {
		client := meta.(*clients.Client).Policy.DefinitionsClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		displayName := d.Get("display_name").(string)
		name := d.Get("name").(string)
		managementGroupName := ""
		if v, ok := d.GetOk("management_group_name"); ok {
			managementGroupName = v.(string)
		}

		subscriptionID := meta.(*clients.Client).Account.SubscriptionId
		var policyDefinition policydefinitions.PolicyDefinition
		var err error
		// one of display_name and name must be non-empty, this is guaranteed by schema
		if displayName != "" {
			policyDefinition, err = getPolicyDefinitionByDisplayName(ctx, client, subscriptionID, displayName, managementGroupName, builtInOnly)
			if err != nil {
				return fmt.Errorf("reading Policy Definition (Display Name %q): %+v", displayName, err)
			}
		}
		if name != "" {
			if builtInOnly && managementGroupName == "" {
				policyDefinition, err = getBuiltInPolicyDefinitionByName(ctx, client, name)
			} else {
				var getResponse policydefinitions.GetOperationResponse
				getResponse, err = getPolicyDefinitionByName(ctx, client, name, managementGroupName, subscriptionID)
				policyDefinition = pointer.From(getResponse.Model)
			}
			if err != nil {
				return fmt.Errorf("reading Policy Definition %q: %+v", name, err)
			}
		}

		if policyDefinition.Id == nil {
			return fmt.Errorf("reading policy definition id as nil")
		}

		id, err := policydefinitions.ParseProviderPolicyDefinitionID(*policyDefinition.Id)
		if err != nil {
			return fmt.Errorf("parsing Policy Definition %q: %+v", *policyDefinition.Id, err)
		}

		d.SetId(id.ID())
		d.Set("name", policyDefinition.Name)
		d.Set("type", policyDefinition.Type)

		if prop := policyDefinition.Properties; prop != nil {
			d.Set("display_name", prop.DisplayName)
			d.Set("description", prop.Description)
			d.Set("policy_type", prop.PolicyType)
			d.Set("mode", prop.Mode)
			if rule := prop.PolicyRule; rule != nil {
				policyRule := (*rule).(map[string]interface{})
				if policyRuleStr := flattenJSON(policyRule); policyRuleStr != "" {
					d.Set("policy_rule", policyRuleStr)
					roleIDs, _ := getPolicyRoleDefinitionIDs(policyRuleStr)
					d.Set("role_definition_ids", roleIDs)
				} else {
					return fmt.Errorf("flattening Policy Definition Rule %q: %+v", name, err)
				}
			}

			if metadataStr := flattenJSON(prop.Metadata); metadataStr != "" {
				d.Set("metadata", metadataStr)
			}

			if parametersStr, err := flattenParameterDefinitionsValueToString(pointer.From(prop.Parameters)); err == nil {
				d.Set("parameters", parametersStr)
			} else {
				return fmt.Errorf("failed to flatten Policy Parameters %q: %+v", name, err)
			}
		}

		return nil
	}
}
