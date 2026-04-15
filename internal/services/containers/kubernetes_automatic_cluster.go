// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	automaticClusterEnabledFields = []string{
		"azure_policy_enabled",
		"local_account_disabled",
		"image_cleaner_enabled",
		// "image_cleaner_interval_hours",
		"oidc_issuer_enabled",
		"workload_identity_enabled",
	}

	// automaticClusterEnabledFieldsInfo = fmt.Sprintf("For clusters with the Automatic SKU, the fields `%s` must be set to `true", strings.Join(automaticClusterEnabledFields, "`, `"))

	networkTemplate = fmt.Sprintf("```%s```", `
  network_profile {
    network_data_plane  = "cilium"
    network_plugin      = "azure"
    network_plugin_mode = "overlay"
    network_policy      = "cilium"
    outbound_type       = "managedNATGateway"
  }
`)

	automaticClusterNetworkProfileInfo = fmt.Sprintf("For clusters with the Automatic SKU, the following `network_profile` block required:\n\n%s",
		networkTemplate)
)

// customizeDiffAutomaticCluster validates that fields required by the Automatic cluster SKU
// are set to their expected fixed values.
func customizeDiffAutomaticCluster(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Get("sku_name").(string) != string(managedclusters.ManagedClusterSKUNameAutomatic) {
		return nil
	}

	var err error

	rawConfig := d.GetRawConfig().AsValueMap()
	// top-level boolean fields
	for _, field := range automaticClusterEnabledFields {
		if v := rawConfig[field]; v.IsKnown() && (v.IsNull() || v.False()) {
			err = multierror.Append(err, fmt.Errorf("`%s` must be `true` for Automatic SKU clusters", field))
		}
	}

	// sku_tier
	for field, expectedValue := range map[string]string{
		"sku_tier":                  "Standard",
		"automatic_upgrade_channel": "stable",
	} {
		if v := rawConfig[field]; v.IsKnown() && (v.IsNull() || v.AsString() != expectedValue) {
			err = multierror.Append(err, fmt.Errorf("`%s` must be `%s` for Automatic SKU clusters", field, expectedValue))
		}
	}

	// image_cleaner_interval_hours
	if v := d.Get("image_cleaner_interval_hours").(int); v != 168 {
		err = multierror.Append(err, fmt.Errorf("`image_cleaner_interval_hours` must be `168` for Automatic SKU clusters"))
	}

	// network_profile
	networkConfigMap := map[string]string{
		"network_data_plane":  "cilium",
		"network_plugin":      "azure",
		"network_plugin_mode": "overlay",
		"network_policy":      "cilium",
		"outbound_type":       "managedNATGateway",
	}
	networkList := rawConfig["network_profile"].AsValueSlice()
	if len(networkList) == 0 {
		err = multierror.Append(err, fmt.Errorf("`network_profile` block is required for Automatic SKU clusters"))
	}
	network := networkList[0].AsValueMap()
	for key, expectedValue := range networkConfigMap {
		if v := network[key]; v.IsKnown() && (v.IsNull() || v.AsString() != expectedValue) {
			err = multierror.Append(err, fmt.Errorf("`network_profile.0.%s` must be `%s` for Automatic SKU clusters", key, expectedValue))
		}
	}

	// default_node_pool
	if v := d.Get("default_node_pool.0.os_disk_type").(string); v != "Ephemeral" {
		err = multierror.Append(err, fmt.Errorf("`default_node_pool.0.os_disk_type` must be `Ephemeral` for Automatic SKU clusters"))
	}

	// workload_autoscaler_profile
	for _, field := range []string{
		"workload_autoscaler_profile.0.keda_enabled",
		"workload_autoscaler_profile.0.vertical_pod_autoscaler_enabled",
		"api_server_access_profile.0.virtual_network_integration_enabled",
		"azure_active_directory_role_based_access_control.0.azure_rbac_enabled",
		"key_vault_secrets_provider.0.secret_rotation_enabled",
	} {
		if v := d.Get(field).(bool); !v {
			err = multierror.Append(err, fmt.Errorf("`%s` must be `true` for Automatic SKU clusters", field))
		}
	}

	return err
}
