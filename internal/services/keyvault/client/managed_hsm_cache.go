package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
)

func managedHSMCacheKey(id managedhsms.ManagedHSMId) string {
	return strings.ToLower(id.ManagedHSMName)
}

func (c *Client) AddManagedHSMToCache(id managedhsms.ManagedHSMId, dataPlaneUri string) {
	c.managedHSMCache.addCache(managedHSMCacheKey(id), id.ID(), id.ResourceGroupName, dataPlaneUri)
}

func (c *Client) BaseUriForManagedHSM(ctx context.Context, id managedhsms.ManagedHSMId) (*string, error) {
	if cached := c.managedHSMCache.getCachedBaseUri(managedHSMCacheKey(id)); cached != nil {
		return cached, nil
	}

	resp, err := c.ManagedHsmClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", id)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	hsmUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.HsmUri != nil {
			hsmUri = *model.Properties.HsmUri
		}
	}
	if hsmUri == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.HsmUri` was nil", id)
	}

	c.AddManagedHSMToCache(id, hsmUri)
	return &hsmUri, nil
}

func (c *Client) ManagedHSMExists(ctx context.Context, id managedhsms.ManagedHSMId) (bool, error) {
	if c.managedHSMCache.getCachedItem(managedHSMCacheKey(id)) != nil {
		return true, nil
	}

	resp, err := c.ManagedHsmClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	hsmUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.HsmUri != nil {
			hsmUri = *model.Properties.HsmUri
		}
	}
	if hsmUri == "" {
		return false, fmt.Errorf("retrieving %s: `properties.HsmUri` was nil", id)
	}
	c.AddManagedHSMToCache(id, hsmUri)

	return true, nil
}

func (c *Client) ManagedHSMIDFromBaseUrl(ctx context.Context, subscriptionId commonids.SubscriptionId, managedHSMBaseUrl string) (*string, error) {
	managedHSMName, err := c.parseNameFromBaseUrl(managedHSMBaseUrl)
	if err != nil {
		return nil, err
	}

	if cached := c.managedHSMCache.getCachedID(strings.ToLower(*managedHSMName)); cached != nil {
		return cached, nil
	}

	opts := managedhsms.DefaultListBySubscriptionOperationOptions()
	results, err := c.ManagedHsmClient.ListBySubscriptionComplete(ctx, subscriptionId, opts)
	if err != nil {
		return nil, fmt.Errorf("listing the Managed HSM within %s: %+v", subscriptionId, err)
	}
	for _, item := range results.Items {
		if item.Id == nil || item.Properties.HsmUri == nil {
			continue
		}

		// Populate the managed HSM into the cache
		managedHSMID, err := managedhsms.ParseManagedHSMIDInsensitively(*item.Id)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a Managed HSM ID: %+v", *item.Id, err)
		}
		hsmUri := *item.Properties.HsmUri
		c.AddManagedHSMToCache(*managedHSMID, hsmUri)
	}

	// Now that the cache has been repopulated, check if we have the managed HSM or not
	if cached := c.managedHSMCache.getCachedID(*managedHSMName); cached != nil {
		return cached, nil
	}

	// We haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) PurgeManagedHSM(id managedhsms.ManagedHSMId) {
	c.managedHSMCache.purge(managedHSMCacheKey(id))
}
