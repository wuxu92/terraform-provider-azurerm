// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmPolicyDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmPolicyDefinitionCreateUpdate,
		Update: resourceArmPolicyDefinitionCreateUpdate,
		Read:   resourceArmPolicyDefinitionRead,
		Delete: resourceArmPolicyDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PolicyDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceArmPolicyDefinitionSchema(),
	}
}

func resourceArmPolicyDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyType := d.Get("policy_type").(string)
	mode := d.Get("mode").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)

	managementGroupName := ""
	if v, ok := d.GetOk("management_group_id"); ok {
		id, err := mgmtGrpParse.ManagementGroupID(v.(string))
		if err != nil {
			return err
		}
		managementGroupName = id.Name
	}

	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	if d.IsNewResource() {
		existing, err := getPolicyDefinitionByName(ctx, client, name, managementGroupName, subscriptionID)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Policy Definition %q: %+v", name, err)
			}
		}

		if model := existing.Model; model != nil {
			if model.Id != nil && *model.Id != "" {
				return tf.ImportAsExistsError("azurerm_policy_definition", *existing.Model.Id)
			}
		}
	}

	properties := policydefinitions.PolicyDefinitionProperties{
		PolicyType:  pointer.To(policydefinitions.PolicyType(policyType)),
		Mode:        utils.String(mode),
		DisplayName: utils.String(displayName),
		Description: utils.String(description),
	}

	if policyRuleString := d.Get("policy_rule").(string); policyRuleString != "" {
		policyRule, err := pluginsdk.ExpandJsonFromString(policyRuleString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_rule`: %+v", err)
		}
		properties.PolicyRule = pointer.To(interface{}(policyRule))
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
		}
		properties.Metadata = pointer.To(interface{}(metaData))
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		parameters, err := expandParameterDefinitionsValueFromString(parametersString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
		}
		properties.Parameters = pointer.To(parameters)
	}

	definition := policydefinitions.PolicyDefinition{
		Name:       utils.String(name),
		Properties: &properties,
	}

	var err error

	if managementGroupName == "" {
		_, err = client.CreateOrUpdate(ctx, policydefinitions.NewProviderPolicyDefinitionID(subscriptionID, name), definition)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, name), definition)
	}

	if err != nil {
		return fmt.Errorf("creating/updating Policy Definition %q: %+v", name, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Definition %q to become available", name)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"404"},
		Target:  []string{"200"},
		Refresh: func() (result interface{}, state string, err error) {
			res, err := getPolicyDefinitionByName(ctx, client, name, managementGroupName, subscriptionID)
			if err != nil {
				return nil, strconv.Itoa(res.HttpResponse.StatusCode), fmt.Errorf("issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q: %+v", name, err)
			}

			return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
		},
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Policy Definition %q to become available: %+v", name, err)
	}

	resp, err := getPolicyDefinitionByName(ctx, client, name, managementGroupName, subscriptionID)
	if err != nil {
		return err
	}

	if model := resp.Model; model != nil {

	}
	if resp.Model == nil || resp.Model.Id == nil || *resp.Model.Id == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Assignment %q", name)
	}

	id, err := policydefinitions.ParsePolicyDefinitionID(*resp.Model.Id)
	if err != nil {
		return fmt.Errorf("failed to flatten Policy Parameters %q: %+v", *resp.Model.Id, err)
	}
	d.SetId(id.ID())

	return resourceArmPolicyDefinitionRead(d, meta)
}

func resourceArmPolicyDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	switch scopeId := id.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
	}

	resp, err := getPolicyDefinitionByName(ctx, client, id.Name, managementGroupName, subscriptionID)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading Policy Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Policy Definition %+v", err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)

		d.Set("management_group_id", managementGroupName)
		if managementGroupName != "" {
			d.Set("management_group_id", managementGroupId.ID())
		}

		if props := model.Properties; props != nil {
			d.Set("policy_type", props.PolicyType)
			d.Set("mode", props.Mode)
			d.Set("display_name", props.DisplayName)
			d.Set("description", props.Description)

			if policyRuleStr := flattenJSON(props.PolicyRule); policyRuleStr != "" {
				d.Set("policy_rule", policyRuleStr)
				roleIDs, _ := getPolicyRoleDefinitionIDs(policyRuleStr)
				d.Set("role_definition_ids", roleIDs)
			}

			if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
				d.Set("metadata", metadataStr)
			}

			if parametersStr, err := flattenParameterDefinitionsValueToString(props.Parameters); err == nil {
				d.Set("parameters", parametersStr)
			} else {
				return fmt.Errorf("flattening policy definition parameters %+v", err)
			}
		}
	}

	return nil
}

func resourceArmPolicyDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	switch scopeId := id.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupName = scopeId.ManagementGroupName
	}

	var resp *http.Response
	if managementGroupName == "" {
		var deleteRes policydefinitions.DeleteOperationResponse
		deleteRes, err = client.Delete(ctx, policydefinitions.NewProviderPolicyDefinitionID(subscriptionID, id.Name))
		resp = deleteRes.HttpResponse
	} else {
		var deleteRes policydefinitions.DeleteAtManagementGroupOperationResponse
		deleteRes, err = client.DeleteAtManagementGroup(ctx, policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, id.Name))
		resp = deleteRes.HttpResponse
	}

	if err != nil {
		if response.WasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting Policy Definition %q: %+v", id.Name, err)
	}

	return nil
}

func flattenJSON(stringMap interface{}) string {
	if stringMap != nil {
		if v, ok := stringMap.(*interface{}); ok {
			stringMap = *v
		}
		value := stringMap.(map[string]interface{})
		jsonString, err := pluginsdk.FlattenJsonToString(value)
		if err == nil {
			return jsonString
		}
	}

	return ""
}

func resourceArmPolicyDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(policy.TypeBuiltIn),
				string(policy.TypeCustom),
				string(policy.TypeNotSpecified),
				string(policy.TypeStatic),
			}, false),
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"All",
					"Indexed",
					"Microsoft.ContainerService.Data",
					"Microsoft.CustomerLockbox.Data",
					"Microsoft.DataCatalog.Data",
					"Microsoft.KeyVault.Data",
					"Microsoft.Kubernetes.Data",
					"Microsoft.MachineLearningServices.Data",
					"Microsoft.Network.Data",
					"Microsoft.Synapse.Data",
				}, false,
			),
		},

		"management_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"policy_rule": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"metadata": metadataSchema(),
	}
}
