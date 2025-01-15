// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurefleet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineProfileSchema(required bool) *pluginsdk.Schema {
	vmProfile := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: required,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"os_profile": osProfileSchema(),

				"source_image_reference": storageProfileSourceImageReferenceSchema(),

				"source_image_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.Any(
						images.ValidateImageID,
						validate.SharedImageID,
						validate.SharedImageVersionID,
						validate.CommunityGalleryImageID,
						validate.CommunityGalleryImageVersionID,
						validate.SharedGalleryImageID,
						validate.SharedGalleryImageVersionID,
					),
					ConflictsWith: []string{
						"virtual_machine_profile.0.source_image_reference",
					},
				},

				"os_disk": storageProfileOsDiskSchema(),

				"data_disk": storageProfileDataDiskSchema(),

				"boot_diagnostic_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"boot_diagnostic_storage_account_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"capacity_reservation_group_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,

					ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
				},

				"extension_operations_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"extension": extensionSchema(),

				"extensions_time_budget": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
				},

				"gallery_application": galleryApplicationSchema(),

				"license_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"RHEL_BYOS",
						"SLES_BYOS",
						"Windows_Client",
						"Windows_Server",
					}, false),
				},

				"network_interface": networkInterfaceSchema(),

				"scheduled_event_os_image_timeout": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "PT15M",
					ValidateFunc: validation.StringInSlice([]string{
						"P1T5M",
					}, false),
				},

				"scheduled_event_termination_timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
				},

				"security_profile": securityProfileSchema(),

				"service_artifact_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"user_data_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsBase64,
				},
			},
		},
	}
	if !required {
		vmProfile.Optional = true
	}
	return vmProfile
}

func galleryApplicationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 100,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"version_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: galleryapplicationversions.ValidateApplicationVersionID,
				},

				// Example: https://mystorageaccount.blob.core.windows.net/configurations/settings_json.config
				"configuration_blob_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntBetween(0, 2147483647),
				},

				// NOTE: Per the service team, "this is a pass through value that we just add to the model but don't depend on. It can be any string."
				"tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"treat_failure_as_deployment_failure_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func extensionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type_handler_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"auto_upgrade_minor_version_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"force_extension_execution_on_change": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"protected_settings_json": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				"protected_settings_from_key_vault": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,

					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"secret_url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: keyVaultValidate.NestedItemId,
							},

							"source_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
						},
					},
				},

				"extensions_to_provision_after_vm_creation": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"settings_json": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},
			},
		},
	}
}

func networkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,

					ValidateFunc: validation.StringIsNotEmpty,
				},

				"accelerated_networking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"auxiliary_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
						string(fleets.NetworkInterfaceAuxiliaryModeFloating),
					}, false),
				},

				"auxiliary_sku": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.NetworkInterfaceAuxiliarySkuATwo),
						string(fleets.NetworkInterfaceAuxiliarySkuAFour),
						string(fleets.NetworkInterfaceAuxiliarySkuAEight),
						string(fleets.NetworkInterfaceAuxiliarySkuAOne),
					}, false),
				},

				"ip_configuration": ipConfigurationSchema(),

				"ip_forwarding_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      fleets.DiskDeleteOptionTypesDelete,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"fpga_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"tcp_state_tracking_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}
}

func ipConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
					},
					Set:      pluginsdk.HashString,
					MaxItems: 20,
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"public_ip_address": publicIPAddressSchema(),

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(fleets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForIPVersion(), false),
				},
			},
		},
	}
}

func publicIPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      fleets.DeleteOptionsDelete,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDeleteOptions(), false),
				},

				"domain_name_label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"domain_name_label_scope": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDomainNameLabelScopeTypes(), false),
				},

				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(4, 32),
				},

				"ip_tag": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:     pluginsdk.TypeString,
								Required: true,

								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Type:     pluginsdk.TypeString,
								Required: true,

								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				// TODO: preview feature
				// $ az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress
				// $ az provider register -n Microsoft.Network
				"public_ip_prefix_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,

					ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
				},

				"sku": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:     pluginsdk.TypeString,
								Optional: true,

								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForPublicIPAddressSkuName(), false),
							},

							"tier": {
								Type:     pluginsdk.TypeString,
								Optional: true,

								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForPublicIPAddressSkuTier(), false),
							},
						},
					},
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,

					Default:      string(fleets.IPVersionIPvFour),
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForIPVersion(), false),
				},
			},
		},
	}
}

func osProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"custom_data_base64": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsBase64,
				},

				"linux_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					// ConflictsWith: []string{"virtual_machine_profile.0.os_profile.0.windows_configuration"},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"admin_username": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validateAdminUsernameLinux,
							},

							"computer_name_prefix": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.LinuxComputerNamePrefix,
							},

							"admin_password": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Sensitive:    true,
								ValidateFunc: validatePasswordComplexityLinux,
							},

							"admin_ssh_key": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"public_key": {
											Type:             pluginsdk.TypeString,
											Required:         true,
											ValidateFunc:     validate.SSHKey,
											DiffSuppressFunc: suppress.SSHKey,
										},
										"username": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},

							"password_authentication_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
							},

							"provision_vm_agent_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
							},

							"vm_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
							},

							// It is not supported in VMSS. need to confirm whether it needs to be exposed.
							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},

							"patch_assessment_mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(fleets.LinuxPatchAssessmentModeImageDefault),
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxPatchAssessmentMode(), false),
							},

							"patch_mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(fleets.LinuxVMGuestPatchModeImageDefault),
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchMode(), false),
							},

							// It is not supported in VMSS. need to confirm whether it needs to be exposed.
							"reboot_setting": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchAutomaticByPlatformRebootSetting(), false),
							},

							"secret": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

										"certificate": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											MinItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"url": {
														Type:         pluginsdk.TypeString,
														Required:     true,
														ValidateFunc: keyVaultValidate.NestedItemId,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},

				"windows_configuration": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					// ConflictsWith: []string{"virtual_machine_profile.0.os_profile.0.linux_configuration"},
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"admin_username": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validateAdminUsernameWindows,
							},

							"admin_password": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								Sensitive:    true,
								ValidateFunc: validatePasswordComplexityWindows,
							},

							"computer_name_prefix": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.WindowsComputerNamePrefix,
							},

							"additional_unattend_content": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"content": {
											Type:      pluginsdk.TypeString,
											Required:  true,
											Sensitive: true,
										},
										"setting": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSettingNames(), false),
										},

										"pass_name": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForPassName(), false),
										},

										"component_name": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForComponentName(), false),
										},
									},
								},
							},

							"automatic_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
							},

							"vm_agent_platform_updates_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
							},

							"patch_assessment_mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(fleets.WindowsPatchAssessmentModeImageDefault),
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxPatchAssessmentMode(), false),
							},

							"patch_mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(fleets.WindowsVMGuestPatchModeAutomaticByOS),
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchMode(), false),
							},

							"hot_patching_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},

							// It is not supported in VMSS. Need to confirm whether it needs to be exposed.
							"bypass_platform_safety_checks_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},

							// It is not supported in VMSS. Need to confirm whether it needs to be exposed.
							"reboot_setting": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForLinuxVMGuestPatchAutomaticByPlatformRebootSetting(), false),
							},

							"provision_vm_agent_enabled": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  true,
							},

							"secret": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

										"certificate": {
											Type:     pluginsdk.TypeSet,
											Required: true,
											MinItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"store": {
														Type:     pluginsdk.TypeString,
														Optional: true,
													},
													"url": {
														Type:         pluginsdk.TypeString,
														Required:     true,
														ValidateFunc: keyVaultValidate.NestedItemId,
													},
												},
											},
										},
									},
								},
							},

							"time_zone": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"winrm_listener": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"protocol": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForProtocolTypes(), false),
										},
										"certificate_url": {
											Type:         pluginsdk.TypeString,
											Optional:     true,
											ValidateFunc: keyVaultValidate.NestedItemId,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func securityProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"encryption_at_host_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"user_assigned_identity_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"proxy_agent": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{

							// remove this property as it is not supported!!!202501113
							// there is another property `enable` exists
							// need to confirm whether the following properties should be set when this feature is enable?
							// key_incarnation_value is required?
							"key_incarnation_value": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
							"mode": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      string(fleets.ModeEnforce),
								ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForMode(), false),
							},
						},
					},
				},

				"security_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityTypes(), false),
				},

				"uefi_secure_boot_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"uefi_vtpm_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func storageProfileDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
				},

				"create_option": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(fleets.DiskCreateOptionTypesEmpty),
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.DiskCreateOptionTypesEmpty),
						string(fleets.DiskCreateOptionTypesFromImage),
					}, false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      fleets.DiskDeleteOptionTypesDelete,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"disk_size_in_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"lun": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 2000),
				},

				"storage_account_type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForStorageAccountTypes(), false),
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func storageProfileOsDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					// whilst this appears in the Update block the API returns this when changing:
					// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
					ForceNew: true,
					// NOTE: OS Disks don't support Ultra SSDs or PremiumV2_LRS
					ValidateFunc: validation.StringInSlice([]string{
						string(fleets.StorageAccountTypesPremiumLRS),
						string(fleets.StorageAccountTypesPremiumZRS),
						string(fleets.StorageAccountTypesStandardLRS),
						string(fleets.StorageAccountTypesStandardSSDLRS),
						string(fleets.StorageAccountTypesStandardSSDZRS),
					}, false),
				},

				"caching": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						// NOTE: because there is a `None` value in the possible values, it's handled in the Create/Update and Read functions.
						string(fleets.CachingTypesReadOnly),
						string(fleets.CachingTypesReadWrite),
					}, false),
				},

				"disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_disk_encryption_set_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"security_encryption_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSecurityEncryptionTypes(), false),
				},

				"delete_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      fleets.DiskDeleteOptionTypesDelete,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiskDeleteOptionTypes(), false),
				},

				"diff_disk_option": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskOptions(), false),
				},
				"diff_disk_placement": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForDiffDiskPlacement(), false),
				},

				"disk_size_in_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"image_uri": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"vhd_containers": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func storageProfileSourceImageReferenceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		ConflictsWith: []string{
			"virtual_machine_profile.0.source_image_id",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"offer": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"sku": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceStorageProfileImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["publisher"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["offer"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["sku"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["version"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["community_gallery_image_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["shared_gallery_image_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func expandBaseVirtualMachineProfileModel(inputList []VirtualMachineProfileModel) (*fleets.BaseVirtualMachineProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}
	input := &inputList[0]
	output := fleets.BaseVirtualMachineProfile{
		OsProfile:                expandOSProfileModel(inputList),
		ApplicationProfile:       expandApplicationProfile(input.GalleryApplicationProfile),
		CapacityReservation:      expandCapacityReservation(input.CapacityReservationGroupId),
		DiagnosticsProfile:       expandDiagnosticsProfile(input.BootDiagnosticEnabled, input.BootDiagnosticStorageAccountEndpoint),
		NetworkProfile:           expandNetworkProfile(input.NetworkInterface),
		ScheduledEventsProfile:   expandScheduledEventsProfile(input),
		SecurityProfile:          expandSecurityProfileModel(input.SecurityProfile),
		ServiceArtifactReference: expandServiceArtifactReference(input.ServiceArtifactId),
	}

	extensionProfileValue, err := expandExtensionProfileModel(input.Extension, input.ExtensionsTimeBudget)
	if err != nil {
		return nil, err
	}
	output.ExtensionProfile = extensionProfileValue

	output.LicenseType = pointer.To("None")
	if input.LicenseType != "" {
		output.LicenseType = pointer.To(input.LicenseType)
	}

	if input.UserDataBase64 != "" {
		output.UserData = pointer.To(input.UserDataBase64)
	}

	storageProfile := &fleets.VirtualMachineScaleSetStorageProfile{
		ImageReference: expandImageReference(input.SourceImageReference, input.SourceImageId),
		OsDisk:         expandOSDiskModel(input.OsDisk),
	}

	dataDisks, err := expandDataDiskModel(input.DataDisks)
	if err != nil {
		return nil, err
	}
	storageProfile.DataDisks = dataDisks

	if len(input.OsProfile) > 0 {
		if len(input.OsProfile[0].LinuxConfiguration) > 0 {
			storageProfile.OsDisk.OsType = pointer.To(fleets.OperatingSystemTypesLinux)
		}
		if len(input.OsProfile[0].WindowsConfiguration) > 0 {
			storageProfile.OsDisk.OsType = pointer.To(fleets.OperatingSystemTypesWindows)
		}
	}
	output.StorageProfile = storageProfile

	return &output, nil
}

func expandServiceArtifactReference(id string) *fleets.ServiceArtifactReference {
	if id == "" {
		return nil
	}
	return &fleets.ServiceArtifactReference{Id: pointer.To(id)}
}

func expandApplicationProfile(inputList []GalleryApplicationModel) *fleets.ApplicationProfile {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VMGalleryApplication
	for _, v := range inputList {
		input := v
		output := fleets.VMGalleryApplication{
			EnableAutomaticUpgrade:          pointer.To(input.AutomaticUpgradeEnabled),
			Order:                           pointer.To(input.Order),
			PackageReferenceId:              input.VersionId,
			TreatFailureAsDeploymentFailure: pointer.To(input.TreatFailureAsDeploymentFailureEnabled),
		}

		if input.ConfigurationBlobUri != "" {
			output.ConfigurationReference = pointer.To(input.ConfigurationBlobUri)
		}

		if input.Tag != "" {
			output.Tags = pointer.To(input.Tag)
		}
		outputList = append(outputList, output)
	}

	return &fleets.ApplicationProfile{
		GalleryApplications: &outputList,
	}
}

func expandCapacityReservation(input string) *fleets.CapacityReservationProfile {
	if input == "" {
		return nil
	}

	return &fleets.CapacityReservationProfile{
		CapacityReservationGroup: expandSubResource(input),
	}
}

func expandSubResource(input string) *fleets.SubResource {
	if input == "" {
		return nil
	}

	return &fleets.SubResource{
		Id: pointer.To(input),
	}
}

func expandSubResources(inputList []string) *[]fleets.SubResource {
	if len(inputList) == 0 {
		return nil
	}
	var outputList []fleets.SubResource

	for _, v := range inputList {
		input := v

		output := expandSubResource(input)
		if output != nil {
			outputList = append(outputList, pointer.From(output))
		}
	}

	return &outputList
}

func expandDiagnosticsProfile(enabled bool, endpoint string) *fleets.DiagnosticsProfile {
	bootDiagnostics := fleets.BootDiagnostics{
		Enabled:    pointer.To(enabled),
		StorageUri: pointer.To(endpoint),
	}

	output := fleets.DiagnosticsProfile{
		BootDiagnostics: pointer.To(bootDiagnostics),
	}

	return &output
}

func expandExtensionProfileModel(inputList []ExtensionModel, timeBudget string) (*fleets.VirtualMachineScaleSetExtensionProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	output := fleets.VirtualMachineScaleSetExtensionProfile{}
	extensionsValue, err := expandExtensionsModel(inputList)
	if err != nil {
		return nil, err
	}
	output.Extensions = extensionsValue

	if timeBudget != "" {
		output.ExtensionsTimeBudget = pointer.To(timeBudget)
	}
	return &output, nil
}

func expandExtensionsModel(inputList []ExtensionModel) (*[]fleets.VirtualMachineScaleSetExtension, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	var outputList []fleets.VirtualMachineScaleSetExtension
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetExtension{}

		if input.Name != "" {
			output.Name = pointer.To(input.Name)
		}

		propertiesValue, err := expandExtensionModel(input)
		if err != nil {
			return nil, err
		}
		output.Properties = propertiesValue

		outputList = append(outputList, output)
	}
	return &outputList, nil
}

func expandExtensionModel(input ExtensionModel) (*fleets.VirtualMachineScaleSetExtensionProperties, error) {
	output := fleets.VirtualMachineScaleSetExtensionProperties{
		AutoUpgradeMinorVersion:       pointer.To(input.AutoUpgradeMinorVersionEnabled),
		EnableAutomaticUpgrade:        pointer.To(input.AutomaticUpgradeEnabled),
		ProtectedSettingsFromKeyVault: expandKeyVaultSecretReferenceModel(input.ProtectedSettingsFromKeyVault),
	}
	if input.ForceExtensionExecutionOnChange != "" {
		output.ForceUpdateTag = pointer.To(input.ForceExtensionExecutionOnChange)
	}

	if input.ProtectedSettingsJson != "" {
		protectedSettingsValue := make(map[string]interface{})
		err := json.Unmarshal([]byte(input.ProtectedSettingsJson), &protectedSettingsValue)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling `protected_settings_json`: %+v", err)
		}
		output.ProtectedSettings = pointer.To(protectedSettingsValue)
	}

	if len(input.ExtensionsToProvisionAfterVmCreation) > 0 {
		output.ProvisionAfterExtensions = pointer.To(input.ExtensionsToProvisionAfterVmCreation)
	}

	if input.Publisher != "" {
		output.Publisher = pointer.To(input.Publisher)
	}

	if input.SettingsJson != "" {
		result := make(map[string]interface{})
		err := json.Unmarshal([]byte(input.SettingsJson), &result)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling `settings_json`: %+v", err)
		}
		output.Settings = pointer.To(result)
	}

	if input.Type != "" {
		output.Type = pointer.To(input.Type)
	}

	if input.TypeHandlerVersion != "" {
		output.TypeHandlerVersion = pointer.To(input.TypeHandlerVersion)
	}

	return &output, nil
}

func expandKeyVaultSecretReferenceModel(inputList []ProtectedSettingsFromKeyVaultModel) *fleets.KeyVaultSecretReference {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.KeyVaultSecretReference{
		SecretURL:   input.SecretUrl,
		SourceVault: pointer.From(expandSubResource(input.SourceVaultId)),
	}

	return &output
}

func expandNetworkProfile(inputList []NetworkInterfaceModel) *fleets.VirtualMachineScaleSetNetworkProfile {
	if len(inputList) == 0 {
		return nil
	}

	output := fleets.VirtualMachineScaleSetNetworkProfile{
		// 2020-11-01 is the only valid value for this value and is only valid for VMSS in Orchestration Mode flex
		NetworkApiVersion:              pointer.To(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
		NetworkInterfaceConfigurations: expandNetworkInterfaceModel(inputList),
	}

	return &output
}

func expandApiEntityReferenceModel(input string) *fleets.ApiEntityReference {
	if input == "" {
		return nil
	}

	return &fleets.ApiEntityReference{
		Id: pointer.To(input),
	}
}

func expandNetworkInterfaceModel(inputList []NetworkInterfaceModel) *[]fleets.VirtualMachineScaleSetNetworkConfiguration {
	var outputList []fleets.VirtualMachineScaleSetNetworkConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetNetworkConfiguration{
			Name:       input.Name,
			Properties: expandNetworkConfigurationPropertiesModel(input),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandNetworkConfigurationPropertiesModel(input NetworkInterfaceModel) *fleets.VirtualMachineScaleSetNetworkConfigurationProperties {
	output := fleets.VirtualMachineScaleSetNetworkConfigurationProperties{
		DisableTcpStateTracking:     pointer.To(!input.TcpStateTrackingEnabled),
		EnableAcceleratedNetworking: pointer.To(input.AcceleratedNetworkingEnabled),
		EnableFpga:                  pointer.To(input.FpgaEnabled),
		EnableIPForwarding:          pointer.To(input.IPForwardingEnabled),
		NetworkSecurityGroup:        expandSubResource(input.NetworkSecurityGroupId),
		Primary:                     pointer.To(input.Primary),
	}

	if len(input.DnsServers) > 0 {
		output.DnsSettings = &fleets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
			DnsServers: pointer.To(input.DnsServers),
		}
	}

	auxiliaryMode := fleets.NetworkInterfaceAuxiliaryModeNone
	if input.AuxiliaryMode != "" {
		auxiliaryMode = fleets.NetworkInterfaceAuxiliaryMode(input.AuxiliaryMode)

	}
	output.AuxiliaryMode = pointer.To(auxiliaryMode)

	if input.DeleteOption != "" {
		output.DeleteOption = pointer.To(fleets.DeleteOptions(input.DeleteOption))
	}

	auxiliarySku := fleets.NetworkInterfaceAuxiliarySkuNone
	if input.AuxiliaryMode != "" {
		auxiliarySku = fleets.NetworkInterfaceAuxiliarySku(input.AuxiliarySku)

	}
	output.AuxiliarySku = pointer.To(auxiliarySku)

	output.IPConfigurations = pointer.From(expandIPConfigurationModel(input.IPConfiguration))

	return &output
}

func expandIPConfigurationModel(inputList []IPConfigurationModel) *[]fleets.VirtualMachineScaleSetIPConfiguration {
	var outputList []fleets.VirtualMachineScaleSetIPConfiguration
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPConfiguration{
			Name: input.Name,
			Properties: &fleets.VirtualMachineScaleSetIPConfigurationProperties{
				ApplicationGatewayBackendAddressPools: expandSubResources(input.ApplicationGatewayBackendAddressPoolIds),
				ApplicationSecurityGroups:             expandSubResources(input.ApplicationSecurityGroupIds),
				LoadBalancerBackendAddressPools:       expandSubResources(input.LoadBalancerBackendAddressPoolIds),
				Primary:                               pointer.To(input.Primary),
				PublicIPAddressConfiguration:          expandPublicIPAddressModel(input.PublicIPAddress),
				Subnet:                                expandApiEntityReferenceModel(input.SubnetId),
			},
		}

		if input.Version != "" {
			output.Properties.PrivateIPAddressVersion = pointer.To(fleets.IPVersion(input.Version))
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPublicIPAddressModel(inputList []PublicIPAddressModel) *fleets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: input.Name,
		Properties: &fleets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			DnsSettings:          expandPublicIPAddressDnsSettings(input.DomainNameLabel, input.DomainNameLabelScope),
			IPTags:               expandIPTagModel(input.IPTag),
			IdleTimeoutInMinutes: pointer.To(input.IdleTimeoutInMinutes),
			PublicIPPrefix:       expandSubResource(input.PublicIPPrefix),
		},
		Sku: expandPublicIPAddressSkuModel(input.Sku),
	}

	if input.DeleteOption != "" {
		output.Properties.DeleteOption = pointer.To(fleets.DeleteOptions(input.DeleteOption))
	}
	if input.Version != "" {
		output.Properties.PublicIPAddressVersion = pointer.To(fleets.IPVersion(input.Version))
	}

	return &output
}

func expandPublicIPAddressDnsSettings(domainNameLabel string, domainNameLabelScope string) *fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings {
	if domainNameLabel == "" && domainNameLabelScope == "" {
		return nil
	}
	output := fleets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{}
	if domainNameLabel != "" {
		output.DomainNameLabel = domainNameLabel
	}

	if domainNameLabelScope != "" {
		output.DomainNameLabelScope = pointer.To(fleets.DomainNameLabelScopeTypes(domainNameLabelScope))
	}

	return &output
}

func expandIPTagModel(inputList []IPTagModel) *[]fleets.VirtualMachineScaleSetIPTag {
	var outputList []fleets.VirtualMachineScaleSetIPTag
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetIPTag{}

		if input.Type != "" {
			output.IPTagType = pointer.To(input.Type)
		}

		if input.Tag != "" {
			output.Tag = pointer.To(input.Tag)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandPublicIPAddressSkuModel(inputList []SkuModel) *fleets.PublicIPAddressSku {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.PublicIPAddressSku{
		Name: pointer.To(fleets.PublicIPAddressSkuName(input.Name)),
		Tier: pointer.To(fleets.PublicIPAddressSkuTier(input.Tier)),
	}

	return &output
}

func expandOSProfileModel(inputList []VirtualMachineProfileModel) *fleets.VirtualMachineScaleSetOSProfile {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}
	osProfile := &inputList[0].OsProfile[0]
	output := fleets.VirtualMachineScaleSetOSProfile{
		AllowExtensionOperations: pointer.To(inputList[0].ExtensionOperationsEnabled),
	}
	if osProfile.CustomDataBase64 != "" {
		output.CustomData = pointer.To(osProfile.CustomDataBase64)
	}

	if lConfig := osProfile.LinuxConfiguration; len(lConfig) > 0 {
		linuxConfig := fleets.LinuxConfiguration{
			ProvisionVMAgent: pointer.To(lConfig[0].ProvisionVMAgentEnabled),
			Ssh:              expandSshConfigurationModel(lConfig[0].AdminSshKeys),
			PatchSettings:    &fleets.LinuxPatchSettings{},
		}

		// AutomaticByPlatformSettings cannot be set if the PatchMode is not 'AutomaticByPlatform'
		if lConfig[0].PatchMode == string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
			linuxConfig.PatchSettings.AutomaticByPlatformSettings = &fleets.LinuxVMGuestPatchAutomaticByPlatformSettings{
				BypassPlatformSafetyChecksOnUserSchedule: pointer.To(lConfig[0].BypassPlatformSafetyChecksEnabled),
			}

			if lConfig[0].RebootSetting != "" {
				linuxConfig.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.To(fleets.LinuxVMGuestPatchAutomaticByPlatformRebootSetting(lConfig[0].RebootSetting))
			}
		}
		if lConfig[0].PatchAssessmentMode != "" {
			linuxConfig.PatchSettings.AssessmentMode = pointer.To(fleets.LinuxPatchAssessmentMode(lConfig[0].PatchAssessmentMode))
		}
		if lConfig[0].PatchMode != "" {
			linuxConfig.PatchSettings.PatchMode = pointer.To(fleets.LinuxVMGuestPatchMode(lConfig[0].PatchMode))
		}
		if lConfig[0].AdminUsername != "" {
			output.AdminUsername = pointer.To(lConfig[0].AdminUsername)
		}
		if lConfig[0].AdminPassword != "" {
			output.AdminPassword = pointer.To(lConfig[0].AdminPassword)
		}
		if lConfig[0].ComputerNamePrefix != "" {
			output.ComputerNamePrefix = pointer.To(lConfig[0].ComputerNamePrefix)
		}
		output.Secrets = expandOsProfileSecretsModel(lConfig[0].Secret)

		linuxConfig.Ssh = expandSshConfigurationModel(lConfig[0].AdminSshKeys)

		output.LinuxConfiguration = &linuxConfig
	}

	if winConfig := osProfile.WindowsConfiguration; len(winConfig) > 0 {
		windowsConfig := fleets.WindowsConfiguration{
			AdditionalUnattendContent:    expandAdditionalUnattendContentModel(winConfig[0].AdditionalUnattendContent),
			EnableAutomaticUpdates:       pointer.To(winConfig[0].AutomaticUpdatesEnabled),
			EnableVMAgentPlatformUpdates: pointer.To(winConfig[0].VMAgentPlatformUpdatesEnabled),
			ProvisionVMAgent:             pointer.To(winConfig[0].ProvisionVMAgentEnabled),
			WinRM:                        expandWinRM(winConfig[0].WinRM),
			PatchSettings:                &fleets.PatchSettings{},
		}
		if winConfig[0].AdminUsername != "" {
			output.AdminUsername = pointer.To(winConfig[0].AdminUsername)
		}
		if winConfig[0].AdminPassword != "" {
			output.AdminPassword = pointer.To(winConfig[0].AdminPassword)
		}
		if winConfig[0].ComputerNamePrefix != "" {
			output.ComputerNamePrefix = pointer.To(winConfig[0].ComputerNamePrefix)
		}

		// AutomaticByPlatformSettings cannot be set if the PatchMode is not 'AutomaticByPlatform'
		if winConfig[0].PatchMode == string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
			windowsConfig.PatchSettings.AutomaticByPlatformSettings = &fleets.WindowsVMGuestPatchAutomaticByPlatformSettings{
				BypassPlatformSafetyChecksOnUserSchedule: pointer.To(winConfig[0].BypassPlatformSafetyChecksEnabled),
			}

			if winConfig[0].RebootSetting != "" {
				windowsConfig.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.To(fleets.WindowsVMGuestPatchAutomaticByPlatformRebootSetting(winConfig[0].RebootSetting))
			}
		}
		if winConfig[0].PatchAssessmentMode != "" {
			windowsConfig.PatchSettings.AssessmentMode = pointer.To(fleets.WindowsPatchAssessmentMode(winConfig[0].PatchAssessmentMode))
		}
		if winConfig[0].PatchMode != "" {
			windowsConfig.PatchSettings.PatchMode = pointer.To(fleets.WindowsVMGuestPatchMode(winConfig[0].PatchMode))
		}
		if winConfig[0].TimeZone != "" {
			windowsConfig.TimeZone = pointer.To(winConfig[0].TimeZone)
		}
		output.WindowsConfiguration = &windowsConfig

		output.Secrets = expandOsProfileSecretsModel(winConfig[0].Secret)
	}

	return &output
}

func validateWindowsSetting(inputList []VirtualMachineProfileModel, d *schema.ResourceDiff) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if len(input.OsProfile[0].WindowsConfiguration) > 0 {
		patchMode := input.OsProfile[0].WindowsConfiguration[0].PatchMode
		patchAssessmentMode := input.OsProfile[0].WindowsConfiguration[0].PatchAssessmentMode
		hotPatchingEnabled := input.OsProfile[0].WindowsConfiguration[0].HotPatchingEnabled
		provisionVMAgentEnabled := input.OsProfile[0].WindowsConfiguration[0].ProvisionVMAgentEnabled

		rebootSetting := input.OsProfile[0].WindowsConfiguration[0].RebootSetting
		bypassPlatformSafetyChecksEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["windows_configuration"].AsValueSlice()[0].AsValueMap()["bypass_platform_safety_checks_enabled"]
		if !bypassPlatformSafetyChecksEnabledExist.IsNull() || rebootSetting != "" {
			if patchMode != string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("`bypass_platform_safety_checks_enabled` and `reboot_setting` cannot be set if the `PatchMode` is not 'AutomaticByPlatform'")
			}
		}

		if input.ExtensionOperationsEnabled && !provisionVMAgentEnabled {
			return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent_enabled` is set to `false`")
		}

		if patchAssessmentMode == string(fleets.WindowsPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgentEnabled {
			return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent_enabled' must always be set to 'true'", fleets.WindowsPatchAssessmentModeAutomaticByPlatform)
		}

		isHotPatchEnabledImage := isValidHotPatchSourceImageReference(input.SourceImageReference)
		hasHealthExtension := false
		if v := input.Extension; len(v) > 0 && (v[0].Type == "ApplicationHealthLinux" || v[0].Type == "ApplicationHealthWindows") {
			hasHealthExtension = true
		}

		if isHotPatchEnabledImage {
			// it is a hot patching enabled image, validate hot patching enabled settings
			if patchMode != string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("when referencing a hot patching enabled image the 'patch_mode' field must always be set to %q", fleets.WindowsVMGuestPatchModeAutomaticByPlatform)
			}
			if !provisionVMAgentEnabled {
				return fmt.Errorf("when referencing a hot patching enabled image the 'provision_vm_agent_enabled' field must always be set to 'true'")
			}
			if !hasHealthExtension {
				return fmt.Errorf("when referencing a hot patching enabled image the 'extension' field must always contain a 'application health extension'")
			}
			if !hotPatchingEnabled {
				return fmt.Errorf("when referencing a hot patching enabled image the 'hotpatching_enabled' field must always be set to 'true'")
			}
		} else {
			// not a hot patching enabled image verify Automatic VM Guest Patching settings
			if patchMode == string(fleets.WindowsVMGuestPatchModeAutomaticByPlatform) {
				if !provisionVMAgentEnabled {
					return fmt.Errorf("when 'patch_mode' is set to %q then 'provision_vm_agent_enabled' must be set to 'true'", patchMode)
				}
				if !hasHealthExtension {
					return fmt.Errorf("when 'patch_mode' is set to %q then the 'extension' field must always contain a 'application health extension'", patchMode)
				}
			}

			if hotPatchingEnabled {
				return fmt.Errorf("'hot_patching_enabled' field is not supported unless you are using one of the following hot patching enable images, '2022-datacenter-azure-edition', '2022-datacenter-azure-edition-core-smalldisk', '2022-datacenter-azure-edition-hotpatch' or '2022-datacenter-azure-edition-hotpatch-smalldisk'")
			}
		}
	}
	return nil
}

func validateLinuxSetting(inputList []VirtualMachineProfileModel, d *schema.ResourceDiff) error {
	if len(inputList) == 0 || len(inputList[0].OsProfile) == 0 {
		return nil
	}

	input := &inputList[0]
	if len(input.OsProfile[0].LinuxConfiguration) > 0 {
		patchMode := input.OsProfile[0].LinuxConfiguration[0].PatchMode
		patchAssessmentMode := input.OsProfile[0].LinuxConfiguration[0].PatchAssessmentMode
		provisionVMAgentEnabled := input.OsProfile[0].LinuxConfiguration[0].ProvisionVMAgentEnabled

		rebootSetting := input.OsProfile[0].LinuxConfiguration[0].RebootSetting
		bypassPlatformSafetyChecksEnabledExist := d.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["os_profile"].AsValueSlice()[0].AsValueMap()["linux_configuration"].AsValueSlice()[0].AsValueMap()["bypass_platform_safety_checks_enabled"]
		if !bypassPlatformSafetyChecksEnabledExist.IsNull() || rebootSetting != "" {
			if patchMode != string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("`bypass_platform_safety_checks_enabled` and `reboot_setting` cannot be set if the `PatchMode` is not 'AutomaticByPlatform'")
			}
		}

		if input.ExtensionOperationsEnabled && !provisionVMAgentEnabled {
			return fmt.Errorf("`extension_operations_enabled` cannot be set to `true` when `provision_vm_agent_enabled` is set to `false`")
		}

		if patchAssessmentMode == string(fleets.WindowsPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgentEnabled {
			return fmt.Errorf("when the 'patch_assessment_mode' field is set to %q the 'provision_vm_agent_enabled' must always be set to 'true'", fleets.LinuxPatchAssessmentModeAutomaticByPlatform)
		}

		hasHealthExtension := false
		if v := input.Extension; len(v) > 0 && (v[0].Type == "ApplicationHealthLinux" || v[0].Type == "ApplicationHealthWindows") {
			hasHealthExtension = true
		}

		if patchMode == string(fleets.LinuxVMGuestPatchModeAutomaticByPlatform) {
			if !provisionVMAgentEnabled {
				return fmt.Errorf("when the 'patch_mode' field is set to %q the 'provision_vm_agent_enabled' field must always be set to 'true', got %q", patchMode, strconv.FormatBool(provisionVMAgentEnabled))
			}

			if !hasHealthExtension {
				return fmt.Errorf("when the 'patch_mode' field is set to %q the 'extension' field must contain at least one 'application health extension', got 0", patchMode)
			}
		}
	}
	return nil
}

func isValidHotPatchSourceImageReference(referenceInput []SourceImageReferenceModel) bool {
	if len(referenceInput) == 0 {
		return false
	}
	raw := referenceInput[0]
	pub := raw.Publisher
	offer := raw.Offer
	sku := raw.Sku

	if pub == "MicrosoftWindowsServer" && offer == "WindowsServer" && (sku == "2022-datacenter-azure-edition-core" || sku == "2022-datacenter-azure-edition-core-smalldisk" || sku == "2022-datacenter-azure-edition-hotpatch" || sku == "2022-datacenter-azure-edition-hotpatch-smalldisk") {
		return true
	}

	return false
}

func expandSshConfigurationModel(inputList []AdminSshKeyModel) *fleets.SshConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	var publicKeys []fleets.SshPublicKey
	for _, v := range inputList {
		input := v
		output := fleets.SshPublicKey{
			Path: pointer.To(fmt.Sprintf("/home/%s/.ssh/authorized_keys", input.Username)),
		}
		if input.PublicKey != "" {
			output.KeyData = pointer.To(input.PublicKey)
		}
		publicKeys = append(publicKeys, output)
	}

	return &fleets.SshConfiguration{
		PublicKeys: pointer.To(publicKeys),
	}
}

func expandOsProfileSecretsModel(inputList []SecretModel) *[]fleets.VaultSecretGroup {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.VaultSecretGroup
	for _, v := range inputList {
		input := v
		output := fleets.VaultSecretGroup{
			SourceVault:       expandSubResource(input.KeyVaultId),
			VaultCertificates: expandVaultCertificateModel(input.Certificates),
		}

		outputList = append(outputList, output)
	}
	return &outputList
}

func expandVaultCertificateModel(inputList []CertificateModel) *[]fleets.VaultCertificate {
	var outputList []fleets.VaultCertificate
	for _, v := range inputList {
		input := v
		output := fleets.VaultCertificate{}

		if input.Store != "" {
			output.CertificateStore = pointer.To(input.Store)
		}

		if input.Url != "" {
			output.CertificateURL = pointer.To(input.Url)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandAdditionalUnattendContentModel(inputList []AdditionalUnattendContentModel) *[]fleets.AdditionalUnattendContent {
	if len(inputList) == 0 {
		return nil
	}

	var outputList []fleets.AdditionalUnattendContent
	for _, v := range inputList {
		input := v
		output := fleets.AdditionalUnattendContent{
			ComponentName: pointer.To(fleets.ComponentName(input.ComponentName)),
			PassName:      pointer.To(fleets.PassName(input.PassName)),
			SettingName:   pointer.To(fleets.SettingNames(input.SettingName)),
			Content:       pointer.To(input.Content),
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandWinRM(inputList []WinRMModel) *fleets.WinRMConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	var listenerList []fleets.WinRMListener
	for _, v := range inputList {
		input := v
		output := fleets.WinRMListener{
			Protocol: pointer.To(fleets.ProtocolTypes(input.Protocol)),
		}

		if input.CertificateUrl != "" {
			output.CertificateURL = pointer.To(input.CertificateUrl)
		}
		listenerList = append(listenerList, output)
	}

	return &fleets.WinRMConfiguration{
		Listeners: pointer.To(listenerList),
	}
}

func expandScheduledEventsProfile(input *VirtualMachineProfileModel) *fleets.ScheduledEventsProfile {
	if input == nil || (input.ScheduledEventTerminationTimeout == "" && input.ScheduledEventOsImageTimeout == "") {
		return nil
	}

	outPut := &fleets.ScheduledEventsProfile{}
	if input.ScheduledEventTerminationTimeout != "" {
		outPut.TerminateNotificationProfile = &fleets.TerminateNotificationProfile{
			Enable:           pointer.To(true),
			NotBeforeTimeout: pointer.To(input.ScheduledEventTerminationTimeout),
		}
	}

	if input.ScheduledEventOsImageTimeout != "" {
		outPut.OsImageNotificationProfile = &fleets.OSImageNotificationProfile{
			Enable:           pointer.To(true),
			NotBeforeTimeout: pointer.To(input.ScheduledEventOsImageTimeout),
		}
	}

	return outPut
}

func expandSecurityProfileModel(inputList []SecurityProfileModel) *fleets.SecurityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.SecurityProfile{
		EncryptionAtHost:   pointer.To(input.EncryptionAtHostEnabled),
		EncryptionIdentity: expandEncryptionIdentityModel(input.UserAssignedIdentityId),
		ProxyAgentSettings: expandProxyAgentModel(input.ProxyAgent),
		SecurityType:       pointer.To(fleets.SecurityTypes(input.SecurityType)),
		UefiSettings: &fleets.UefiSettings{
			SecureBootEnabled: pointer.To(input.UefiSecureBootEnabled),
			VTpmEnabled:       pointer.To(input.UefiVTpmEnabled),
		},
	}

	return &output
}

func expandEncryptionIdentityModel(id string) *fleets.EncryptionIdentity {
	if id == "" {
		return nil
	}

	return &fleets.EncryptionIdentity{
		UserAssignedIdentityResourceId: pointer.To(id),
	}
}

func expandProxyAgentModel(inputList []ProxyAgentModel) *fleets.ProxyAgentSettings {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.ProxyAgentSettings{
		Enabled:          pointer.To(true),
		KeyIncarnationId: pointer.To(input.KeyIncarnationValue),
		Mode:             pointer.To(fleets.Mode(input.mode)),
	}

	return &output
}

func expandDataDiskModel(inputList []DataDiskModel) (*[]fleets.VirtualMachineScaleSetDataDisk, error) {
	var outputList []fleets.VirtualMachineScaleSetDataDisk
	for _, v := range inputList {
		input := v
		output := fleets.VirtualMachineScaleSetDataDisk{
			CreateOption:            fleets.DiskCreateOptionTypes(input.CreateOption),
			DeleteOption:            pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption)),
			Lun:                     input.Lun,
			WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
		}

		if input.DiskSizeInGB > 0 {
			output.DiskSizeGB = pointer.To(input.DiskSizeInGB)
		}

		caching := string(fleets.CachingTypesNone)
		if input.Caching != "" {
			caching = input.Caching
		}
		output.Caching = pointer.To(fleets.CachingTypes(caching))

		if input.Name != "" {
			output.Name = &input.Name
		}

		managedDisk := &fleets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: pointer.To(fleets.StorageAccountTypes(input.StorageAccountType)),
		}
		if input.DiskEncryptionSetId != "" {
			managedDisk.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
				Id: pointer.To(input.DiskEncryptionSetId),
			}
		}
		output.ManagedDisk = managedDisk

		outputList = append(outputList, output)
	}
	return &outputList, nil
}

func expandVMDiskSecurityProfileModel(securityEncryptionType string, securityDiskEncryptionSetId string) *fleets.VMDiskSecurityProfile {
	if securityEncryptionType == "" && securityDiskEncryptionSetId == "" {
		return nil
	}

	output := fleets.VMDiskSecurityProfile{}
	if securityDiskEncryptionSetId != "" {
		output.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
			Id: pointer.To(securityDiskEncryptionSetId),
		}
	}

	if securityEncryptionType != "" {
		output.SecurityEncryptionType = pointer.To(fleets.SecurityEncryptionTypes(securityEncryptionType))
	}

	return &output
}

func expandImageReference(inputList []SourceImageReferenceModel, imageId string) *fleets.ImageReference {
	if imageId != "" {
		// With Version            : "/communityGalleries/publicGalleryName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/communityGalleries/publicGalleryName/images/myGalleryImageName"
		if _, errors := validation.Any(validate.CommunityGalleryImageID, validate.CommunityGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &fleets.ImageReference{
				CommunityGalleryImageId: pointer.To(imageId),
			}
		}

		// With Version            : "/sharedGalleries/galleryUniqueName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/sharedGalleries/galleryUniqueName/images/myGalleryImageName"
		if _, errors := validation.Any(validate.SharedGalleryImageID, validate.SharedGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &fleets.ImageReference{
				SharedGalleryImageId: pointer.To(imageId),
			}
		}

		return &fleets.ImageReference{
			Id: pointer.To(imageId),
		}
	}

	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	return &fleets.ImageReference{
		Publisher: pointer.To(input.Publisher),
		Offer:     pointer.To(input.Offer),
		Sku:       pointer.To(input.Sku),
		Version:   pointer.To(input.Version),
	}
}

func expandOSDiskModel(inputList []OSDiskModel) *fleets.VirtualMachineScaleSetOSDisk {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]

	output := fleets.VirtualMachineScaleSetOSDisk{
		// these have to be hard-coded so there's no point exposing them
		CreateOption:            fleets.DiskCreateOptionTypesFromImage,
		DiffDiskSettings:        expandDiffDiskSettingsModel(input),
		WriteAcceleratorEnabled: pointer.To(input.WriteAcceleratorEnabled),
	}
	managedDisk := &fleets.VirtualMachineScaleSetManagedDiskParameters{
		StorageAccountType: pointer.To(fleets.StorageAccountTypes(input.StorageAccountType)),
	}
	if input.DiskEncryptionSetId != "" {
		managedDisk.DiskEncryptionSet = &fleets.DiskEncryptionSetParameters{
			Id: pointer.To(input.DiskEncryptionSetId),
		}
	}
	if input.SecurityEncryptionType != "" && input.SecurityDiskEncryptionSetId != "" {
		managedDisk.SecurityProfile = expandVMDiskSecurityProfileModel(input.SecurityEncryptionType, input.SecurityDiskEncryptionSetId)
	}
	output.ManagedDisk = managedDisk

	if input.DiskSizeInGB > 0 {
		output.DiskSizeGB = pointer.To(input.DiskSizeInGB)
	}

	caching := string(fleets.CachingTypesNone)
	if v := input.Caching; v != "" {
		caching = v
	}
	output.Caching = pointer.To(fleets.CachingTypes(caching))

	if input.DeleteOption != "" {
		output.DeleteOption = pointer.To(fleets.DiskDeleteOptionTypes(input.DeleteOption))
	}
	if input.ImageUri != "" {
		output.Image = expandImage(input.ImageUri)
	}
	if input.Name != "" {
		output.Name = pointer.To(input.Name)
	}
	if len(input.VhdContainers) > 0 {
		output.VhdContainers = pointer.To(input.VhdContainers)
	}

	return &output
}

func expandDiffDiskSettingsModel(input *OSDiskModel) *fleets.DiffDiskSettings {
	if input == nil || (input.DiffDiskOption == "" && input.DiffDiskPlacement == "") {
		return nil
	}

	output := fleets.DiffDiskSettings{}
	if input.DiffDiskOption != "" {
		output.Option = pointer.To(fleets.DiffDiskOptions(input.DiffDiskOption))
	}
	if input.DiffDiskPlacement != "" {
		output.Placement = pointer.To(fleets.DiffDiskPlacement(input.DiffDiskPlacement))
	}

	return &output
}

func expandImage(input string) *fleets.VirtualHardDisk {
	if input == "" {
		return nil
	}

	return &fleets.VirtualHardDisk{Uri: pointer.To(input)}
}

func flattenVirtualMachineProfileModel(input *fleets.BaseVirtualMachineProfile, metadata sdk.ResourceMetaData) ([]VirtualMachineProfileModel, error) {
	var outputList []VirtualMachineProfileModel
	if input == nil {
		return outputList, nil
	}
	output := VirtualMachineProfileModel{
		GalleryApplicationProfile: flattenApplicationProfileModel(input.ApplicationProfile),
		NetworkInterface:          flattenNetworkInterfaceModel(input.NetworkProfile),
		SecurityProfile:           flattenSecurityProfileModel(input.SecurityProfile),
	}

	if v := input.OsProfile; v != nil {
		osProfile, err := flattenOSProfileModel(v, metadata.ResourceData)
		if err != nil {
			return outputList, err
		}
		output.OsProfile = osProfile
		output.ExtensionOperationsEnabled = pointer.From(v.AllowExtensionOperations)
	}

	if v := input.ServiceArtifactReference; v != nil {
		output.ServiceArtifactId = pointer.From(v.Id)
	}

	if v := input.StorageProfile; v != nil {
		output.DataDisks = flattenDataDiskModel(v.DataDisks)
		var storageImageId string
		if v.ImageReference != nil && v.ImageReference.Id != nil {
			storageImageId = *v.ImageReference.Id
		}
		if v.ImageReference != nil && v.ImageReference.CommunityGalleryImageId != nil {
			storageImageId = *v.ImageReference.CommunityGalleryImageId
		}
		if v.ImageReference != nil && v.ImageReference.SharedGalleryImageId != nil {
			storageImageId = *v.ImageReference.SharedGalleryImageId
		}
		output.SourceImageId = storageImageId
		output.SourceImageReference = flattenImageReference(v.ImageReference, storageImageId != "")
		output.OsDisk = flattenOSDiskModel(v.OsDisk)
	}

	if se := input.ScheduledEventsProfile; se != nil {
		if v := se.TerminateNotificationProfile; v != nil {
			output.ScheduledEventTerminationTimeout = pointer.From(v.NotBeforeTimeout)
		}
		if v := se.OsImageNotificationProfile; v != nil {
			output.ScheduledEventOsImageTimeout = pointer.From(v.NotBeforeTimeout)
		}
	}

	if cr := input.CapacityReservation; cr != nil {
		if v := cr.CapacityReservationGroup; v != nil {
			output.CapacityReservationGroupId = pointer.From(v.Id)
		}
	}

	if dp := input.DiagnosticsProfile; dp != nil {
		if v := dp.BootDiagnostics; v != nil {
			output.BootDiagnosticEnabled = pointer.From(v.Enabled)
			output.BootDiagnosticStorageAccountEndpoint = pointer.From(v.StorageUri)
		}
	}

	extensionProfileValue, err := flattenExtensionModel(input.ExtensionProfile, metadata)
	if err != nil {
		return nil, err
	}
	output.Extension = extensionProfileValue

	licenseType := ""
	if v := pointer.From(input.LicenseType); v != "None" {
		licenseType = v
	}
	output.LicenseType = licenseType

	if input.UserData != nil {
		output.UserDataBase64 = *input.UserData
	}

	return append(outputList, output), nil
}

func flattenAdminSshKeyModel(input *fleets.SshConfiguration) ([]AdminSshKeyModel, error) {
	var outputList []AdminSshKeyModel
	if input == nil || input.PublicKeys == nil {
		return outputList, nil
	}

	for _, input := range *input.PublicKeys {

		output := AdminSshKeyModel{}
		username := parseUsernameFromAuthorizedKeysPath(*input.Path)
		if username == nil {
			return nil, fmt.Errorf("parsing username from %q", *input.Path)
		}

		output.PublicKey = pointer.From(input.KeyData)
		output.Username = pointer.From(username)

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func parseUsernameFromAuthorizedKeysPath(input string) *string {
	// the Azure VM agent hard-codes this to `/home/username/.ssh/authorized_keys`
	// as such we can hard-code this for a better UX
	r := regexp.MustCompile("(/home/)+(?P<username>.*?)(/.ssh/authorized_keys)+")

	keys := r.SubexpNames()
	values := r.FindStringSubmatch(input)

	if values == nil {
		return nil
	}

	for i, k := range keys {
		if k == "username" {
			value := values[i]
			return &value
		}
	}

	return nil
}

func flattenApplicationProfileModel(input *fleets.ApplicationProfile) []GalleryApplicationModel {
	var outputList []GalleryApplicationModel
	if input == nil {
		return outputList
	}

	for _, input := range *input.GalleryApplications {
		output := GalleryApplicationModel{}
		output.VersionId = input.PackageReferenceId
		output.ConfigurationBlobUri = pointer.From(input.ConfigurationReference)
		output.AutomaticUpgradeEnabled = pointer.From(input.EnableAutomaticUpgrade)
		output.Order = pointer.From(input.Order)
		output.Tag = pointer.From(input.Tags)
		output.TreatFailureAsDeploymentFailureEnabled = pointer.From(input.TreatFailureAsDeploymentFailure)

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenNetworkInterfaceModel(input *fleets.VirtualMachineScaleSetNetworkProfile) []NetworkInterfaceModel {
	var outputList []NetworkInterfaceModel
	if input == nil {
		return outputList
	}

	for _, input := range *input.NetworkInterfaceConfigurations {
		output := NetworkInterfaceModel{
			Name: input.Name,
		}

		if props := input.Properties; props != nil {
			auxiliaryMode := ""
			if v := props.AuxiliaryMode; v != nil && *v != fleets.NetworkInterfaceAuxiliaryModeNone {
				auxiliaryMode = string(*v)
			}
			output.AuxiliaryMode = auxiliaryMode

			auxiliarySku := ""
			if v := props.AuxiliarySku; v != nil && *v != fleets.NetworkInterfaceAuxiliarySkuNone {
				auxiliarySku = string(*v)
			}
			output.AuxiliarySku = auxiliarySku

			output.DeleteOption = string(pointer.From(props.DeleteOption))

			output.TcpStateTrackingEnabled = !pointer.From(props.DisableTcpStateTracking)

			if v := props.DnsSettings; v != nil {
				output.DnsServers = pointer.From(v.DnsServers)
			}

			output.TcpStateTrackingEnabled = !pointer.From(props.DisableTcpStateTracking)

			output.AcceleratedNetworkingEnabled = pointer.From(props.EnableAcceleratedNetworking)

			output.FpgaEnabled = pointer.From(props.EnableFpga)

			output.IPForwardingEnabled = pointer.From(props.EnableIPForwarding)

			output.IPConfiguration = flattenIPConfigurationModel(props.IPConfigurations)

			if v := props.NetworkSecurityGroup; v != nil {
				output.NetworkSecurityGroupId = pointer.From(v.Id)
			}

			output.Primary = pointer.From(props.Primary)
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenOsProfileSecretsModel(inputList *[]fleets.VaultSecretGroup) []SecretModel {
	var outputList []SecretModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := SecretModel{
			Certificates: flattenVaultCertificateModel(input.VaultCertificates),
		}
		if v := input.SourceVault; v != nil {
			output.KeyVaultId = pointer.From(v.Id)
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenOSProfileModel(input *fleets.VirtualMachineScaleSetOSProfile, d *schema.ResourceData) ([]OSProfileModel, error) {
	var outputList []OSProfileModel
	if input == nil {
		return outputList, nil
	}

	output := OSProfileModel{}
	output.CustomDataBase64 = pointer.From(input.CustomData)

	var windowsConfigs []WindowsConfigurationModel
	if v := input.WindowsConfiguration; v != nil {
		windowsConfig := WindowsConfigurationModel{
			AdditionalUnattendContent:     flattenAdditionalUnattendContentModel(v.AdditionalUnattendContent),
			WinRM:                         flattenWinRMModel(v.WinRM),
			AdminUsername:                 pointer.From(input.AdminUsername),
			AdminPassword:                 d.Get("virtual_machine_profile.0.os_profile.0.windows_configuration.0.admin_password").(string),
			AutomaticUpdatesEnabled:       pointer.From(v.EnableAutomaticUpdates),
			ComputerNamePrefix:            pointer.From(input.ComputerNamePrefix),
			VMAgentPlatformUpdatesEnabled: pointer.From(v.ProvisionVMAgent),
			ProvisionVMAgentEnabled:       pointer.From(v.ProvisionVMAgent),
			TimeZone:                      pointer.From(v.TimeZone),
			Secret:                        flattenOsProfileSecretsModel(input.Secrets),
		}

		if p := v.PatchSettings; p != nil {
			windowsConfig.PatchMode = string(pointer.From(p.PatchMode))
			windowsConfig.PatchAssessmentMode = string(pointer.From(p.AssessmentMode))
			if a := p.AutomaticByPlatformSettings; a != nil {
				windowsConfig.BypassPlatformSafetyChecksEnabled = pointer.From(a.BypassPlatformSafetyChecksOnUserSchedule)
				windowsConfig.RebootSetting = string(pointer.From(a.RebootSetting))
			}
			windowsConfig.HotPatchingEnabled = pointer.From(p.EnableHotpatching)
		}
		windowsConfigs = append(windowsConfigs, windowsConfig)
	}
	output.WindowsConfiguration = windowsConfigs

	var linuxConfigs []LinuxConfigurationModel
	if v := input.LinuxConfiguration; v != nil {
		linuxConfig := LinuxConfigurationModel{
			AdminUsername:                 pointer.From(input.AdminUsername),
			AdminPassword:                 d.Get("virtual_machine_profile.0.os_profile.0.linux_configuration.0.admin_password").(string),
			ComputerNamePrefix:            pointer.From(input.ComputerNamePrefix),
			PasswordAuthenticationEnabled: !pointer.From(v.DisablePasswordAuthentication),
			VMAgentPlatformUpdatesEnabled: pointer.From(v.ProvisionVMAgent),
			ProvisionVMAgentEnabled:       pointer.From(v.ProvisionVMAgent),
			Secret:                        flattenOsProfileSecretsModel(input.Secrets),
		}

		if p := v.PatchSettings; p != nil {
			linuxConfig.PatchAssessmentMode = string(pointer.From(p.AssessmentMode))
			linuxConfig.PatchMode = string(pointer.From(p.PatchMode))
			if a := p.AutomaticByPlatformSettings; a != nil {
				linuxConfig.BypassPlatformSafetyChecksEnabled = pointer.From(a.BypassPlatformSafetyChecksOnUserSchedule)
				linuxConfig.RebootSetting = string(pointer.From(a.RebootSetting))
			}
		}

		flattenedSSHKeys, err := flattenAdminSshKeyModel(v.Ssh)
		if err != nil {
			return nil, fmt.Errorf("flattening `linux_configuration.0.admin_ssh_key`: %+v", err)
		}
		linuxConfig.AdminSshKeys = flattenedSSHKeys

		linuxConfigs = append(linuxConfigs, linuxConfig)
	}

	output.LinuxConfiguration = linuxConfigs

	return append(outputList, output), nil
}

func flattenVaultCertificateModel(inputList *[]fleets.VaultCertificate) []CertificateModel {
	var outputList []CertificateModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := CertificateModel{}

		output.Store = pointer.From(input.CertificateStore)
		output.Url = pointer.From(input.CertificateURL)

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenAdditionalUnattendContentModel(inputList *[]fleets.AdditionalUnattendContent) []AdditionalUnattendContentModel {
	var outputList []AdditionalUnattendContentModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := AdditionalUnattendContentModel{}
		output.ComponentName = string(pointer.From(input.ComponentName))
		output.Content = pointer.From(input.Content)
		output.PassName = string(pointer.From(input.PassName))
		output.SettingName = string(pointer.From(input.SettingName))
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenWinRMModel(input *fleets.WinRMConfiguration) []WinRMModel {
	var outputList []WinRMModel
	if input == nil || input.Listeners == nil {
		return outputList
	}

	for _, input := range *input.Listeners {
		output := WinRMModel{}
		output.CertificateUrl = pointer.From(input.CertificateURL)
		output.Protocol = string(pointer.From(input.Protocol))
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenSecurityProfileModel(input *fleets.SecurityProfile) []SecurityProfileModel {
	var outputList []SecurityProfileModel
	if input == nil {
		return outputList
	}
	output := SecurityProfileModel{
		ProxyAgent: flattenProxyAgentModel(input.ProxyAgentSettings),
	}

	output.EncryptionAtHostEnabled = pointer.From(input.EncryptionAtHost)

	if v := input.EncryptionIdentity; v != nil {
		output.UserAssignedIdentityId = pointer.From(v.UserAssignedIdentityResourceId)
	}

	if v := input.UefiSettings; v != nil {
		output.UefiSecureBootEnabled = pointer.From(v.SecureBootEnabled)
		output.UefiVTpmEnabled = pointer.From(v.VTpmEnabled)
	}

	output.SecurityType = string(pointer.From(input.SecurityType))

	return append(outputList, output)
}

func flattenProxyAgentModel(input *fleets.ProxyAgentSettings) []ProxyAgentModel {
	var outputList []ProxyAgentModel
	if input == nil {
		return outputList
	}

	output := ProxyAgentModel{}
	output.KeyIncarnationValue = pointer.From(input.KeyIncarnationId)
	output.mode = string(pointer.From(input.Mode))

	return append(outputList, output)
}

func flattenDataDiskModel(inputList *[]fleets.VirtualMachineScaleSetDataDisk) []DataDiskModel {
	var outputList []DataDiskModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := DataDiskModel{
			CreateOption: string(input.CreateOption),
			Lun:          input.Lun,
		}

		caching := ""
		if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
			caching = string(*v)
		}
		output.Caching = caching

		output.DeleteOption = string(pointer.From(input.DeleteOption))
		// Property 'dataDisk.diskMBpsReadWrite' can be enabled only on VMs in a Virtual Machine Scale Set?
		//output.DiskIOPSReadWrite = pointer.From(input.DiskIOPSReadWrite)
		//output.DiskMbpsReadWrite = pointer.From(input.DiskMBpsReadWrite)
		output.DiskSizeInGB = pointer.From(input.DiskSizeGB)
		output.Name = pointer.From(input.Name)
		output.WriteAcceleratorEnabled = pointer.From(input.WriteAcceleratorEnabled)

		if v := input.ManagedDisk; v != nil {
			if v := v.DiskEncryptionSet; v != nil {
				output.DiskEncryptionSetId = pointer.From(v.Id)
			}
			//if sp := v.SecurityProfile; sp != nil {
			//	if v := sp.DiskEncryptionSet; v != nil {
			//		output.SecurityDiskEncryptionSetId = pointer.From(v.Id)
			//	}
			//	output.SecurityEncryptionType = string(pointer.From(sp.SecurityEncryptionType))
			//}
			output.StorageAccountType = string(pointer.From(v.StorageAccountType))
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenImageReference(input *fleets.ImageReference, hasImageId bool) []SourceImageReferenceModel {
	var outputList []SourceImageReferenceModel
	if input == nil || hasImageId {
		return outputList
	}
	output := SourceImageReferenceModel{}
	output.Version = pointer.From(input.ExactVersion)
	output.Offer = pointer.From(input.Offer)
	output.Publisher = pointer.From(input.Publisher)
	output.Sku = pointer.From(input.Sku)
	output.Version = pointer.From(input.Version)

	return append(outputList, output)
}

func flattenOSDiskModel(input *fleets.VirtualMachineScaleSetOSDisk) []OSDiskModel {
	var outputList []OSDiskModel
	if input == nil {
		return outputList
	}

	output := OSDiskModel{}

	if v := input.DiffDiskSettings; v != nil {
		output.DiffDiskOption = string(pointer.From(v.Option))
		output.DiffDiskPlacement = string(pointer.From(v.Placement))
	}

	if v := input.Image; v != nil {
		output.ImageUri = pointer.From(v.Uri)
	}

	caching := ""
	if v := input.Caching; v != nil && *v != fleets.CachingTypesNone {
		caching = string(*v)
	}
	output.Caching = caching
	output.DeleteOption = string(pointer.From(input.DeleteOption))
	output.DiskSizeInGB = pointer.From(input.DiskSizeGB)
	output.Name = pointer.From(input.Name)
	output.VhdContainers = pointer.From(input.VhdContainers)

	if v := input.ManagedDisk; v != nil {
		if v := v.DiskEncryptionSet; v != nil {
			output.DiskEncryptionSetId = pointer.From(v.Id)
		}
		if sp := v.SecurityProfile; sp != nil {
			if v := sp.DiskEncryptionSet; v != nil {
				output.SecurityDiskEncryptionSetId = pointer.From(v.Id)
			}
			output.SecurityEncryptionType = string(pointer.From(sp.SecurityEncryptionType))
		}
		output.StorageAccountType = string(pointer.From(v.StorageAccountType))
	}

	output.VhdContainers = pointer.From(input.VhdContainers)
	output.WriteAcceleratorEnabled = pointer.From(input.WriteAcceleratorEnabled)

	return append(outputList, output)
}

func flattenIPConfigurationModel(inputList []fleets.VirtualMachineScaleSetIPConfiguration) []IPConfigurationModel {
	var outputList []IPConfigurationModel
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := IPConfigurationModel{
			Name: input.Name,
		}
		if props := input.Properties; props != nil {
			output.Primary = pointer.From(props.Primary)
			output.Version = string(pointer.From(props.PrivateIPAddressVersion))

			addressPools := make([]string, 0)
			if v := props.ApplicationGatewayBackendAddressPools; v != nil {
				addressPools = flattenSubResourceId(*v)
			}
			output.ApplicationGatewayBackendAddressPoolIds = addressPools

			lbAddressPools := make([]string, 0)
			if v := props.LoadBalancerBackendAddressPools; v != nil {
				lbAddressPools = flattenSubResourceId(*v)
			}
			output.LoadBalancerBackendAddressPoolIds = lbAddressPools

			groupIds := make([]string, 0)
			if v := props.ApplicationSecurityGroups; v != nil {
				groupIds = flattenSubResourceId(*v)
			}
			output.ApplicationSecurityGroupIds = groupIds

			if v := props.PublicIPAddressConfiguration; v != nil {
				output.PublicIPAddress = flattenPublicIPAddressModel(v)
			}

			if v := props.Subnet; v != nil {
				output.SubnetId = pointer.From(v.Id)
			}
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func validateAdminUsernameLinux(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	// **Disallowed values:**
	invalidUserNames := []string{
		" ", "abrt", "adm", "admin", "audio", "backup", "bin", "cdrom", "cgred", "console", "crontab", "daemon", "dbus", "dialout", "dip",
		"disk", "fax", "floppy", "ftp", "fuse", "games", "gnats", "gopher", "haldaemon", "halt", "irc", "kmem", "landscape", "libuuid", "list",
		"lock", "lp", "mail", "maildrop", "man", "mem", "messagebus", "mlocate", "modem", "netdev", "news", "nfsnobody", "nobody", "nogroup",
		"ntp", "operator", "oprofile", "plugdev", "polkituser", "postdrop", "postfix", "proxy", "public", "qpidd", "root", "rpc", "rpcuser",
		"sasl", "saslauth", "shadow", "shutdown", "slocate", "src", "ssh", "sshd", "staff", "stapdev", "stapusr", "sudo", "sync", "sys", "syslog",
		"tape", "tcpdump", "test", "trusted", "tty", "users", "utempter", "utmp", "uucp", "uuidd", "vcsa", "video", "voice", "wheel", "whoopsie",
		"www", "www-data", "wwwrun", "xok",
	}

	for _, str := range invalidUserNames {
		if strings.EqualFold(v, str) {
			errors = append(errors, fmt.Errorf("%q can not be one of %s, got %q", key, azure.QuotedStringSlice(invalidUserNames), v))
			return warnings, errors
		}
	}

	if len(v) < 1 || len(v) > 64 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length, got %q(%d characters)", key, v, len(v)))
		return warnings, errors
	}

	return
}

func validatePasswordComplexityWindows(input interface{}, key string) (warnings []string, errors []error) {
	return validatePasswordComplexity(input, key, 8, 123)
}

func validatePasswordComplexityLinux(input interface{}, key string) (warnings []string, errors []error) {
	return validatePasswordComplexity(input, key, 6, 72)
}

func validatePasswordComplexity(input interface{}, key string, min int, max int) (warnings []string, errors []error) {
	password, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return warnings, errors
	}

	complexityMatch := 0
	re := regexp.MustCompile(`[a-z]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[A-Z]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[0-9]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	re = regexp.MustCompile(`[\W_]{1,}`)
	if re != nil && re.MatchString(password) {
		complexityMatch++
	}

	if complexityMatch < 3 {
		errors = append(errors, fmt.Errorf("%q did not meet minimum password complexity requirements. A password must contain at least 3 of the 4 following conditions: a lower case character, a upper case character, a digit and/or a special character. Got %q", key, password))
		return warnings, errors
	}

	if len(password) < min || len(password) > max {
		errors = append(errors, fmt.Errorf("%q must be at least 6 characters long and less than 72 characters long. Got %q(%d characters)", key, password, len(password)))
		return warnings, errors
	}

	// NOTE: I realize that some of these will not pass the above complexity checks, but they are in the API so I am checking
	// the same values that the API is...
	disallowedValues := []string{
		"abc@123", "P@$$w0rd", "P@ssw0rd", "P@ssword123", "Pa$$word", "pass@word1", "Password!", "Password1", "Password22", "iloveyou!",
	}

	for _, str := range disallowedValues {
		if password == str {
			errors = append(errors, fmt.Errorf("%q can not be one of %s, got %q", key, azure.QuotedStringSlice(disallowedValues), password))
			return warnings, errors
		}
	}

	return warnings, errors
}

func validateAdminUsernameWindows(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	// **Disallowed values:**
	invalidUserNames := []string{
		" ", "administrator", "admin", "user", "user1", "test", "user2", "test1", "user3", "admin1", "1", "123", "a",
		"actuser", "adm", "admin2", "aspnet", "backup", "console", "david", "guest", "john", "owner", "root", "server",
		"sql", "support", "support_388945a0", "sys", "test2", "test3", "user4", "user5",
	}

	for _, str := range invalidUserNames {
		if strings.EqualFold(v, str) {
			errors = append(errors, fmt.Errorf("%q can not be one of %v, got %q", key, invalidUserNames, v))
			return warnings, errors
		}
	}

	// Cannot end in "."
	if strings.HasSuffix(input.(string), ".") {
		errors = append(errors, fmt.Errorf("%q can not end with a '.', got %q", key, v))
		return warnings, errors
	}

	if len(v) < 1 || len(v) > 20 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 20 characters in length, got %q(%d characters)", key, v, len(v)))
		return warnings, errors
	}

	return
}
