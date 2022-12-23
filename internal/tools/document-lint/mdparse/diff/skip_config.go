package diff

import "strings"

var skipProps = []string{
	"",
	"all.serialization.field_delimiter",
	"all.timezone",
	"all.time_zone",
	"all.time_zone_id",
	"azurerm_virtual_network_gateway_connection.type",
	"azurerm_nginx_deployment.identity.type",
	"azurerm_managed_disk.create_option",
	"azurerm_synapse_role_assignment.role_name",
	"all.advanced_filter",
}

var skipPropMap = map[string]struct{}{}

func init() {
	for _, k := range skipProps {
		skipPropMap[k] = struct{}{}
	}
}

func isSkipProp(rt, prop string) bool {
	if _, ok := skipPropMap[rt]; ok {
		return true
	}
	if _, ok := skipPropMap[rt+"."+prop]; ok {
		return true
	}
	if idx := strings.LastIndex(prop, "."); idx > 0 {
		prop = prop[idx+1:]
	}
	if _, ok := skipPropMap["all."+prop]; ok {
		return true
	}
	return false
}
