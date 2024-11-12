// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/customermanagedkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	hsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppAccountEncryptionResource struct{}

var _ sdk.Resource = NetAppAccountEncryptionResource{}

func (r NetAppAccountEncryptionResource) ModelObject() interface{} {
	return &netAppModels.NetAppAccountEncryption{}
}

func (r NetAppAccountEncryptionResource) ResourceType() string {
	return "azurerm_netapp_account_encryption"
}

func (r NetAppAccountEncryptionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return netappaccounts.ValidateNetAppAccountID
}

func (r NetAppAccountEncryptionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"netapp_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the NetApp Account where encryption will be set.",
			ValidateFunc: netAppValidate.ValidateNetAppAccountID,
		},

		"user_assigned_identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			Description:   "The resource ID of the User Assigned Identity to use for encryption.",
			ConflictsWith: []string{"system_assigned_identity_principal_id"},
		},

		"system_assigned_identity_principal_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			Description:   "The Principal ID of the System Assigned Identity to use for encryption.",
			ConflictsWith: []string{"user_assigned_identity_id"},
		},

		"encryption_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
			ExactlyOneOf: []string{"encryption_managed_hsm_key"},
			Description:  "The versionless encryption key url.",
		},

		"encryption_managed_hsm_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: hsmValidate.ManagedHSMDataPlaneVersionlessKeyID,
			ExactlyOneOf: []string{"encryption_key"},
			Description:  "The versionless managed HSM key id.",
		},
	}
}

func (r NetAppAccountEncryptionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppAccountEncryptionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient

			var model netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountID, err := netappaccounts.ParseNetAppAccountID(model.NetAppAccountID)
			if err != nil {
				return fmt.Errorf("error parsing account id %s: %+v", model.NetAppAccountID, err)
			}

			metadata.Logger.Infof("Import check for %s", accountID.ID())

			locks.ByID(accountID.ID())
			defer locks.UnlockByID(accountID.ID())

			existing, err := client.AccountsGet(ctx, pointer.From(accountID))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("not found %s: %s", accountID.ID(), err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				if existing.Model.Properties.Encryption != nil && existing.Model.Properties.Encryption.KeySource != nil && pointer.From(existing.Model.Properties.Encryption.KeySource) == netappaccounts.KeySourceMicrosoftPointKeyVault {

					return tf.ImportAsExistsError(r.ResourceType(), accountID.ID())
				}
			}

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			encryptionExpanded, err := expandEncryption(ctx, metadata.ResourceData, metadata.Client)
			if err != nil {
				return err
			}

			update.Properties.Encryption = encryptionExpanded

			if err := client.AccountsUpdateThenPoll(ctx, pointer.From(accountID), update); err != nil {
				return fmt.Errorf("updating %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)

			return nil
		},
	}
}

func (r NetAppAccountEncryptionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			metadata.Logger.Infof("Updating %s", id)

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_id") ||
				metadata.ResourceData.HasChange("system_assigned_identity_principal_id") ||
				metadata.ResourceData.HasChange("encryption_key") ||
				metadata.ResourceData.HasChange("encryption_managed_hsm_key") {

				encryptionExpanded, err := expandEncryption(ctx, metadata.ResourceData, metadata.Client)
				if err != nil {
					return err
				}
				update.Properties.Encryption = encryptionExpanded

				if err := metadata.Client.NetApp.AccountClient.AccountsUpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				metadata.SetID(id)
			}

			return nil
		},
	}
}

