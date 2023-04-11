package firewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FirewallPublicIpAssociationResource struct{}

func (a FirewallPublicIpAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallPublicIPAddressAssociationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Firewall.AzureFirewallsClient.Get(ctx, id.FirewallID.ResourceGroup, id.FirewallID.AzureFirewallName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	if resp.AzureFirewallPropertiesFormat == nil || resp.AzureFirewallPropertiesFormat.IPConfigurations == nil {
		for _, ip := range *resp.AzureFirewallPropertiesFormat.IPConfigurations {
			if ip.PublicIPAddress != nil && *ip.PublicIPAddress.ID == id.PublicIPAddressID.ID() {
				return pointer.To(true), nil
			}
		}
	}
	return pointer.To(false), nil
}

func TestAccFirewallPublicIpAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, firewall.PublicIPAssociationResource{}.ResourceType(), "test")
	r := FirewallPublicIpAssociationResource{}
	fmt.Printf("%s", r.basic(data))
	t.Fatal()
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

func (a FirewallPublicIpAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_firewall_public_ip_association" "test" {
  firewall_id          = azurerm_firewall.test.id
  name                 = "configuration"
  public_ip_address_id = azurerm_public_ip.test.id
}
`, a.template(data))
}

func (a FirewallPublicIpAssociationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s


resource "azurerm_firewall_public_ip_association" "test" {
  firewall_id          = azurerm_firewall.test.id
  name                 = "configurationFoo"
  public_ip_address_id = azurerm_public_ip.test.id
}
`, a.template(data))
}

func (a FirewallPublicIpAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  threat_intel_mode = "Deny"
}
`, data.RandomInteger, data.Locations.Primary)
}
