// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	// nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmPolicySetDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmPolicySetDefinitionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"parameters": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"policy_definitions": { // TODO -- remove in the next major version
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"policy_definition_reference": { // TODO -- rename this back to `policy_definition` after the deprecation
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"policy_definition_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"parameters": { // TODO -- remove this attribute in the next major version
							Type:     pluginsdk.TypeMap,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"parameter_values": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"reference_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"policy_group_names": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"policy_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"policy_definition_group": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"display_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"category": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"additional_metadata_resource_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmPolicySetDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	managementGroupID := d.Get("management_group_name").(string)

	var setDefinition policysetdefinitions.PolicySetDefinition
	var err error

	// we marked `display_name` and `name` as `ExactlyOneOf`, therefore there will only be one of display_name and name that have non-empty value here
	if displayName != "" {
		query := newQueryPolicy(subscriptionID, displayName, managementGroupID)
		setDefinition, err = getPolicySetDefinitionByDisplayName(ctx, client, query)
		if err != nil {
			return fmt.Errorf("reading Policy Set Definition (Display Name %q): %+v", displayName, err)
		}
	}
	if name != "" {
		query := newQueryPolicy(subscriptionID, name, managementGroupID)
		getRes, err := getPolicySetDefinitionByName(ctx, client, query)
		if err != nil {
			return fmt.Errorf("reading Policy Set Definition %q: %+v", name, err)
		}
		if getRes.Model != nil {
			setDefinition = *getRes.Model
		}
	}

	if pointer.From(setDefinition.Id) == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Set Definition %q", name)
	}

	id, err := policysetdefinitions.ParsePolicySetDefinitionID(*setDefinition.Id)
	if err != nil {
		return fmt.Errorf("parsing Policy Set Definition %q: %+v", *setDefinition.Id, err)
	}

	d.SetId(id.ID())
	d.Set("name", setDefinition.Name)
	if prop := setDefinition.Properties; prop != nil {
		d.Set("display_name", prop.DisplayName)
		d.Set("description", prop.Description)
		d.Set("policy_type", prop.PolicyType)
		d.Set("metadata", flattenJSON(prop.Metadata))

		if paramsStr, err := flattenPolicySetParameterDefinitionsValueToString(prop.Parameters); err != nil {
			return fmt.Errorf("flattening JSON for `parameters`: %+v", err)
		} else {
			d.Set("parameters", paramsStr)
		}

		definitionBytes, err := json.Marshal(prop.PolicyDefinitions)
		if err != nil {
			return fmt.Errorf("flattening JSON for `policy_defintions`: %+v", err)
		}
		d.Set("policy_definitions", string(definitionBytes))

		references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitions(prop.PolicyDefinitions)
		if err != nil {
			return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
		}
		if err := d.Set("policy_definition_reference", references); err != nil {
			return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
		}

		if err := d.Set("policy_definition_group", flattenAzureRMPolicySetDefinitionPolicyGroups(prop.PolicyDefinitionGroups)); err != nil {
			return fmt.Errorf("setting `policy_definition_group`: %+v", err)
		}
	}

	return nil
}
