// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func testAccKeyVaultManagedHardwareSecurityModule_basicKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_key", "test")
	r := KeyVaultManagedHSMKeyResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

type KeyVaultManagedHSMKeyResource struct{}

// Exists implements types.TestResource.
func (KeyVaultManagedHSMKeyResource) Exists(ctx context.Context, cli *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := cli.KeyVault
	subscriptionId := cli.Account.SubscriptionId

	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	managedHSMIdRaw, err := client.ManagedHSMIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil || managedHSMIdRaw == nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	managedHSMId, err := managedhsms.ParseManagedHSMID(*managedHSMIdRaw)
	if err != nil {
		return nil, err
	}

	ok, err := client.ManagedHSMExists(ctx, *managedHSMId)
	if err != nil || !ok {
		return nil, fmt.Errorf("checking if managed HSM %q for name %q in Vault at url %q exists: %v", *managedHSMId, id.Name, id.KeyVaultBaseUrl, err)
	}

	resp, err := client.ManagedHSMManagementClient.GetKey(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Key Vault Key %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Key != nil), nil
}

func (k *KeyVaultManagedHSMKeyResource) basic(data acceptance.TestData) string {
	hsm := KeyVaultManagedHardwareSecurityModuleResource{}.download(data, 3)

	return fmt.Sprintf(`
	



%s

locals {
  assignmentUserName = "706c03c7-69ad-33e5-2796-b3380d3a6e1a"
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "user" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "%s"
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "officer" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "%s"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "user" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name               = local.assignmentUserName
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "officer" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name               = "d9a89332-9ec9-11ee-a8a5-00155dbdfff5"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}


resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "hsmkey%s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "EC-HSM"

  key_opts = [
    "sign",
    "verify",
  ]

  depends_on = [azurerm_key_vault_managed_hardware_security_module_role_assignment.user, azurerm_key_vault_managed_hardware_security_module_role_assignment.officer]
}
`, hsm, managedHSMCryptoUserRoleID, managedHSMCryptoOfficerRoleID, data.RandomString)
}
