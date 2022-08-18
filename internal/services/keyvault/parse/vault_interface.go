package parse

type Vaulter interface {
	GetSubscriptionID() string
	GetResourceGroup() string
	GetName() string
	GetCacheKey() string
}

var (
	_ Vaulter = VaultId{}
	_ Vaulter = ManagedHSMId{}
)

func (id VaultId) GetSubscriptionID() string {
	return id.SubscriptionId
}

func (id VaultId) GetResourceGroup() string {
	return id.ResourceGroup
}

func (id VaultId) GetName() string {
	return id.Name
}

func (id VaultId) GetCacheKey() string {
	return "keyvault:" + id.GetName()
}

func (id ManagedHSMId) GetSubscriptionID() string {
	return id.SubscriptionId
}

func (id ManagedHSMId) GetResourceGroup() string {
	return id.ResourceGroup
}

func (id ManagedHSMId) GetName() string {
	return id.Name
}

func (id ManagedHSMId) GetCacheKey() string {
	return "mhsm:" + id.GetName()
}
