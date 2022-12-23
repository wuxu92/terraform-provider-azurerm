# 

# henglu

## Purview

### azurerm_purview_account

- [ ] identity.identity_ids missed in doc

## Datadog

### azurerm_datadog_monitor_tag_rule

- [ ] name missed in doc
- [ ] rule_set_name missed in code

## Data Factory

### azurerm_data_factory_linked_service_web

- [ ] password missed in doc
- [ ] username missed in doc

### azurerm_data_factory

- [ ] purview_id missed in doc

### azurerm_data_factory_data_flow

- [ ] sink.flowlet.dataset_parameters missed in doc
- [ ] source.flowlet.dataset_parameters missed in doc
- [ ] source.rejected_linked_service missed in doc
- [ ] transformation.flowlet.dataset_parameters missed in doc

### azurerm_data_factory_linked_service_azure_file_storage

- [ ] host missed in doc
- [ ] password missed in doc
- [ ] user_id missed in doc

### azurerm_data_factory_linked_service_azure_databricks

- [ ] new_cluster_config.max_number_of_workers missed in doc
- [ ] new_cluster_config.min_number_of_workers missed in doc

### azurerm_data_factory_flowlet_data_flow

- [ ] sink.flowlet.dataset_parameters missed in doc
- [ ] source.flowlet.dataset_parameters missed in doc
- [ ] source.rejected_linked_service missed in doc
- [ ] transformation.flowlet.dataset_parameters missed in doc
- [ ] transformation.name missed in doc

### azurerm_data_factory_dataset_parquet

- [ ] compression_level missed in doc

## Spring Cloud

### azurerm_spring_cloud_service

- [ ] required_network_traffic_rules.ips missed in code

## Container Services

### azurerm_container_registry_token_password

- [ ] password missed in code

### azurerm_kubernetes_fleet_manager

- [ ] hub_profile not block missed in doc

### azurerm_container_registry_scope_map

- [ ] description missed in doc

### azurerm_kubernetes_cluster

- [ ] default_node_pool.host_group_id missed in doc
- [ ] default_node_pool.node_taints missed in doc
- [ ] default_node_pool.proximity_placement_group_id missed in doc
- [ ] enable_pod_security_policy missed in doc
- [ ] ingress_application_gateway.ingress_application_gateway_identity not block missed in doc
- [ ] key_vault_secrets_provider.secret_identity not block missed in doc
- [ ] linux_profile.ssh_key not block missed in doc
- [ ] monitor_metrics.labels_allowed missed in doc
- [ ] oms_agent.oms_agent_identity not block missed in doc

# jiaweitao

## HDInsight

### azurerm_hdinsight_hadoop_cluster

- [ ] compute_isolation.compute_isolation_enabled missed in doc
- [ ] compute_isolation.enable_compute_isolation missed in code
- [ ] disk_encryption missed in doc
- [ ] min_tls_version missed in code
- [ ] roles.edge_node.https_endpoints missed in doc
- [ ] roles.edge_node.target_instance_count missed in doc
- [ ] roles.edge_node.uninstall_script_actions missed in doc
- [ ] roles.head_node.script_actions not block missed in doc
- [ ] roles.worker_node.script_actions missed in doc
- [ ] roles.zookeeper_node.script_actions missed in doc
- [ ] tls_min_version missed in doc

### azurerm_hdinsight_hbase_cluster

- [ ] compute_isolation.compute_isolation_enabled missed in doc
- [ ] compute_isolation.enable_compute_isolation missed in code
- [ ] disk_encryption missed in doc
- [ ] min_tls_version missed in code
- [ ] roles.head_node.script_actions not block missed in doc
- [ ] roles.worker_node.script_actions missed in doc
- [ ] roles.zookeeper_node.script_actions missed in doc
- [ ] tls_min_version missed in doc

### azurerm_hdinsight_interactive_query_cluster

