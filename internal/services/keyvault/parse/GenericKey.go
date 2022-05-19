package parse

// IVaultID defiinterface for vaultid and mhsmid process
type IVaultID interface {
	ID() string
	GetKey() string
	GetName() string
	GetSubscription() string
	GetResourceGroup() string
}

var _ IVaultID = new(VaultId)
var _ IVaultID = new(ManagedHSMId)

func IsHSM(id IVaultID) bool {
	switch id.(type) {
	case ManagedHSMId, *ManagedHSMId:
		return true
	}
	return false
}

func (id VaultId) GetKey() string {
	return "vault-" + id.Name
}

func (id VaultId) GetName() string {
	return id.Name
}

func (id VaultId) GetSubscription() string {
	return id.SubscriptionId
}

func (id VaultId) GetResourceGroup() string {
	return id.ResourceGroup
}

func (id ManagedHSMId) GetKey() string {
	return "hsm-" + id.Name
}

func (id ManagedHSMId) GetName() string {
	return id.Name
}

func (id ManagedHSMId) GetSubscription() string {
	return id.SubscriptionId
}

func (id ManagedHSMId) GetResourceGroup() string {
	return id.ResourceGroup
}
