package firewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type FirewallPublicIpAssociationResource struct{}

func (a FirewallPublicIpAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallPublicIpAssociationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AzureFirewallsClient.Get(ctx, id.ResourceGroup, id.XXName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.FirewallPublicIpAssociationProperties != nil), nil
}

func (a FirewallPublicIpAssociationResource) template(data acceptance.TestData) string {
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

func (a FirewallPublicIpAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`




%s

resource "azurerm_firewall_public_ip_association" "test" {
  firewall_id  = "example"
  public_ip_id = "example"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a FirewallPublicIpAssociationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`




%s

resource "azurerm_firewall_public_ip_association" "test" {

  firewall_id  = "example"
  public_ip_id = "example"
}
`, a.template(data), data.RandomInteger)
}

func TestAccFirewallPublicIpAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, network.FirewallPublicIpAssociationResource{}.ResourceType(), "test")
	r := FirewallPublicIpAssociationResource{}
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