func (r NetAppAccountEncryptionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.AccountClient

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.AccountsGet(ctx, pointer.From(id))
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model.Properties.Encryption == nil {
				return fmt.Errorf("encryption information does not exist for %s", id)
			}

			anfAccountIdentityFlattened, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(existing.Model.Identity)
			if err != nil {
				return err
			}

			model := netAppModels.NetAppAccountEncryption{
				NetAppAccountID: id.ID(),
			}

			encryptionKey, err := flattenEncryption(existing.Model.Properties.Encryption, metadata.Client.Account.Environment.ManagedHSM)
			if err != nil {
				return err
			}

			if encryptionKey != nil {
				if encryptionKey.KeyVaultKeyId != nil {
					model.EncryptionKey = encryptionKey.KeyVaultKeyId.VersionlessID()
				} else if encryptionKey.ManagedHSMKeyVersionlessId != nil {
					model.EncryptionManagedHSMKey = encryptionKey.ManagedHSMKeyVersionlessId.ID()
				}
			}

			if len(anfAccountIdentityFlattened) > 0 {

				if anfAccountIdentityFlattened[0].Type == identity.TypeSystemAssigned {
					model.SystemAssignedIdentityPrincipalID = anfAccountIdentityFlattened[0].PrincipalId
				}

				if anfAccountIdentityFlattened[0].Type == identity.TypeUserAssigned {
					if len(anfAccountIdentityFlattened[0].IdentityIds) > 0 {
						model.UserAssignedIdentityID = anfAccountIdentityFlattened[0].IdentityIds[0]
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppAccountEncryptionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			update.Properties.Encryption = &netappaccounts.AccountEncryption{}
			// ATTENTION: This actually cannot remove the encryption and nothing is updated in the service.
			// But it seems just fine to terraform as the encryption is removed from the state.
			if err := client.AccountsUpdateThenPoll(ctx, pointer.From(id), update); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandEncryption(ctx context.Context, d *pluginsdk.ResourceData, clientHub *clients.Client) (*netappaccounts.AccountEncryption, error) {
	cmk, err := customermanagedkeys.ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d,
		customermanagedkeys.VersionTypeVersionless, "encryption_key", "encryption_managed_hsm_key",
		clientHub.Account.Environment.KeyVault, clientHub.Account.Environment.ManagedHSM)

	if err != nil {
		return nil, fmt.Errorf("expanding customermanagedkeys: %+v", err)
	}

	encryptionProperty := netappaccounts.AccountEncryption{
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointNetApp),
	}

	if cmk == nil {
		return &encryptionProperty, nil
	}

	encryptionProperty.KeySource = pointer.To(netappaccounts.KeySourceMicrosoftPointKeyVault)
	if userIdentity := d.Get("user_assigned_identity_id").(string); userIdentity != "" {
		encryptionProperty.Identity = &netappaccounts.EncryptionIdentity{
			UserAssignedIdentity: pointer.To(userIdentity),
		}
	}

	subscriptionID := commonids.NewSubscriptionID(clientHub.Account.SubscriptionId)
	encryptionProperty.KeyVaultProperties = &netappaccounts.KeyVaultProperties{
		KeyVaultUri: cmk.BaseUri(),
	}

	if cmk.KeyVaultKeyId != nil {
		encryptionProperty.KeyVaultProperties.KeyName = cmk.KeyVaultKeyId.Name
		keyVaultID, err := clientHub.KeyVault.KeyVaultIDFromBaseUrl(ctx, subscriptionID, cmk.KeyVaultKeyId.KeyVaultBaseUrl)
		if err != nil {
			return nil, fmt.Errorf("retrieving the resource id the key vault at url %q: %s", cmk.KeyVaultKeyId.KeyVaultBaseUrl, err)
		}
		encryptionProperty.KeyVaultProperties.KeyVaultResourceId = pointer.From(keyVaultID)
	} else if baseUri := cmk.BaseUri(); baseUri != "" {
		encryptionProperty.KeyVaultProperties.KeyName = cmk.ManagedHSMKeyVersionlessId.KeyName
		hsmID, err := clientHub.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionID, baseUri, nil)
		if err != nil {
			return nil, fmt.Errorf("retrieving the resource id the managed hsm at url %q: %s", baseUri, err)
		}
		encryptionProperty.KeyVaultProperties.KeyVaultResourceId = hsmID.ID()
	}
	return &encryptionProperty, nil
}

func flattenEncryption(prop *netappaccounts.AccountEncryption, hsmEnv environments.Api) (*customermanagedkeys.KeyVaultOrManagedHSMKey, error) {
	if prop == nil || prop.KeyVaultProperties == nil || pointer.From(prop.KeySource) == netappaccounts.KeySourceMicrosoftPointNetApp {
		return nil, nil
	}

	cmk, err := customermanagedkeys.FlattenKeyVaultOrManagedHSMIDByComponents(prop.KeyVaultProperties.KeyVaultUri, prop.KeyVaultProperties.KeyName, "", hsmEnv)
	if err != nil {
		return nil, fmt.Errorf("flattening key vault or managed hsm id: %+v", err)
	}
	return cmk, nil
}