- [ ] component_version.interactive_hive missed in doc
- [ ] component_version.interactive_query missed in code
- [ ] compute_isolation.compute_isolation_enabled missed in doc
- [ ] compute_isolation.enable_compute_isolation missed in code
- [ ] disk_encryption missed in doc
- [ ] min_tls_version missed in code
- [ ] roles.head_node.script_actions not block missed in doc
- [ ] roles.worker_node.script_actions missed in doc
- [ ] roles.zookeeper_node.script_actions missed in doc
- [ ] tls_min_version missed in doc

### azurerm_hdinsight_kafka_cluster

- [ ] compute_isolation.compute_isolation_enabled missed in doc
- [ ] compute_isolation.enable_compute_isolation missed in code
- [ ] disk_encryption missed in doc
- [ ] roles.head_node.script_actions missed in doc
- [ ] roles.kafka_management_node.script_actions missed in doc
- [ ] roles.worker_node.script_actions missed in doc
- [ ] roles.zookeeper_node.script_actions missed in doc
- [ ] storage_account.storage_resource_id missed in doc

### azurerm_hdinsight_spark_cluster

- [ ] compute_isolation.compute_isolation_enabled missed in doc
- [ ] compute_isolation.enable_compute_isolation missed in code
- [ ] disk_encryption missed in doc
- [ ] min_tls_version missed in code
- [ ] roles.head_node.script_actions not block missed in doc
- [ ] roles.worker_node.script_actions missed in doc
- [ ] roles.zookeeper_node.script_actions missed in doc
- [ ] tls_min_version missed in doc

## Stream Analytics

### azurerm_stream_analytics_job

- [ ] account_key missed in code
- [ ] account_name missed in code
- [ ] authentication_mode missed in code
- [ ] job_storage_account not block missed in doc

## ServiceConnector

### azurerm_app_service_connection

- [ ] authentication not block missed in doc
- [ ] certificate missed in code
- [ ] client_id missed in code
- [ ] principal_id missed in code
- [ ] secret missed in code
- [ ] subscription_id missed in code
- [ ] type missed in code

### azurerm_spring_cloud_connection

- [ ] authentication not block missed in doc
- [ ] certificate missed in code
- [ ] client_id missed in code
- [ ] principal_id missed in code
- [ ] secret missed in code
- [ ] subscription_id missed in code
- [ ] type missed in code

## Dev Test

### azurerm_dev_test_global_vm_shutdown_schedule

- [ ] notification_settings missed in doc

### azurerm_dev_test_schedule

- [ ] daily_recurrence missed in doc
- [ ] hourly_recurrence missed in doc
- [ ] notification_settings missed in doc
- [ ] weekly_recurrence missed in doc

### azurerm_dev_test_policy

- [ ] location missed in code

# v-cheye

## Confidential Ledger

### azurerm_confidential_ledger

- [ ] azuread_based_service_principal missed in doc
- [ ] azuread_service_principal missed in code
- [ ] cert_based_security_principals missed in code
- [ ] certificate_based_security_principal missed in doc

## CosmosDB

### azurerm_cosmosdb_mongo_collection

- [ ] account_name missed in doc

### azurerm_cosmosdb_account

- [ ] capabilities not block missed in doc
- [ ] consistency_policy not block missed in doc
- [ ] virtual_network_rule not block missed in doc

## Log Analytics

### azurerm_log_analytics_cluster

- [ ] type missed in code

## NetApp

### azurerm_netapp_volume

- [ ] azure_vmware_data_store_enabled missed in doc

### azurerm_netapp_snapshot_policy

- [ ] daily_schedule not block missed in doc
- [ ] hourly_schedule not block missed in doc
- [ ] monthly_schedule not block missed in doc
- [ ] tags missed in doc
- [ ] weekly_schedule not block missed in doc

## Service Fabric

### azurerm_service_fabric_cluster

- [ ] client_certificate_common_name.certificate_issuer_thumbprint missed in code
- [ ] client_certificate_common_name.issuer_thumbprint missed in doc
- [ ] upgrade_policy.force_restart missed in code
- [ ] upgrade_policy.force_restart_enabled missed in doc

