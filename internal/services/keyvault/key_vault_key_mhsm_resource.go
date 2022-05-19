package keyvault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/Azure/go-autorest/autorest/date"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// mhsm key just like popular vault key, but with different id format
// implement with Resource Interface

type KeyVaultKey struct {
	ID               string                       `tfschema:"id"`
	Name             string                       `tfschema:"name"`
	KeyVaultID       string                       `tfschema:"key_vault_id"`
	KeyType          keyvault.JSONWebKeyType      `tfschema:"key_type"`
	KeySize          int                          `tfschema:"key_size"`
	KeyOpts          []string                     `tfschema:"key_opts"`
	Version          string                       `tfschema:"version"`
	VersionlessID    string                       `tfschema:"versionless_id"`
	Curve            keyvault.JSONWebKeyCurveName `tfschema:"curve"`
	N                string                       `tfschema:"n"`
	E                string                       `tfschema:"e"`
	X                string                       `tfschema:"x"`
	Y                string                       `tfschema:"y"`
	NotBeforDate     string                       `tfschema:"not_before_date"`
	ExpirationDate   string                       `tfschema:"expiration_date"`
	PubilcKeyPem     string                       `tfschema:"public_key_pem"`
	PublicKeyOpenSSH string                       `tfschema:"public_key_openssh"`
	Tags             map[string]*string           `tfschema:"tags"`
}

func (k *KeyVaultKey) setPublicKey(pubilcKey interface{}) error {
	pem, openssh, err := parsePublicKey(pubilcKey)
	if err != nil {
		return fmt.Errorf("failed to read public key: %+v", err)
	}
	k.PubilcKeyPem = pem
	k.PublicKeyOpenSSH = openssh
	return nil
}

type allStr interface {
	~string
}

func convertStrSlice[T allStr](input []string) (res *[]T) {
	var s []T
	for _, v := range input {
		s = append(s, T(v))
	}
	return &s
}

func revertToStrSlice[T allStr](input *[]T) (res []string) {
	if input == nil {
		return nil
	}
	for _, v := range *input {
		res = append(res, string(v))
	}
	return
}

type pubKey interface {
	rsa.PublicKey | ecdsa.PublicKey
}

func parsePublicKey(pubKey interface{}) (pemStr, sshKeyStr string, err error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		err = fmt.Errorf("failed to marshal public key error: %s", err)
		return
	}
	pubKeyPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	pemStr = string(pem.EncodeToMemory(pubKeyPemBlock))

	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err == nil {
		// Not all EC types can be SSH keys, so we'll produce this only
		// if an appropriate type was selected.
		sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
		sshKeyStr = string(sshPubKeyBytes)
	}
	return
}

type MHSMKey struct{}

var _ sdk.ResourceWithUpdate = &MHSMKey{}

func (m *MHSMKey) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyVaultValidate.NestedItemName,
		},

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyVaultValidate.ManagedHSMID,
		},

		"key_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// turns out Azure's *really* sensitive about the casing of these
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
			ValidateFunc: validation.StringInSlice([]string{
				// TODO: add `oct` back in once this is fixed
				// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257
				string(keyvault.EC),
				string(keyvault.ECHSM),
				string(keyvault.RSA),
				string(keyvault.RSAHSM),
			}, false),
		},

		"key_size": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"curve"},
		},

		"key_opts": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.Decrypt),
					string(keyvault.Encrypt),
					string(keyvault.Sign),
					string(keyvault.UnwrapKey),
					string(keyvault.Verify),
					string(keyvault.WrapKey),
				}, false),
			},
		},

		"curve": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return old == "SECP256K1" && new == string(keyvault.P256K)
			},
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				out := []string{
					string(keyvault.P256),
					string(keyvault.P256K),
					string(keyvault.P384),
					string(keyvault.P521),
				}
				return validation.StringInSlice(out, false)
			}(),
			// TODO: the curve name should probably be mandatory for EC in the future,
			// but handle the diff so that we don't break existing configurations and
			// imported EC keys
			ConflictsWith: []string{"key_size"},
		},

		"not_before_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"expiration_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"tags": tags.Schema(),
	}
}

func (m *MHSMKey) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"versionless_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"n": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"e": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"x": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"y": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_key_pem": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_key_openssh": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m *MHSMKey) ModelObject() interface{} {
	return &KeyVaultKey{}
}

