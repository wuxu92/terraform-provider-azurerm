package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

type FirewallPublicIPAddressAssociationId struct {
	FirewallID        FirewallId
	PublicIPAddressID parse.PublicIpAddressId
}

func (f FirewallPublicIPAddressAssociationId) ID() string {
	return fmt.Sprintf("%s|%s", f.FirewallID, f.PublicIPAddressID)
}

func (f FirewallPublicIPAddressAssociationId) String() string {
	return f.ID()
}

func NewFirewallPublicIPAddressAssociationID(firewallID, publicIP string) (*FirewallPublicIPAddressAssociationId, error) {
	return FirewallPublicIPAddressAssociationID(strings.Join([]string{firewallID, publicIP}, "|"))
}

func FirewallPublicIPAddressAssociationID(input string) (*FirewallPublicIPAddressAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format `{firewallID}|{publicIPAddressID} but got %q", input)
	}

	firewallID, err := FirewallID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Firewall ID %q: %+v", segments[0], err)
	}

	// whilst we need the Resource ID, we may as well validate it
	publicIPAddress := segments[1]
	ipAddressID, err := parse.PublicIpAddressID(publicIPAddress)
	if err != nil {
		return nil, fmt.Errorf("parsing Public IP Address ID %q: %+v", publicIPAddress, err)
	}

	return &FirewallPublicIPAddressAssociationId{
		FirewallID:        *firewallID,
		PublicIPAddressID: *ipAddressID,
	}, nil
}
