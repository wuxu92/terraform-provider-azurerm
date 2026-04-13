// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customizeDiffAutomaticCluster validates that fields required by the Automatic cluster SKU
// are set to their expected fixed values.
func customizeDiffAutomaticCluster(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	rawConfig := d.GetRawConfig().AsValueMap()
	skuName := rawConfig["sku_name"].AsString()
	if skuName != string(managedclusters.ManagedClusterSKUNameAutomatic) {
		return nil
	}

	// top-level boolean fields
	if v := d.Get("azure_policy_enabled").(bool); !v {
		return fmt.Errorf("`azure_policy_enabled` must be `true` for Automatic SKU clusters")
	}
	if v := d.Get("local_account_disabled").(bool); !v {
		return fmt.Errorf("`local_account_disabled` must be `true` for Automatic SKU clusters")
	}
	if v := d.Get("oidc_issuer_enabled").(bool); !v {
		return fmt.Errorf("`oidc_issuer_enabled` must be `true` for Automatic SKU clusters")
	}
	if v := d.Get("workload_identity_enabled").(bool); !v {
		return fmt.Errorf("`workload_identity_enabled` must be `true` for Automatic SKU clusters")
	}
	if v := d.Get("image_cleaner_enabled").(bool); !v {
		return fmt.Errorf("`image_cleaner_enabled` must be `true` for Automatic SKU clusters")
	}

	// sku_tier
	if v := d.Get("sku_tier").(string); v != "Standard" {
		return fmt.Errorf("`sku_tier` must be `Standard` for Automatic SKU clusters")
	}

	// automatic_upgrade_channel
	if v := d.Get("automatic_upgrade_channel").(string); v != "stable" {
		return fmt.Errorf("`automatic_upgrade_channel` must be `stable` for Automatic SKU clusters")
	}

	// image_cleaner_interval_hours
	if v := d.Get("image_cleaner_interval_hours").(int); v != 168 {
		return fmt.Errorf("`image_cleaner_interval_hours` must be `168` for Automatic SKU clusters")
	}

	// network_profile
	if v := d.Get("network_profile.0.network_data_plane").(string); v != "cilium" {
		return fmt.Errorf("`network_profile.0.network_data_plane` must be `cilium` for Automatic SKU clusters")
	}
	if v := d.Get("network_profile.0.network_plugin").(string); v != "azure" {
		return fmt.Errorf("`network_profile.0.network_plugin` must be `azure` for Automatic SKU clusters")
	}
	if v := d.Get("network_profile.0.network_plugin_mode").(string); v != "overlay" {
		return fmt.Errorf("`network_profile.0.network_plugin_mode` must be `overlay` for Automatic SKU clusters")
	}
	if v := d.Get("network_profile.0.network_policy").(string); v != "cilium" {
		return fmt.Errorf("`network_profile.0.network_policy` must be `cilium` for Automatic SKU clusters")
	}
	if v := d.Get("network_profile.0.outbound_type").(string); v != "managedNATGateway" {
		return fmt.Errorf("`network_profile.0.outbound_type` must be `managedNATGateway` for Automatic SKU clusters")
	}

	// default_node_pool
	if v := d.Get("default_node_pool.0.os_disk_type").(string); v != "Ephemeral" {
		return fmt.Errorf("`default_node_pool.0.os_disk_type` must be `Ephemeral` for Automatic SKU clusters")
	}

	// workload_autoscaler_profile
	if v := d.Get("workload_autoscaler_profile.0.keda_enabled").(bool); !v {
		return fmt.Errorf("`workload_autoscaler_profile.0.keda_enabled` must be `true` for Automatic SKU clusters")
	}
	if v := d.Get("workload_autoscaler_profile.0.vertical_pod_autoscaler_enabled").(bool); !v {
		return fmt.Errorf("`workload_autoscaler_profile.0.vertical_pod_autoscaler_enabled` must be `true` for Automatic SKU clusters")
	}

	// api_server_access_profile
	if v := d.Get("api_server_access_profile.0.virtual_network_integration_enabled").(bool); !v {
		return fmt.Errorf("`api_server_access_profile.0.virtual_network_integration_enabled` must be `true` for Automatic SKU clusters")
	}

	// azure_active_directory_role_based_access_control
	if v := d.Get("azure_active_directory_role_based_access_control.0.azure_rbac_enabled").(bool); !v {
		return fmt.Errorf("`azure_active_directory_role_based_access_control.0.azure_rbac_enabled` must be `true` for Automatic SKU clusters")
	}

	// key_vault_secrets_provider
	if v := d.Get("key_vault_secrets_provider.0.secret_rotation_enabled").(bool); !v {
		return fmt.Errorf("`key_vault_secrets_provider.0.secret_rotation_enabled` must be `true` for Automatic SKU clusters")
	}

	return nil
}