func (m *MHSMKey) ResourceType() string {
	return "azurerm_managed_hsm_key"
}

func (m *MHSMKey) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			name := meta.ResourceData.Get("name").(string)
			id, err := parse.ManagedHSMID(meta.ResourceData.Get("key_vault_id").(string))
			if err != nil {
				return nil
			}
			cli := meta.Client.KeyVault
			//hsmCli := cli.ManagedHsmClient
			baseUri, err := cli.BaseUriForKeyVaultV2(ctx, *id)
			if err != nil {
				return fmt.Errorf("looking up Key %q managed hsm from id %q: %+v", name, id, err)
			}
			client := cli.MHSMManagementClient
			existing, err := client.GetKey(ctx, *baseUri, name, "")
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing key from %s (MHSM vault %q): %s", name, baseUri, err)
				}
			}

			if existing.Key != nil && existing.Key.Kid != nil && *existing.Key.Kid != "" {
				return tf.ImportAsExistsError("azurerm_key_vault_key", *existing.Key.Kid)
			}
			// code from key_vault_key_resource
			var model KeyVaultKey
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}
			d := meta.ResourceData
			//keyType := d.Get("key_type").(string)
			//keyOptions := expandKeyVaultKeyOptions(d)
			//t := d.Get("tags").(map[string]interface{})

			// TODO: support Importing Keys once this is fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/1747
			parameters := keyvault.KeyCreateParameters{
				Kty:    keyvault.JSONWebKeyType(model.KeyType),
				KeyOps: convertStrSlice[keyvault.JSONWebKeyOperation](model.KeyOpts),
				KeyAttributes: &keyvault.KeyAttributes{
					Enabled: utils.Bool(true),
				},

				Tags: model.Tags,
			}

			if parameters.Kty == keyvault.EC || parameters.Kty == keyvault.ECHSM {
				parameters.Curve = model.Curve
			} else if parameters.Kty == keyvault.RSA || parameters.Kty == keyvault.RSAHSM {
				if model.KeySize == 0 {
					return fmt.Errorf("Key size is required when creating an RSA key")
				}
				parameters.KeySize = utils.Int32(int32(model.KeySize))
			}
			// TODO: support `oct` once this is fixed
			// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257

			if model.NotBeforDate != "" {
				notBeforeDate, _ := time.Parse(time.RFC3339, model.NotBeforDate) // validated by schema
				notBeforeUnixTime := date.UnixTime(notBeforeDate)
				parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
			}

			if model.ExpirationDate != "" {
				expirationDate, _ := time.Parse(time.RFC3339, model.ExpirationDate) // validated by schema
				expirationUnixTime := date.UnixTime(expirationDate)
				parameters.KeyAttributes.Expires = &expirationUnixTime
			}
			if resp, err := client.CreateKey(ctx, *baseUri, name, parameters); err != nil {
				if meta.Client.Features.KeyVault.RecoverSoftDeletedKeys && utils.ResponseWasConflict(resp.Response) {
					recoveredKey, err := client.RecoverDeletedKey(ctx, *baseUri, name)
					if err != nil {
						return err
					}
					log.Printf("[DEBUG] Recovering Key %q with ID: %q", name, *recoveredKey.Key.Kid)
					if kid := recoveredKey.Key.Kid; kid != nil {
						stateConf := &pluginsdk.StateChangeConf{
							Pending:                   []string{"pending"},
							Target:                    []string{"available"},
							Refresh:                   keyVaultChildItemRefreshFunc(*kid),
							Delay:                     30 * time.Second,
							PollInterval:              10 * time.Second,
							ContinuousTargetOccurence: 10,
							Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
						}

						if _, err := stateConf.WaitForStateContext(ctx); err != nil {
							return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
						}
						log.Printf("[DEBUG] Key %q recovered with ID: %q", name, *kid)
					}
				} else {
					return fmt.Errorf("Creating Key: %+v", err)
				}
			}

			// "" indicates the latest version
			read, err := client.GetKey(ctx, *baseUri, name, "")
			if err != nil {
				return err
			}

			d.SetId(*read.Key.Kid)

			return m.Read().Func(ctx, meta)
		},
	}
}

