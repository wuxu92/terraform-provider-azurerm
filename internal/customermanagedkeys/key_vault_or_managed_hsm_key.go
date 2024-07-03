package customermanagedkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	hsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyVaultOrManagedHSMKey struct {
	KeyVaultKeyId              *parse.NestedItemId
	ManagedHSMKeyId            *hsmParse.ManagedHSMDataPlaneVersionedKeyId
	ManagedHSMKeyVersionlessId *hsmParse.ManagedHSMDataPlaneVersionlessKeyId
}

func (k *KeyVaultOrManagedHSMKey) IsSet() bool {
	return k != nil && (k.KeyVaultKeyId != nil || k.ManagedHSMKeyId != nil || k.ManagedHSMKeyVersionlessId != nil)
}

func (k *KeyVaultOrManagedHSMKey) ID() string {
	if k == nil {
		return ""
	}

	if k.KeyVaultKeyId != nil {
		return k.KeyVaultKeyId.ID()
	}

	if k.ManagedHSMKeyId != nil {
		return k.ManagedHSMKeyId.ID()
	}

	if k.ManagedHSMKeyVersionlessId != nil {
		return k.ManagedHSMKeyVersionlessId.ID()
	}

	return ""
}

func (k *KeyVaultOrManagedHSMKey) KeyVaultKeyID() string {
	if k != nil && k.KeyVaultKeyId != nil {
		return k.KeyVaultKeyId.ID()
	}
	return ""
}

func (k *KeyVaultOrManagedHSMKey) ManagedHSMKeyID() string {
	if k != nil && k.ManagedHSMKeyId != nil {
		return k.ManagedHSMKeyId.ID()
	}

	if k != nil && k.ManagedHSMKeyVersionlessId != nil {
		return k.ManagedHSMKeyVersionlessId.ID()
	}

	return ""
}

func (k *KeyVaultOrManagedHSMKey) BaseUri() string {
	if k.KeyVaultKeyId != nil {
		return k.KeyVaultKeyId.KeyVaultBaseUrl
	}

	if k.ManagedHSMKeyId != nil {
		return k.ManagedHSMKeyId.BaseUri()
	}

	if k.ManagedHSMKeyVersionlessId != nil {
		return k.ManagedHSMKeyVersionlessId.BaseUri()
	}

	return ""
}

func parseKeyvauleID(keyRaw string, hasVersion *bool) (*parse.NestedItemId, error) {
	keyID, err := parse.ParseOptionallyVersionedNestedKeyID(keyRaw)
	if err != nil {
		return nil, err
	}

	if pointer.From(hasVersion) && keyID.Version == "" {
		return nil, fmt.Errorf("expected a key vault versioned ID but no version information was found in: %q", keyRaw)
	}

	if hasVersion != nil && !*hasVersion && keyID.Version != "" {
		return nil, fmt.Errorf("expected a key vault versionless ID but version information was found in: %q", keyRaw)
	}

	return keyID, nil
}

func parseManagedHSMKey(keyRaw string, hasVersion *bool, hsmEnv environments.Api) (*hsmParse.ManagedHSMDataPlaneVersionedKeyId, *hsmParse.ManagedHSMDataPlaneVersionlessKeyId, error) {
	// if specified with hasVersion == True, then it has to be parsed as versionedKeyID
	var domainSuffix *string
	if hsmEnv != nil {
		domainSuffix, _ = hsmEnv.DomainSuffix()
	}
	// versioned or optional version
	if hasVersion == nil || pointer.From(hasVersion) {
		versioned, err := hsmParse.ManagedHSMDataPlaneVersionedKeyID(keyRaw, domainSuffix)
		if err == nil {
			return versioned, nil, nil
		}
		// if required versioned but got error
		if pointer.From(hasVersion) {
			return nil, nil, err
		}
	}

	// versionless or optional version
	if versionless, err := hsmParse.ManagedHSMDataPlaneVersionlessKeyID(keyRaw, domainSuffix); err == nil {
		return nil, versionless, nil
	} else {
		return nil, nil, err
	}
}

func ExpandKeyVaultOrManagedHSMKey(d interface{}, hasVersion *bool, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	return ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d, hasVersion, "key_vault_key_id", "managed_hsm_key_id", hsmEnv)
}

// ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey
// d: should be one of *pluginsdk.ResourceData or map[string]interface{}
// hasVersion:
//   - nil: both versioned or versionless are ok
//   - true: must have version
//   - false: must not have vesrion
//
// if return nil, nil, it means no key_vault_key_id or managed_hsm_key_id is specified
func ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d interface{}, hasVersion *bool, keyVaultFieldName, hsmFieldName string, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	key := &KeyVaultOrManagedHSMKey{}
	var err error
	var vaultKeyStr, hsmKeyStr string
	if rd, ok := d.(*pluginsdk.ResourceData); ok {
		if keyRaw, ok := rd.GetOk(keyVaultFieldName); ok {
			vaultKeyStr = keyRaw.(string)
		} else if keyRaw, ok = rd.GetOk(hsmFieldName); ok {
			hsmKeyStr = keyRaw.(string)
		}
	} else if obj, ok := d.(map[string]interface{}); ok {
		if keyRaw, ok := obj[keyVaultFieldName]; ok {
			vaultKeyStr, _ = keyRaw.(string)
		}
		if keyRaw, ok := obj[hsmFieldName]; ok {
			hsmKeyStr, _ = keyRaw.(string)
		}
	} else {
		return nil, fmt.Errorf("not supported data type to parse CMK: %T", d)
	}

	switch {
	case vaultKeyStr != "":
		if key.KeyVaultKeyId, err = parseKeyvauleID(vaultKeyStr, hasVersion); err != nil {
			return nil, err
		}
	case hsmKeyStr != "":
		if key.ManagedHSMKeyId, key.ManagedHSMKeyVersionlessId, err = parseManagedHSMKey(hsmKeyStr, hasVersion, hsmEnv); err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	return key, err
}

// FlattenKeyVaultOrManagedHSMID uses `KeyVaultOrManagedHSMKey.SetState()` to save the state, which this function is designed not to do.
func FlattenKeyVaultOrManagedHSMID(id string, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	if id == "" {
		return nil, nil
	}

	key := &KeyVaultOrManagedHSMKey{}
	var err error
	key.KeyVaultKeyId, err = parse.ParseOptionallyVersionedNestedKeyID(id)
	if err == nil {
		return key, nil
	}

	var domainSuffix *string
	if hsmEnv != nil {
		domainSuffix, _ = hsmEnv.DomainSuffix()
	}
	if key.ManagedHSMKeyId, err = hsmParse.ManagedHSMDataPlaneVersionedKeyID(id, domainSuffix); err == nil {
		return key, nil
	}

	if key.ManagedHSMKeyVersionlessId, err = hsmParse.ManagedHSMDataPlaneVersionlessKeyID(id, domainSuffix); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("cannot parse given id to key vault key nor managed hsm key: %s", id)
}
