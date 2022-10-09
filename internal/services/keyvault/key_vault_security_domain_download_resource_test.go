package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultSecurityDomainDownloadResource struct{}

func (a KeyVaultSecurityDomainDownloadResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.KeyVaultSecurityDomainDownloadID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.KeyVault.MHSMSDClient.Get(ctx, id.ResourceGroup, id.XXName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.KeyVaultSecurityDomainDownloadProperties != nil), nil
}

func (a KeyVaultSecurityDomainDownloadResource) template(data acceptance.TestData) string {
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

func (a KeyVaultSecurityDomainDownloadResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_security_domain_download" "test" {
  certificate  {
key_vault_key_id = "example"
}
required  =  %[2]d
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a KeyVaultSecurityDomainDownloadResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_security_domain_download" "test" {
  
  certificate  {
key_vault_key_id = "example"
}
  required  =  %[2]d
}
`, a.template(data), data.RandomInteger)
}

func TestAccKeyVaultSecurityDomainDownload_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, keyvault.KeyVaultSecurityDomainDownloadResource{}.ResourceType(), "test")
	r := KeyVaultSecurityDomainDownloadResource{}
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
