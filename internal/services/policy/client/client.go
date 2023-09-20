// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	// "github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy"

	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2020-06-25/guestconfigurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-07-01-preview/policyexemptions" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policyassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policydefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-04-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssignmentsClient                   *policyassignments.PolicyAssignmentsClient
	DefinitionsClient                   *policydefinitions.PolicyDefinitionsClient
	ExemptionsClient                    *policyexemptions.PolicyExemptionsClient
	GuestConfigurationAssignmentsClient *guestconfigurationassignments.GuestConfigurationAssignmentsClient
	RemediationsClient                  *remediations.RemediationsClient
	SetDefinitionsClient                *policysetdefinitions.PolicySetDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	assignmentsClient, err := policyassignments.NewPolicyAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PolicyAssignments client: %+v", err)
	}
	o.Configure(assignmentsClient.Client, o.Authorizers.ResourceManager)

	definitionsClient, err := policydefinitions.NewPolicyDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PolicyDefinitions client: %+v", err)
	}
	o.Configure(definitionsClient.Client, o.Authorizers.ResourceManager)

	exemptionsClient, err := policyexemptions.NewPolicyExemptionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PolicyExemptions client: %+v", err)
	}
	o.Configure(exemptionsClient.Client, o.Authorizers.ResourceManager)

	guestConfigurationAssignmentsClient, err := guestconfigurationassignments.NewGuestConfigurationAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Guest Configuration Assignments Client:  %+v", err)
	}
	o.Configure(guestConfigurationAssignmentsClient.Client, o.Authorizers.ResourceManager)

	remediationsClient, err := remediations.NewRemediationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Remediations client: %+v", err)
	}
	o.Configure(remediationsClient.Client, o.Authorizers.ResourceManager)

	setDefinitionsClient, err := policysetdefinitions.NewPolicySetDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PolicySetDefinitions client: %+v", err)
	}
	o.Configure(setDefinitionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AssignmentsClient:                   assignmentsClient,
		DefinitionsClient:                   definitionsClient,
		ExemptionsClient:                    exemptionsClient,
		GuestConfigurationAssignmentsClient: guestConfigurationAssignmentsClient,
		RemediationsClient:                  remediationsClient,
		SetDefinitionsClient:                setDefinitionsClient,
	}, nil
}
