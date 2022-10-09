package keyvault

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/sdk/v7.3/keyvault"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Certificate struct {
	CertificateId string `tfschema:"certificate_id"`
	Alg           string `tfschema:"alg"`
}

type KeyVaultSecurityDomainDownloadModel struct {
	VaultBaseUrl string        `tfschema:"vault_base_url"`
	Certificate  []Certificate `tfschema:"certificate"`
	Required     int           `tfschema:"required"`
	EncData      string        `tfschema:"enc_data"`
}

type KeyVaultSecurityDomainDownloadResource struct{}

var _ sdk.Resource = (*KeyVaultSecurityDomainDownloadResource)(nil)

func (m KeyVaultSecurityDomainDownloadResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"vault_base_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"certificate": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 3,
			MaxItems: 10,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.NestedItemId,
					},

					"alg": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(func() (res []string) {
							for _, v := range keyvault.PossibleJSONWebKeyEncryptionAlgorithmValues() {
								res = append(res, string(v))
							}
							return
						}(), false),
					},
				},
			},
		},

		"required": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      2,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 10),
		},
	}
}

func (m KeyVaultSecurityDomainDownloadResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"enc_data": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m KeyVaultSecurityDomainDownloadResource) ModelObject() interface{} {
	return &KeyVaultSecurityDomainDownloadModel{}
}

func (m KeyVaultSecurityDomainDownloadResource) ResourceType() string {
	return "azurerm_key_vault_security_domain_download"
}

func (m KeyVaultSecurityDomainDownloadResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.KeyVault.MHSMSDClient
			keyClient := meta.Client.KeyVault.ManagementClient

			var model KeyVaultSecurityDomainDownloadModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			//subscriptionID := meta.Client.Account.SubscriptionId
			// get keys from key vault key
			var param keyvault.CertificateInfoObject
			param.Required = utils.Int32(int32(model.Required))
			var certs []keyvault.SecurityDomainJSONWebKey
			for _, cert := range model.Certificate {
				keyID, _ := parse.ParseNestedItemID(cert.CertificateId)
				certRes, err := keyClient.GetCertificate(ctx, keyID.KeyVaultBaseUrl, keyID.Name, keyID.Version)
				if err != nil {
					return fmt.Errorf("retriving key %s: %v", cert.CertificateId, err)
				}
				if certRes.Cer == nil {
					return fmt.Errorf("got nil key for %s", cert.CertificateId)
				}
				cert := keyvault.SecurityDomainJSONWebKey{
					Kty:    pointer.FromString("RSA"),
					KeyOps: &[]string{""},
					Alg:    pointer.FromString(cert.Alg),
				}
				if certRes.Policy != nil && certRes.Policy.KeyProperties != nil {
					cert.Kty = pointer.FromString(string(certRes.Policy.KeyProperties.KeyType))
				}
				x5c := ""
				if contents := certRes.Cer; contents != nil {
					x5c = base64.StdEncoding.EncodeToString(*contents)
				}
				cert.X5c = &[]string{x5c}

				sum1 := sha1.Sum([]byte(x5c))
				x5tDst := make([]byte, base64.StdEncoding.EncodedLen(len(sum1)))
				base64.URLEncoding.Encode(x5tDst, sum1[:])
				cert.X5t = pointer.FromString(string(x5tDst))

				sum256 := sha256.Sum256([]byte(x5c))
				s256Dst := make([]byte, base64.StdEncoding.EncodedLen(len(sum256)))
				base64.URLEncoding.Encode(s256Dst, sum256[:])
				cert.X5tS256 = pointer.FromString(string(s256Dst))
				certs = append(certs, cert)
			}
			param.Certificates = &certs

			future, err := client.Download(ctx, model.VaultBaseUrl, param)
			originResponse := future.Response()
			data, err := io.ReadAll(originResponse.Body)
			if err != nil {
				return err
			}
			var EncData struct {
				Value string `json:"value"`
			}
			err = json.Unmarshal(data, &EncData)
			if err != nil {
				return err
			}
			if err != nil {
				return fmt.Errorf("downloading %s: %v", model.VaultBaseUrl, err)
			}
			// wait download code has bug will never return
			// limit ctx to wait 5 second(value from azcli)
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasStatusCode(future.Response(), http.StatusOK) {
					return fmt.Errorf("waiting for download of %s: %v", model.VaultBaseUrl, err)
				}
			}
			result, err := future.Result(*client)
			if result.Value != nil {
				EncData.Value = pointer.ToString(result.Value)
			}

			// set id
			//meta.SetID()
			return meta.Encode(result)
		},
	}
}

func (m KeyVaultSecurityDomainDownloadResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			var model KeyVaultSecurityDomainDownloadModel
			if err := meta.Decode(&model); err != nil {
				return err
			}
			return meta.Encode(&model)
		},
	}
}

func (m KeyVaultSecurityDomainDownloadResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (m KeyVaultSecurityDomainDownloadResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, s string) ([]string, []error) {
		return nil, nil
	}
}