func (m *MHSMKey) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			id, err := parse.ParseNestedItemID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			keyVaultIdRaw, err := meta.Client.KeyVault.KeyVaultIDFromBaseUrl(ctx, meta.Client.Resource, id.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
			}
			state := KeyVaultKey{
				Name:          id.Name,
				Version:       id.Version,
				VersionlessID: id.VersionlessID(),
			}
			defer func() {
				if err == nil {
					err = meta.Encode(&state)
				}
			}()
			if keyVaultIdRaw == nil {
				log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
				state.ID = ""
				return nil
			}
			hsmID, err := parse.ManagedHSMID(*keyVaultIdRaw)
			if err != nil {
				return err
			}
			ok, err := meta.Client.KeyVault.ExistsV2(ctx, hsmID)
			if err != nil {
				return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v",
					*hsmID, id.Name, id.KeyVaultBaseUrl, err)
			}
			if !ok {
				log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state",
					id.Name, *hsmID, id.KeyVaultBaseUrl)
				return nil
			}

			resp, err := meta.Client.KeyVault.MHSMManagementClient.GetKey(ctx, id.KeyVaultBaseUrl, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					log.Printf("[DEBUG] Key %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
					return nil
				}
				return err
			}
			str := func(ptr *string) string {
				if ptr != nil {
					return *ptr
				}
				return ""
			}
			state.Name = id.Name
			if key := resp.Key; key != nil {
				state.KeyType = key.Kty
				if key.KeyOps != nil {
					state.KeyOpts = *key.KeyOps
				}
				state.E = str(key.E)
				state.X = str(key.X)
				state.Y = str(key.Y)
				if key.N != nil {
					nBytes, err := base64.RawStdEncoding.DecodeString(*key.N)
					if err != nil {
						return fmt.Errorf("Could not decode N: %+v", err)
					}
					state.N = *key.N
					state.KeySize = len(nBytes) * 8
				}
				state.Curve = key.Crv
				// convert pem and openssh key
				var publicKey interface{}
				switch key.Kty {
				case keyvault.RSA, keyvault.RSAHSM:
					nBytes, err := base64.RawURLEncoding.DecodeString(state.N)
					if err != nil {
						return fmt.Errorf("failed to decode N: %+v", err)
					}
					eBytes, err := base64.RawURLEncoding.DecodeString(state.E)
					if err != nil {
						return fmt.Errorf("failed to decode E: %+v", err)
					}
					publicKey = &rsa.PublicKey{
						N: big.NewInt(0).SetBytes(nBytes),
						E: int(big.NewInt(0).SetBytes(eBytes).Uint64()),
					}
				case keyvault.EC, keyvault.ECHSM:
					xBytes, err := base64.RawURLEncoding.DecodeString(state.X)
					if err != nil {
						return fmt.Errorf("failed to decode X: %+v", err)
					}
					yBytes, err := base64.RawURLEncoding.DecodeString(state.Y)
					if err != nil {
						return fmt.Errorf("failed to decode Y: %+v", err)
					}
					publicKey := &ecdsa.PublicKey{
						X: big.NewInt(0).SetBytes(xBytes),
						Y: big.NewInt(0).SetBytes(yBytes),
					}
					switch key.Crv {
					case keyvault.P256:
						publicKey.Curve = elliptic.P256()
					case keyvault.P384:
						publicKey.Curve = elliptic.P384()
					case keyvault.P521:
						publicKey.Curve = elliptic.P521()
					}
				}
				// read publicKey info
				if err := state.setPublicKey(publicKey); err != nil {
					return err
				}
			}
			if attr := resp.Attributes; attr != nil {
				timeStr := func(ptr *date.UnixTime) string {
					if ptr != nil {
						return time.Time(*ptr).Format(time.RFC3339)
					}
					return ""
				}
				state.NotBeforDate = timeStr(attr.NotBefore)
				state.ExpirationDate = timeStr(attr.Expires)
			}
			state.Version = id.Version
			state.VersionlessID = id.VersionlessID()
			// set tags
			state.Tags = resp.Tags
			return nil
		},
	}
}

func (m *MHSMKey) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (m *MHSMKey) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (m *MHSMKey) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, s string) ([]string, []error) {
		_, err := parse.ParseNestedItemID(i.(string))
		return nil, []error{err}
	}
}