## Cost Management

### azurerm_billing_account_cost_management_export

- [ ] recurrence_period_end missed in code
- [ ] recurrence_period_end_date missed in doc

### azurerm_subscription_cost_management_export

- [ ] recurrence_period_end missed in code
- [ ] recurrence_period_end_date missed in doc

# v-elenaxin

## Media

### azurerm_media_asset_filter

- [ ] track_selection.condition not block missed in doc

### azurerm_media_streaming_endpoint

- [ ] access_control.ip_allow not block missed in doc

### azurerm_media_content_key_policy

- [ ] policy_option.fairplay_configuration.offline_rental_configuration not block missed in doc

### azurerm_media_live_event_output

- [ ] output_snap_time_in_seconds missed in doc
- [ ] output_snap_timestamp_in_seconds missed in code

## SQL

### azurerm_sql_database

- [ ] extended_auditing_policy missed in code

### azurerm_sql_server

- [ ] principal_id missed in code
- [ ] tenant_id missed in code

## Analysis Services

### azurerm_analysis_services_server

- [ ] tags missed in doc

## Microsoft SQL Server / Azure SQL

### azurerm_mssql_managed_instance

- [ ] identity.identity_ids missed in doc

### azurerm_mssql_virtual_machine

- [ ] storage_configuration.data_settings not block missed in doc
- [ ] storage_configuration.log_settings not block missed in doc

### azurerm_mssql_server

- [ ] principal_id missed in code
- [ ] tenant_id missed in code

### azurerm_mssql_managed_instance_vulnerability_assessment

- [ ] managed_instance_id missed in doc
- [ ] manged_instance_id missed in code

### azurerm_mssql_database

- [ ] :L126 missed in code

# wangta

## Monitor

### azurerm_monitor_scheduled_query_rules_log

- [ ] authorized_resource_ids missed in doc

### azurerm_monitor_scheduled_query_rules_alert

- [ ] query_type missed in doc
- [ ] trigger not block missed in doc

## App Configuration

### azurerm_app_configuration_feature

- [ ] end missed in code
- [ ] start missed in code
- [ ] timewindow_filter.default_rollout_percentage missed in code
- [ ] timewindow_filter.end missed in doc
- [ ] timewindow_filter.groups missed in code
- [ ] timewindow_filter.start missed in doc
- [ ] timewindow_filter.users missed in code

### azurerm_app_configuration

- [ ] encrption missed in code
- [ ] encryption missed in doc

## DataProtection

### azurerm_data_protection_backup_policy_disk

- [ ] resource_group_name missed in code

### azurerm_data_protection_backup_vault

- [ ] principal_id missed in code
- [ ] tenant_id missed in code

# xiaxin.yi

## EventGrid

### azurerm_eventgrid_event_subscription

- [ ] topic_name deprecated missed in code

## EventHub

### azurerm_eventhub_namespace_disaster_recovery_config

- [ ] wait_for_replication missed in code

### azurerm_eventhub_namespace

- [ ] identity.identity_ids missed in doc
- [ ] network_rulesets. missed in code

### azurerm_eventhub_namespace_schema_group

- [ ] name missed in doc

## Health Care

### azurerm_healthcare_service

- [ ] access_policy_ids missed in code
- [ ] access_policy_object_ids missed in doc

### azurerm_healthcare_workspace

- [ ] tags missed in doc

### azurerm_healthcare_dicom_service

- [ ] tags missed in doc

### azurerm_healthcare_fhir_service

- [ ] authentication.smart_proxy_enabled missed in doc
- [ ] resource_group_name missed in doc
- [ ] tags missed in doc

### azurerm_healthcare_medtech_service

- [ ] tags missed in doc

### azurerm_healthcare_medtech_service_fhir_destination

- [ ] destination_fhir_service_id missed in doc

## Web

### azurerm_app_service_slot

- [ ] auth_settings.twitter not block missed in doc
- [ ] logs.http_logs.azure_blob_storage.level missed in code
- [ ] site_config.ip_restriction.subnet_mask missed in code
- [ ] site_config.scm_ip_restriction not block missed in doc

