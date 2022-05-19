package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

type notFoundError string

func (n notFoundError) Error() string {
	return string(n)
}

func newNotFoundError(id interface{}) notFoundError {
	n := fmt.Sprintf("%s was not found", id)
	return notFoundError(n)
}

func (c *Client) BaseUriForKeyVaultV2(ctx context.Context, id parse.IVaultID) (*string, error) {
	defer c.lockKey(id.GetKey()).Unlock()
	subscription := id.GetSubscription()
	var e error
	var res autorest.Response
	if parse.IsHSM(id) {
		resp, err := c.GetHSMClient(subscription).Get(ctx, id.GetResourceGroup(), id.GetName())
		if resp.Properties != nil && resp.Properties.HsmURI != nil {
			return resp.Properties.HsmURI, nil
		}
		e, res = err, resp.Response
	} else {
		resp, err := c.GetVaultClient(subscription).Get(ctx, id.GetResourceGroup(), id.GetName())
		if resp.Properties != nil && resp.Properties.VaultURI != nil {
			return resp.Properties.VaultURI, nil
		}
		e, res = err, resp.Response
	}
	if e != nil {
		if utils.ResponseWasNotFound(res) {
			return nil, newNotFoundError(id)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, e)
	}
	return nil, fmt.Errorf("`properties` was nil for %s", id)
}

func (c *Client) ExistsV2(ctx context.Context, id parse.IVaultID) (bool, error) {
	uri, err := c.BaseUriForKeyVaultV2(ctx, id)
	if err != nil || uri == nil {
		if _, ok := err.(notFoundError); ok {
			return false, nil
		}
		return false, err
	}
	c.AddToCacheV2(id, *uri)
	return true, nil
}
