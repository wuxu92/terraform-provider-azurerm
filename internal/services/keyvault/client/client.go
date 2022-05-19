package client

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault"
	"sync"

	keyvaultmgmt "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ManagedHsmClient     *keyvault.ManagedHsmsClient
	ManagementClient     *keyvaultmgmt.BaseClient
	MHSMManagementClient *keyvaultmgmt.BaseClient
	VaultsClient         *keyvault.VaultsClient
	options              *common.ClientOptions

	lock     *sync.RWMutex
	keyLocks map[string]*sync.RWMutex
	keyCache map[string]keyVaultDetails // vault key cache both keyvaultkey and hsmkey
}

func NewClient(o *common.ClientOptions) *Client {
	managedHsmClient := keyvault.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedHsmClient.Client, o.ResourceManagerAuthorizer)

	managementClient := keyvaultmgmt.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	mhsmManagementClient := keyvaultmgmt.New()
	o.ConfigureClient(&mhsmManagementClient.Client, o.MangedHSMAuthorizer)

	vaultsClient := keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedHsmClient:     &managedHsmClient,
		ManagementClient:     &managementClient,
		MHSMManagementClient: &mhsmManagementClient,
		VaultsClient:         &vaultsClient,
		options:              o,

		lock:     &sync.RWMutex{},
		keyLocks: map[string]*sync.RWMutex{},
		keyCache: map[string]keyVaultDetails{},
	}
}

type Unlocker interface {
	Unlock()
}

func (client Client) lockKey(key string) Unlocker {
	client.lock.RLock()
	l := client.keyLocks[key]
	client.lock.RUnlock()
	if l == nil {
		client.lock.Lock()
		l = &sync.RWMutex{}
		client.keyLocks[key] = l
		client.lock.Unlock()
	}
	l.Lock()
	return l
}

func (c *Client) GetVaultClient(subscription string) *keyvault.VaultsClient {
	if subscription != "" || subscription != c.VaultsClient.SubscriptionID {
		c.VaultsClient = c.KeyVaultClientForSubscription(subscription)
	}
	return c.VaultsClient
}

func (c *Client) GetHSMClient(subscription string) *keyvault.ManagedHsmsClient {
	if subscription != "" && subscription != c.ManagedHsmClient.SubscriptionID {
		cli := keyvault.NewManagedHsmsClient(subscription)
		c.options.ConfigureClient(&cli.Client, c.options.ResourceManagerAuthorizer)
		c.ManagedHsmClient = &cli
	}
	return c.ManagedHsmClient
}

func (client Client) KeyVaultClientForSubscription(subscriptionId string) *keyvault.VaultsClient {
	vaultsClient := keyvault.NewVaultsClientWithBaseURI(client.options.ResourceManagerEndpoint, subscriptionId)
	client.options.ConfigureClient(&vaultsClient.Client, client.options.ResourceManagerAuthorizer)
	return &vaultsClient
}
