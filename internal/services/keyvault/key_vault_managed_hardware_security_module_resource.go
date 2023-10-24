// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	kv74 "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultManagedHardwareSecurityModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmKeyVaultManagedHardwareSecurityModuleCreate,
		Read:   resourceArmKeyVaultManagedHardwareSecurityModuleRead,
		Delete: resourceArmKeyVaultManagedHardwareSecurityModuleDelete,
		Update: resourceArmKeyVaultManagedHardwareSecurityModuleUpdate,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedhsms.ParseManagedHSMID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(120 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(120 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(keyVaultHSMCustomizeDiff),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagedHardwareSecurityModuleName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(managedhsms.ManagedHsmSkuNameStandardBOne),
				}, false),
			},

			"admin_object_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"soft_delete_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      90,
				ValidateFunc: validation.IntBetween(7, 90),
			},

			"hsm_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				// Computed: true,
				Default:  true,
				ForceNew: true,
			},

			// replication has to after hsm activated
			// or error like: Security domain is not downloaded for the pool
			"replication_regions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:             pluginsdk.TypeString,
					ValidateFunc:     validation.StringIsNotEmpty,
					DiffSuppressFunc: location.DiffSuppressFunc,
				},
			},

			"network_acls": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedhsms.NetworkRuleActionAllow),
								string(managedhsms.NetworkRuleActionDeny),
							}, false),
						},
						"bypass": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedhsms.NetworkRuleBypassOptionsNone),
								string(managedhsms.NetworkRuleBypassOptionsAzureServices),
							}, false),
						},
					},
				},
			},

			"security_domain_key_vault_certificate_ids": {
				Type:         pluginsdk.TypeList,
				MinItems:     3,
				MaxItems:     10,
				Optional:     true,
				RequiredWith: []string{"security_domain_quorum"},
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.NestedItemId,
				},
			},

			"security_domain_quorum": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				RequiredWith: []string{"security_domain_key_vault_certificate_ids"},
				ValidateFunc: validation.IntBetween(2, 10),
			},

			"security_domain_encrypted_data": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			// https://github.com/Azure/azure-rest-api-specs/issues/13365
			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceArmKeyVaultManagedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	kvClient := meta.(*clients.Client).KeyVault
	hsmClient := kvClient.ManagedHsmClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managedhsms.NewManagedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := hsmClient.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_key_vault_managed_hardware_security_module", id.ID())
	}

	publicNetworkAccessEnabled := managedhsms.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccessEnabled = managedhsms.PublicNetworkAccessDisabled
	}
	hsm := managedhsms.ManagedHsm{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &managedhsms.ManagedHsmProperties{
			InitialAdminObjectIds:     utils.ExpandStringSlice(d.Get("admin_object_ids").(*pluginsdk.Set).List()),
			CreateMode:                pointer.To(managedhsms.CreateModeDefault),
			EnableSoftDelete:          utils.Bool(true),
			SoftDeleteRetentionInDays: utils.Int64(int64(d.Get("soft_delete_retention_days").(int))),
			EnablePurgeProtection:     utils.Bool(d.Get("purge_protection_enabled").(bool)),
			PublicNetworkAccess:       pointer.To(publicNetworkAccessEnabled),
			NetworkAcls:               expandMHSMNetworkAcls(d.Get("network_acls").([]interface{})),
		},
		Sku: &managedhsms.ManagedHsmSku{
			Family: managedhsms.ManagedHsmSkuFamilyB,
			Name:   managedhsms.ManagedHsmSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if tenantId := d.Get("tenant_id").(string); tenantId != "" {
		hsm.Properties.TenantId = pointer.To(tenantId)
	}

	if err := hsmClient.CreateOrUpdateThenPoll(ctx, id, hsm); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// security domain download to activate this module
	if ok := d.HasChange("security_domain_key_vault_certificate_ids"); ok {
		// get hsm uri
		resp, err := hsmClient.Get(ctx, id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
			return fmt.Errorf("got nil HSMUri for %s: %+v", id, err)
		}
		encData, err := securityDomainDownload(ctx, kvClient, id, *resp.Model.Properties.HsmUri,
			d.Get("security_domain_key_vault_certificate_ids").([]interface{}),
			d.Get("security_domain_quorum").(int))
		if err != nil {
			return fmt.Errorf("downloading security domain for %q: %+v", id, err)
		}
		d.Set("security_domain_encrypted_data", encData)
	}

	// add regions after Security domain is downloaded for the pool
	if replications := d.Get("replication_regions").([]interface{}); len(replications) > 0 {
		hsm.Properties.Regions = expandMHSMRegions(replications)
		if err := hsmClient.CreateOrUpdateThenPoll(ctx, id, hsm); err != nil {
			return fmt.Errorf("adding replication regions for %s: %+v", id, err)
		}
	}

	return resourceArmKeyVaultManagedHardwareSecurityModuleRead(d, meta)
}

// update to re-activate the security module
func resourceArmKeyVaultManagedHardwareSecurityModuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	kvClient := meta.(*clients.Client).KeyVault
	hsmClient := kvClient.ManagedHsmClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := hsmClient.Get(ctx, *id)
	if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// security domain download to activate this module
	if ok := d.HasChange("security_domain_key_vault_certificate_ids"); ok {
		// get hsm uri
		resp, err := hsmClient.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.HsmUri == nil {
			return fmt.Errorf("got nil HSMUri for %s: %+v", id, err)
		}
		encData, err := securityDomainDownload(ctx, kvClient, *id, *resp.Model.Properties.HsmUri,
			d.Get("security_domain_key_vault_certificate_ids").([]interface{}),
			d.Get("security_domain_quorum").(int))
		if err != nil {
			return fmt.Errorf("downloading security domain for %q: %+v", id, err)
		}
		d.Set("security_domain_encrypted_data", encData)
	}

	// NOTE: cannot removing and adding regions at the same time.
	// add regions after Security domain is downloaded for the pool
	if d.HasChange("replication_regions") {
		hsm := *resp.Model
		hsm.Properties.Regions = expandMHSMRegions(d.Get("replication_regions").([]interface{}))
		if err := hsmClient.CreateOrUpdateThenPoll(ctx, *id, *resp.Model); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return nil
}

func resourceArmKeyVaultManagedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	hsmClient := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := hsmClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[ERROR] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ManagedHSMName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			tenantId := ""
			if props.TenantId != nil {
				tenantId = *props.TenantId
			}
			d.Set("tenant_id", tenantId)
			d.Set("admin_object_ids", utils.FlattenStringSlice(props.InitialAdminObjectIds))
			d.Set("hsm_uri", props.HsmUri)
			d.Set("soft_delete_retention_days", props.SoftDeleteRetentionInDays)
			d.Set("purge_protection_enabled", props.EnablePurgeProtection)

			publicAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess != managedhsms.PublicNetworkAccessEnabled {
				publicAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicAccessEnabled)

			if err := d.Set("network_acls", flattenMHSMNetworkAcls(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_acls`: %+v", err)
			}

			d.Set("replication_regions", flattenMHSMRegions(props.Regions))
		}

		skuName := ""
		if sku := model.Sku; sku != nil {
			skuName = string(sku.Name)
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceArmKeyVaultManagedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	hsmClient := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	// We need to grab the keyvault hsm to see if purge protection is enabled prior to deletion
	resp, err := hsmClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	loc := ""
	purgeProtectionEnabled := false
	if model := resp.Model; model != nil {
		loc = location.NormalizeNilable(model.Location)
		if props := model.Properties; props != nil {
			if props.EnablePurgeProtection != nil {
				purgeProtectionEnabled = *props.EnablePurgeProtection
			}

			if props.Regions != nil && len(*props.Regions) > 0 {
				// Have to remove all replication regions before delete the managed HSM resource
				// https://learn.microsoft.com/en-us/azure/key-vault/managed-hsm/multi-region-replication#soft-delete-behavior
				props.Regions = pointer.To([]managedhsms.MHSMGeoReplicatedRegion{})
				if err = hsmClient.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
					return fmt.Errorf("deleting replication regions for %s: %+v", id, err)
				}
			}
		}
	}

	if err := hsmClient.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedHSMsOnDestroy {
		if purgeProtectionEnabled {
			log.Printf("[DEBUG] cannot purge %s because purge protection is enabled", id)
			return nil
		}
	}

	purgedId := managedhsms.NewDeletedManagedHSMID(id.SubscriptionId, loc, id.ManagedHSMName)
	if err := hsmClient.PurgeDeletedThenPoll(ctx, purgedId); err != nil {
		return fmt.Errorf("purging %s: %+v", id, err)
	}

	return nil
}

func expandMHSMNetworkAcls(input []interface{}) *managedhsms.MHSMNetworkRuleSet {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	res := &managedhsms.MHSMNetworkRuleSet{
		Bypass:        pointer.To(managedhsms.NetworkRuleBypassOptions(v["bypass"].(string))),
		DefaultAction: pointer.To(managedhsms.NetworkRuleAction(v["default_action"].(string))),
	}

	return res
}

func expandMHSMRegions(regions []interface{}) *[]managedhsms.MHSMGeoReplicatedRegion {
	var res []managedhsms.MHSMGeoReplicatedRegion
	for _, v := range regions {
		res = append(res, managedhsms.MHSMGeoReplicatedRegion{
			Name: pointer.To(v.(string)),
		})
	}

	return &res
}

func flattenMHSMNetworkAcls(acl *managedhsms.MHSMNetworkRuleSet) []interface{} {
	bypass := string(managedhsms.NetworkRuleBypassOptionsAzureServices)
	defaultAction := string(managedhsms.NetworkRuleActionAllow)

	if acl != nil {
		if acl.Bypass != nil {
			bypass = string(*acl.Bypass)
		}
		if acl.DefaultAction != nil {
			defaultAction = string(*acl.DefaultAction)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"bypass":         bypass,
			"default_action": defaultAction,
		},
	}
}

func flattenMHSMRegions(regions *[]managedhsms.MHSMGeoReplicatedRegion) (res []string) {
	res = make([]string, 0)
	if regions == nil {
		return
	}

	for _, region := range *regions {
		if !pointer.From(region.IsPrimary) {
			res = append(res, pointer.From(region.Name))
		}
	}

	return res
}

func securityDomainDownload(ctx context.Context, cli *client.Client, id managedhsms.ManagedHSMId, vaultBaseUrl string, certIds []interface{}, quorum int) (encDataStr string, err error) {
	sdClient := cli.MHSMSDClient
	keyClient := cli.ManagementClient

	var param kv74.CertificateInfoObject

	param.Required = utils.Int32(int32(quorum))
	var certs []kv74.SecurityDomainJSONWebKey
	for _, certID := range certIds {
		certIDStr, ok := certID.(string)
		if !ok {
			continue
		}
		keyID, _ := parse.ParseNestedItemID(certIDStr)
		certRes, err := keyClient.GetCertificate(ctx, keyID.KeyVaultBaseUrl, keyID.Name, keyID.Version)
		if err != nil {
			return "", fmt.Errorf("retreiving key %s: %v", certID, err)
		}
		if certRes.Cer == nil {
			return "", fmt.Errorf("got nil key for %s", certID)
		}
		cert := kv74.SecurityDomainJSONWebKey{
			Kty:    pointer.FromString("RSA"),
			KeyOps: &[]string{""},
			Alg:    pointer.FromString("RSA-OAEP-256"),
		}
		if certRes.Policy != nil && certRes.Policy.KeyProperties != nil {
			cert.Kty = pointer.FromString(string(certRes.Policy.KeyProperties.KeyType))
		}
		x5c := ""
		if contents := certRes.Cer; contents != nil {
			x5c = base64.StdEncoding.EncodeToString(*contents)
		}
		cert.X5c = &[]string{x5c}

		sum256 := sha256.Sum256([]byte(x5c))
		s256Dst := make([]byte, base64.StdEncoding.EncodedLen(len(sum256)))
		base64.URLEncoding.Encode(s256Dst, sum256[:])
		cert.X5tS256 = pointer.FromString(string(s256Dst))
		certs = append(certs, cert)
	}
	param.Certificates = &certs

	future, err := sdClient.Download(ctx, vaultBaseUrl, param)
	if err != nil {
		return "", fmt.Errorf("downloading for %s: %v", vaultBaseUrl, err)
	}

	originResponse := future.Response()
	data, err := io.ReadAll(originResponse.Body)
	if err != nil {
		return "", err
	}
	var encData struct {
		Value string `json:"value"`
	}

	err = json.Unmarshal(data, &encData)
	if err != nil {
		return "", fmt.Errorf("unmarshal EncData: %v", err)
	}

	pollerType := custompollers.NewHSMDownloadPoller(sdClient, vaultBaseUrl)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return "", fmt.Errorf("waiting for security domain to download: %+v", err)
	}

	// The GET request may delay, so we need to wait for a while
	conf := pluginsdk.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Finish"},
		Refresh: func() (result interface{}, state string, err error) {
			hsm, err := cli.ManagedHsmClient.Get(ctx, id)
			if err != nil {
				return nil, "Pending", err
			}
			if hsm.Model != nil && hsm.Model.Properties != nil && &hsm.Model.Properties.SecurityDomainProperties != nil {
				prop := *hsm.Model.Properties.SecurityDomainProperties
				status := pointer.From(prop.ActivationStatus)
				switch status {
				case managedhsms.ActivationStatusActive:
					return prop, "Finish", nil
				case managedhsms.ActivationStatusFailed:
					return nil, "Pending", fmt.Errorf("security domain download failed: %+v", prop.ActivationStatusMessage)
				}
			}
			return hsm, "Pending", nil
		},
		Timeout:      time.Minute,
		PollInterval: time.Second * 10,
	}
	if _, err = conf.WaitForStateContext(ctx); err != nil {
		return "", fmt.Errorf("waiting for security domain to download finish: %+v", err)
	}

	return encData.Value, err
}

func keyVaultHSMCustomizeDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	if oldVal, newVal := d.GetChange("security_domain_key_vault_certificate_ids"); len(oldVal.([]interface{})) != 0 && len(newVal.([]interface{})) == 0 {
		if err := d.ForceNew("security_domain_key_vault_certificate_ids"); err != nil {
			return err
		}
	}

	if oldVal, newVal := d.GetChange("security_domain_quorum"); oldVal.(int) != 0 && newVal.(int) == 0 {
		if err := d.ForceNew("security_domain_quorum"); err != nil {
			return err
		}
	}

	return nil
}
