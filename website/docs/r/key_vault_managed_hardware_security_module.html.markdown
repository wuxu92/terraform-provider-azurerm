---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module"
description: |-
  Manages a Key Vault Managed Hardware Security Module.
---

# azurerm_key_vault_managed_hardware_security_module

Manages a Key Vault Managed Hardware Security Module.

~> **Note:** the Azure Provider includes a Feature Toggle which will purge a Key Vault Managed Hardware Security Module resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_hardware_security_modules_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block#purge_soft_deleted_hardware_security_modules_on_destroy) for more information.

~> **Note:** To create a key for the Manage HSM, you can use the `[azurerm_key_vault_key](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key)`.

## Example Usage

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_hardware_security_modules_on_destroy = true
    }
  }
}
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault_managed_hardware_security_module" "example" {
  name                       = "exampleKVHsm"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  sku_name                   = "Standard_B1"
  purge_protection_enabled   = false
  soft_delete_retention_days = 90
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]

  tags = {
    Env = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Key Vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `admin_object_ids` - (Required) Specifies a list of administrators object IDs for the key vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `sku_name` - (Required) The Name of the SKU used for this Key Vault Managed Hardware Security Module. Possible value is `Standard_B1`. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The Azure Active Directory Tenant ID that should be used for authenticating requests to the key vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `purge_protection_enabled` - (Optional) Is Purge Protection enabled for this Key Vault Managed Hardware Security Module? Changing this forces a new resource to be created.

* `soft_delete_retention_days` - (Optional) The number of days that items should be retained for once soft-deleted. This value can be between `7` and `90` days. Defaults to `90`. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether traffic from public networks is permitted. Defaults to `true`. Changing this forces a new resource to be created.

* `network_acls` - (Optional) A `network_acls` block as defined below.

* `activate_config` - (Optional) A `activate_config` block used to activate this Managed HSM as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

A `network_acls` block supports the following:

* `bypass` - (Required) Specifies which traffic can bypass the network rules. Possible values are `AzureServices` and `None`.

* `default_action` - (Required) The Default Action to use. Possible values are `Allow` and `Deny`.

---

A `activate_config` block supports the following:

* `security_domain_certificate` - (Required) A list of KeyVault certificates resource ID(minimum of three and up to a maximum of 10) to activate this Managed HSM. More information see [activate-your-managed-hsm](https://learn.microsoft.com/en-us/azure/key-vault/managed-hsm/quick-create-cli#activate-your-managed-hsm)

* `security_domain_quorum` - (Required) Specifies the minimum number of shares required to decrypt the security domain for recovery. This is required the `security_domain_certificate` is provided. The value must between 2 and 10 (inclusive).

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Secret Managed Hardware Security Module ID.

* `hsm_uri` - The URI of the Key Vault Managed Hardware Security Module, used for performing operations on keys.

* `security_domain_enc_data` - The sensitive data will be used for disaster recovery or for creating another Managed HSM that shares same security domain so the two can share keys.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Key Vault Managed Hardware Security Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module.
* `delete` - (Defaults to 60 minutes) Used when deleting the Key Vault Managed Hardware Security Module.

## Import

Key Vault Managed Hardware Security Module can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/managedHSMs/hsm1
```