### azurerm_app_service

- [ ] auth_settings.twitter not block missed in doc
- [ ] logs.http_logs.azure_blob_storage.level missed in code
- [ ] site_config.auto_swap_slot_name missed in doc

### azurerm_function_app

- [ ] site_config.auto_swap_slot_name missed in doc

### azurerm_function_app_slot

- [ ] site_config.health_check_path missed in doc
- [ ] site_config.java_version missed in doc
- [ ] site_config.scm_use_main_ip_restriction missed in doc

### azurerm_app_service_certificate

- [ ] hosting_environment_profile_id missed in code
- [ ] tags missed in doc

### azurerm_app_service_plan

- [ ] is_xenon missed in doc

## SignalR

### azurerm_web_pubsub

- [ ] ip_address missed in code

## AppService

### azurerm_source_control_token

- [ ] token_secret missed in doc

### azurerm_linux_function_app

- [ ] site_config.application_stack.docker not block missed in doc

### azurerm_linux_web_app

- [ ] logs.application_logs.azure_blob_storage not block missed in doc
- [ ] logs.http_logs.azure_blob_storage.level missed in code

### azurerm_linux_web_app_slot

- [ ] :L57 missed in code
- [ ] connection_string.name missed in doc
- [ ] logs.application_logs.azure_blob_storage not block missed in doc
- [ ] logs.http_logs.azure_blob_storage.level missed in code
- [ ] site_config.websockets missed in code
- [ ] site_config.websockets_enabled missed in doc

### azurerm_app_service_source_control

- [ ] github_action_configuration.generate_workflow_file missed in doc

### azurerm_windows_web_app

- [ ] logs.application_logs.azure_blob_storage not block missed in doc
- [ ] logs.http_logs.azure_blob_storage.level missed in code
- [ ] site_config.api_definition_url missed in doc
- [ ] site_config.auto_heal_setting.action.custom_action not block missed in doc

### azurerm_windows_web_app_slot

- [ ] :L57 missed in code
- [ ] app_metadata missed in code
- [ ] connection_string.name missed in doc
- [ ] logs.application_logs.azure_blob_storage not block missed in doc
- [ ] logs.http_logs.azure_blob_storage.level missed in code
- [ ] site_config.auto_heal_setting.action.custom_action not block missed in doc
- [ ] site_config.websockets missed in code
- [ ] site_config.websockets_enabled missed in doc

## ServiceBus

### azurerm_servicebus_namespace

- [ ] identity.default_primary_connection_string missed in code
- [ ] identity.default_primary_key missed in code
- [ ] identity.default_secondary_connection_string missed in code
- [ ] identity.default_secondary_key missed in code

# xuwu1

## Fluid Relay

### azurerm_fluid_relay_server

- [ ] principal_id missed in code
- [ ] tenant_id missed in code

## Orbital

### azurerm_orbital_spacecraft

- [ ] bandwidth_mhz missed in code
- [ ] center_frequency_mhz missed in code
- [ ] direction missed in code
- [ ] links not block missed in doc
- [ ] polarization missed in code
- [ ] tags missed in doc

### azurerm_orbital_contact_profile

- [ ] links.channels not block missed in doc
- [ ] tags missed in doc

## Network

### azurerm_route_server

- [ ] sku missed in doc
- [ ] tags missed in doc

### azurerm_virtual_network_gateway

- [ ] custom_route missed in doc
- [ ] vpn_client_configuration.revoked_certificate not block missed in doc

### azurerm_subnet

- [ ] address_prefix deprecated missed in code

### azurerm_application_gateway

- [ ] ssl_policy not block missed in doc
- [ ] ssl_profile.ssl_policy not block missed in doc

### azurerm_route_filter

- [ ] rule not block missed in doc

### azurerm_express_route_port

- [ ] link missed in code

### azurerm_nat_gateway

- [ ] public_ip_address_ids deprecated missed in code
- [ ] public_ip_prefix_ids deprecated missed in code

