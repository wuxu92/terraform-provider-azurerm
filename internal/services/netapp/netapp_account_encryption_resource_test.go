// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppAccountEncryptionResource struct{}

func TestAccNetAppAccountEncryption_cmkSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkSystemAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_cmkUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkUserAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_cmkManagedHSMKey(t *testing.T) {
	if os.Getenv("ARM_TEST_HSM_KEY") == "" {
		t.Skip("Skipping as ARM_TEST_HSM_KEY is not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	data.Locations.Primary = "centralus"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkManagedHSM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_managed_hsm_key").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	regexInitialKey := regexp.MustCompile(`anfenckey\d+$`)
	regexNewKey := regexp.MustCompile(`.*anfenckey-new.*`)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyUpdate1(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").MatchesRegex(regexInitialKey),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyUpdate2(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_key").MatchesRegex(regexNewKey),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppAccountEncryptionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := netappaccounts.ParseNetAppAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.AccountClient.AccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp Account (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r NetAppAccountEncryptionResource) cmkSystemAssigned(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = "%[3]s"

  sku_name = "standard"
}

resource "azurerm_key_vault_access_policy" "test-currentuser" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "Create",
    "Delete",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.test-currentuser
  ]
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "test-systemassigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
  object_id    = azurerm_netapp_account.test.identity.0.principal_id

  key_permissions = [
    "Get",
    "Encrypt",
    "Decrypt"
  ]
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id = azurerm_netapp_account.test.id

  system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id

  encryption_key = azurerm_key_vault_key.test.versionless_id

  depends_on = [
    azurerm_key_vault_access_policy.test-systemassigned
  ]
}
`, r.template(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) cmkManagedHSM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "user-assigned-identity-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                       = "acckv%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
    ]
    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
  tags = {
    environment = "Production"
  }
}
resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
  name         = "acchsmcert${count.index}"
  key_vault_id = azurerm_key_vault.test.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    lifetime_action {
      action {
        action_type = "AutoRenew"
      }
      trigger {
        days_before_expiry = 30
      }
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      extended_key_usage = []
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled   = true
  soft_delete_retention_days = 7

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 3
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "officer" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "1e243909-064c-6ac3-84e9-1c8bf8d6ad23"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "cryptor" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "1e243909-064c-6ac3-84e9-1c8bf8d6ad22"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "identity" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "ff248f45-efda-4122-bec0-71460ecc8add"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[3]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 4096
  key_opts       = ["encrypt", "decrypt"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.officer,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.cryptor,
  ]
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id          = azurerm_netapp_account.test.id
  user_assigned_identity_id  = azurerm_user_assigned_identity.test.id
  encryption_managed_hsm_key = azurerm_key_vault_managed_hardware_security_module_key.test.id
  depends_on                 = [azurerm_key_vault_managed_hardware_security_module_role_assignment.identity]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r NetAppAccountEncryptionResource) cmkUserAssigned(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "user-assigned-identity-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = "%[3]s"

  sku_name = "standard"

  access_policy {
    tenant_id = "%[3]s"
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
      "Create",
      "Delete",
      "WrapKey",
      "UnwrapKey",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }

  access_policy {
    tenant_id = "%[3]s"
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "Get",
      "Encrypt",
      "Decrypt"
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id = azurerm_netapp_account.test.id

  user_assigned_identity_id = azurerm_user_assigned_identity.test.id

  encryption_key = azurerm_key_vault_key.test.versionless_id
}
`, r.template(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) keyUpdate1(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = "%[3]s"

  sku_name = "standard"
}

resource "azurerm_key_vault_access_policy" "test-currentuser" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "Create",
    "Delete",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.test-currentuser
  ]
}

resource "azurerm_key_vault_key" "test-new-key" {
  name         = "anfenckey-new%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_key.test,
    azurerm_key_vault_access_policy.test-currentuser
  ]
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "test-systemassigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
  object_id    = azurerm_netapp_account.test.identity.0.principal_id

  key_permissions = [
    "Get",
    "Encrypt",
    "Decrypt"
  ]
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id = azurerm_netapp_account.test.id

  system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id

  encryption_key = azurerm_key_vault_key.test.versionless_id

  depends_on = [
    azurerm_key_vault_access_policy.test-systemassigned
  ]
}
`, r.template(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) keyUpdate2(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[2]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  tenant_id                       = "%[3]s"

  sku_name = "standard"
}

resource "azurerm_key_vault_access_policy" "test-currentuser" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "Create",
    "Delete",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_access_policy.test-currentuser
  ]
}

resource "azurerm_key_vault_key" "test-new-key" {
  name         = "anfenckey-new%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault_key.test,
    azurerm_key_vault_access_policy.test-currentuser
  ]
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "test-systemassigned" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
  object_id    = azurerm_netapp_account.test.identity.0.principal_id

  key_permissions = [
    "Get",
    "Encrypt",
    "Decrypt"
  ]
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id = azurerm_netapp_account.test.id

  system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id

  encryption_key = azurerm_key_vault_key.test-new-key.versionless_id

  depends_on = [
    azurerm_key_vault_access_policy.test-systemassigned
  ]
}
`, r.template(data), data.RandomInteger, tenantID)
}

func (NetAppAccountEncryptionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }

    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
