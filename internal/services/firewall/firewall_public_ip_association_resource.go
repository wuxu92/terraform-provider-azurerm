package firewall

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type PublicIpAssociationModel struct {
	FirewallId        string `tfschema:"firewall_id"`
	Name              string `tfschema:"name"`
	PublicIPAddressId string `tfschema:"public_ip_address_id"`
}

type publicIPAssociationOperation int

func (p publicIPAssociationOperation) String() string {
	if p >= endOfPublicIPAssociationOperation {
		return ""
	}
	return []string{
		"creating",
		"reading",
		"updating",
		"deleting",
	}[p]
}

const (
	createPublicIPAssociation publicIPAssociationOperation = iota
	readPublicIPAssociation
	updatePublicIPAssociation
	deletePublicIPAssociation
	endOfPublicIPAssociationOperation
)

type PublicIPAssociationResource struct{}

var _ sdk.Resource = (*PublicIPAssociationResource)(nil)

func (m PublicIPAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"firewall_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_ip_address_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate2.PublicIpAddressID,
		},
	}
}

func (m PublicIPAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m PublicIPAssociationResource) ModelObject() interface{} {
	return &PublicIpAssociationModel{}
}

func (m PublicIPAssociationResource) ResourceType() string {
	return "azurerm_firewall_public_ip_association"
}

func (m PublicIPAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := m.updateFirewallResource(ctx, meta, createPublicIPAssociation)
			if err != nil {
				return err
			}
			meta.SetID(id)
			return nil
		},
	}
}

func (m PublicIPAssociationResource) updateFirewallResource(ctx context.Context, meta sdk.ResourceMetaData, op publicIPAssociationOperation) (req *parse.FirewallPublicIPAddressAssociationId, err error) {
	client := meta.Client.Firewall.AzureFirewallsClient

	var model PublicIpAssociationModel
	if op == updatePublicIPAssociation || op == createPublicIPAssociation {
		if err := meta.Decode(&model); err != nil {
			return nil, err
		}
	}

	var id *parse.FirewallPublicIPAddressAssociationId
	if op == createPublicIPAssociation {
		id, err = parse.NewFirewallPublicIPAddressAssociationID(model.FirewallId, model.PublicIPAddressId)
	} else {
		id, err = parse.FirewallPublicIPAddressAssociationID(meta.ResourceData.Id())
	}
	if err != nil {
		return nil, err
	}

	firewallResource, err := client.Get(ctx, id.FirewallID.ResourceGroup, id.FirewallID.AzureFirewallName)
	if err != nil {
		return nil, fmt.Errorf("retreiving %s: %v", id, err)
	}

	if utils.ResponseWasNotFound(firewallResource.Response) || firewallResource.AzureFirewallPropertiesFormat == nil {
		return nil, fmt.Errorf("no such firewall: %s", id.FirewallID)
	}

	// check if public configured already
	var readIPConfiguration *network.AzureFirewallIPConfiguration
	var publicIPs []network.AzureFirewallIPConfiguration
	if ips := firewallResource.AzureFirewallPropertiesFormat.IPConfigurations; ips != nil {
		for _, ip := range *ips {
			if ip.PublicIPAddress != nil && ip.PublicIPAddress.ID != nil && *ip.PublicIPAddress.ID == model.PublicIPAddressId {
				switch op {
				case createPublicIPAssociation:
					return nil, tf.ImportAsExistsError(m.ResourceType(), id.ID())
				case updatePublicIPAssociation:
					ip.Name = pointer.To(model.Name)
					publicIPs = append(publicIPs, ip)
				case deletePublicIPAssociation:
				// skip
				case readPublicIPAssociation:
					readIPConfiguration = &ip
				}
			} else {
				publicIPs = append(publicIPs, ip)
			}
		}
	}

	if op == readPublicIPAssociation {
		if readIPConfiguration == nil {
			err = meta.MarkAsGone(id)
		} else {
			// encode meta
			model.Name = pointer.From(readIPConfiguration.Name)
			model.FirewallId = id.FirewallID.ID()
			model.PublicIPAddressId = pointer.From(readIPConfiguration.PublicIPAddress.ID)
			err = meta.Encode(model)
		}
		return nil, err
	}

	// create/update/delete
	if op == createPublicIPAssociation {
		publicIPs = append(publicIPs, network.AzureFirewallIPConfiguration{
			AzureFirewallIPConfigurationPropertiesFormat: &network.AzureFirewallIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.SubResource{
					ID: pointer.To(model.PublicIPAddressId),
				},
			},
			Name: pointer.To(model.Name),
		})
	}

	firewallResource.IPConfigurations = &publicIPs
	future, err := client.CreateOrUpdate(ctx, id.FirewallID.ResourceGroup, id.FirewallID.AzureFirewallName, firewallResource)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %+v", op, id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return nil, fmt.Errorf("waiting for %s of %s: %+v", op, id, err)
	}
	return id, nil
}

func (m PublicIPAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			_, err := m.updateFirewallResource(ctx, meta, readPublicIPAssociation)
			return err
		},
	}
}

func (m PublicIPAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			_, err := m.updateFirewallResource(ctx, meta, deletePublicIPAssociation)
			return err
		},
	}
}

func (m PublicIPAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FirewallPublicIpAssociationID
}