## Automation

### azurerm_automation_software_update_configuration

- [ ] schedule.monthly_occurrence not block missed in doc
- [ ] target.azure_query.tags not block missed in doc

### azurerm_automation_runbook

- [ ] draft.content_link not block missed in doc
- [ ] draft.parameter missed in code
- [ ] draft.parameters missed in doc
- [ ] publish_content_link not block missed in doc

### azurerm_automation_module

- [ ] module_link not block missed in doc

### azurerm_automation_schedule

- [ ] monthly_occurrence not block missed in doc

# xuzhang3

## Machine Learning

### azurerm_machine_learning_workspace

- [ ] encryption missed in doc

## Redis

### azurerm_redis_cache

- [ ] patch_schedule not block missed in doc

## Redis Enterprise

### azurerm_redis_enterprise_cluster

- [ ] version missed in code

## API Management

### azurerm_api_management

- [ ] hostname_configuration.certificate_source missed in code
- [ ] hostname_configuration.certificate_status missed in code
- [ ] hostname_configuration.expiry missed in code
- [ ] hostname_configuration.subject missed in code
- [ ] hostname_configuration.thumbprint missed in code
- [ ] security.disable_backend_ssl30 missed in code
- [ ] security.disable_backend_tls10 missed in code
- [ ] security.disable_backend_tls11 missed in code
- [ ] security.disable_frontend_ssl30 missed in code
- [ ] security.disable_frontend_tls10 missed in code
- [ ] security.disable_frontend_tls11 missed in code

### azurerm_api_management_custom_domain

- [ ] developer_portal.ssl_keyvault_identity_client_id missed in doc
- [ ] gateway.ssl_keyvault_identity_client_id missed in doc
- [ ] management.ssl_keyvault_identity_client_id missed in doc
- [ ] portal.ssl_keyvault_identity_client_id missed in doc
- [ ] scm.ssl_keyvault_identity_client_id missed in doc

### azurerm_api_management_gateway

- [ ] api_management_id missed in doc
- [ ] api_management_name missed in code
- [ ] resource_group_name missed in code

### azurerm_api_management_diagnostic

- [ ] backend_request.data_masking missed in doc
- [ ] backend_response.data_masking missed in doc
- [ ] frontend_request.data_masking missed in doc
- [ ] frontend_response.data_masking missed in doc

# yicma

## IoT Hub

### azurerm_iothub

- [ ] key_name missed in code
- [ ] permissions missed in code
- [ ] primary_key missed in code
- [ ] secondary_key missed in code
- [ ] shared_access_policy not block missed in doc

### azurerm_iothub_enrichment

- [ ] iothub_name missed in doc
- [ ] resource_group_name missed in doc

## Legacy

### azurerm_virtual_machine

- [ ] identity.tenant_id missed in code

### azurerm_virtual_machine_scale_set

- [ ] os_profile_linux_config.ssh_keys not block missed in doc
- [ ] os_profile_secrets.certificate_store missed in code
- [ ] os_profile_secrets.certificate_url missed in code
- [ ] os_profile_secrets.vault_certificates not block missed in doc

## Compute

### azurerm_orchestrated_virtual_machine_scale_set

- [ ] data_disk.disk_encryption_set_id missed in doc
- [ ] data_disk.write_accelerator_enabled missed in doc
- [ ] encryption_at_host_enabled missed in doc
- [ ] extension.settings missed in doc
- [ ] os_disk.write_accelerator_enabled missed in doc
- [ ] os_profile.linux_configuration.secret.certificate.store missed in code
- [ ] os_profile.windows_configuration.automatic_instance_repair missed in code
- [ ] os_profile.windows_configuration.winrm_listener.protocol missed in doc
- [ ] store missed in code
- [ ] url missed in code
- [ ] zone_balance missed in doc

### azurerm_managed_disk

- [ ] Copy missed in code
- [ ] Empty missed in code
- [ ] FromImage missed in code
- [ ] Import missed in code
- [ ] Restore missed in code
- [ ] Upload missed in code

