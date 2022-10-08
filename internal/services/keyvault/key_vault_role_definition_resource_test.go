package keyvault_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultRoleDefinitionResource struct{}

func (a KeyVaultRoleDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	baseURL := state.Attributes["vault_base_url"]
	id := strings.Split(state.ID, "/")
	if len(id) == 0 {
		return utils.Bool(false), nil
	}
	definitionID := id[len(id)-1]
	if definitionID == "" {
		return utils.Bool(false), fmt.Errorf("no role definition id")
	}
	resp, err := client.KeyVault.MHSMRoleClient.Get(ctx, baseURL, "/", id[len(id)-1])
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.RoleDefinitionProperties != nil), nil
}

func (a KeyVaultRoleDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a KeyVaultRoleDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_role_assignment" "test" {
  name  =  "acctest-%[2]d"
role_definition_id  =  "example"
scope  =  "example"
vault_base_url  =  "example"
description  =  "example"
role_type  =  "example"
assignable_scopes  =  ["example"]
permission  {
actions = ["example"]
not_actions = ["example"]
data_actions = ["example"]
not_data_actions = ["example"]
}
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a KeyVaultRoleDefinitionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_role_assignment" "test" {
  
  name  =  "acctest-%[2]d"
  role_definition_id  =  "example"
  scope  =  "example"
  vault_base_url  =  "example"
  description  =  "example"
  role_type  =  "example"
  assignable_scopes  =  ["example"]
  permission  {
actions = ["example"]
not_actions = ["example"]
data_actions = ["example"]
not_data_actions = ["example"]
}
}
`, a.template(data), data.RandomInteger)
}

// We cannot run this Test for now, because we cannot activate mhsm keyvault by terraform now.
func testAccKeyVaultRoleAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, keyvault.KeyVaultRoleDefinitionResource{}.ResourceType(), "test")
	r := KeyVaultRoleDefinitionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_global").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_global").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}
