// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	dataplane "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	keyvaultCache   *vaultCache
	managedHSMCache *vaultCache

	ManagedHsmClient *managedhsms.ManagedHsmsClient
	ManagementClient *dataplane.BaseClient
	VaultsClient     *vaults.VaultsClient

	ManagedHSMSDClient              *dataplane.HSMSecurityDomainClient
	ManagedHSMRoleClient            *dataplane.RoleDefinitionsClient
	ManagedHSMRoleAssignmentsClient *dataplane.RoleAssignmentsClient
	ManagedHSMManagementClient      *dataplane.BaseClient
}

func NewClient(o *common.ClientOptions) *Client {
	managedHsmClient := managedhsms.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&managedHsmClient.Client, o.ResourceManagerAuthorizer)

	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	vaultsClient := vaults.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint)

	sdClient := dataplane.NewHSMSecurityDomainClient()
	o.ConfigureClient(&sdClient.Client, o.ManagedHSMAuthorizer)

	mhsmRoleDefineClient := dataplane.NewRoleDefinitionsClient()
	o.ConfigureClient(&mhsmRoleDefineClient.Client, o.ManagedHSMAuthorizer)

	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	mhsmRoleAssignClient := dataplane.NewRoleAssignmentsClient()
	o.ConfigureClient(&mhsmRoleAssignClient.Client, o.ManagedHSMAuthorizer)

	mhsmManagement := dataplane.New()
	o.ConfigureClient(&mhsmManagement.Client, o.ManagedHSMAuthorizer)

	return &Client{
		keyvaultCache:   newVaultCache(),
		managedHSMCache: newVaultCache(),

		ManagedHsmClient:                &managedHsmClient,
		ManagementClient:                &managementClient,
		VaultsClient:                    &vaultsClient,
		ManagedHSMSDClient:              &sdClient,
		ManagedHSMRoleClient:            &mhsmRoleDefineClient,
		ManagedHSMRoleAssignmentsClient: &mhsmRoleAssignClient,
		ManagedHSMManagementClient:      &mhsmManagement,
	}
}