### azurerm_windows_virtual_machine_scale_set

- [ ] data_disk.disk_iops_read_write missed in code
- [ ] data_disk.disk_mbps_read_write missed in code

### azurerm_dedicated_host_group

- [ ] zone missed in doc
- [ ] zones missed in code

### azurerm_linux_virtual_machine_scale_set

- [ ] data_disk.disk_iops_read_write missed in code
- [ ] data_disk.disk_mbps_read_write missed in code

### azurerm_windows_virtual_machine

- [ ] winrm_listener.Protocol missed in code
- [ ] winrm_listener.protocol missed in doc

## Time Series Insights

### azurerm_iot_time_series_insights_access_policy

- [ ] resource_group_name missed in code

# yunliu1

## Batch

### azurerm_batch_certificate

- [ ] thumbprint_algorithm missed in doc

### azurerm_batch_pool

- [ ] container_configuration not block missed in doc
- [ ] network_configuration.endpoint_configuration.network_security_group_rules not block missed in doc
- [ ] start_task.container.registry not block missed in doc
- [ ] stop_pending_resize_operation missed in doc
- [ ] storage_image_reference not block missed in doc

### azurerm_batch_account

- [ ] encryption not block missed in doc

## Private DNS

### azurerm_private_dns_aaaa_record

- [ ] TTL missed in code
- [ ] ttl missed in doc

### azurerm_private_dns_a_record

- [ ] TTL missed in code
- [ ] ttl missed in doc

## LoadTestService

### azurerm_load_test

- [ ] identity not block missed in doc

# zhaoting.weng

## DomainServices

### azurerm_active_directory_domain_service

- [ ] initial_replica_set.replica_set_id missed in code
- [ ] secure_ldap.external_access_ip_address missed in code

## Storage

### azurerm_storage_account

- [ ] principal_id missed in code
- [ ] tenant_id missed in code

### azurerm_storage_management_policy

- [ ] rule.filters not block missed in doc

### azurerm_storage_blob_inventory_policy

- [ ] filter missed in code
- [ ] rules.filter missed in doc

## CDN

### azurerm_cdn_frontdoor_firewall_policy

- [ ] location missed in code

### azurerm_cdn_frontdoor_rule

- [ ] conditions.host_name_condition.negate_condition missed in doc

## FrontDoor

### azurerm_frontdoor_custom_https_configuration

- [ ] custom_https_configuration.custom_https_configuration missed in code
- [ ] custom_https_configuration.custom_https_provisioning_enabled missed in code
- [ ] custom_https_configuration.frontend_endpoint_id missed in code
- [ ] custom_https_provisioning_enabled missed in doc
- [ ] frontend_endpoint_id missed in doc

### azurerm_frontdoor_rules_engine

- [ ] enabled missed in doc
- [ ] rule.action not block missed in doc

### azurerm_frontdoor

- [ ] backend_pool_settings missed in doc
- [ ] provisioning_state missed in code
- [ ] provisioning_substate missed in code

# zhenteng

## Data Share

### azurerm_data_share_dataset_kusto_cluster

- [ ] public_network_access_enabled missed in code

## Sentinel

### azurerm_sentinel_alert_rule_scheduled

- [ ] incident_configuration.grouping.gorup_by_alert_details missed in code
- [ ] incident_configuration.grouping.gorup_by_custom_details missed in code
- [ ] incident_configuration.grouping.group_by_alert_details missed in doc
- [ ] incident_configuration.grouping.group_by_custom_details missed in doc

## Recovery Services

### azurerm_site_recovery_replicated_vm

- [ ] recovery_replication_policy_id missed in doc

# zhhu

## Resources

### azurerm_resource_deployment_script_azure_power_shell

- [ ] identity.identity_ids missed in doc
- [ ] identity.user_assigned_identities missed in code

### azurerm_resource_deployment_script_azure_cli

- [ ] identity.identity_ids missed in doc
- [ ] identity.user_assigned_identities missed in code

