package schema

//
//var identityTypeValues = map[string][]string{
//	"UserAssignedIdentityRequiredForceNew": {"UserAssigned"},
//	"UserAssignedIdentityOptional":         {"UserAssigned"},
//	"UserAssignedIdentityOptionalForceNew": {"UserAssigned"},
//	"UserAssignedIdentityComputed":         {"UserAssigned"},
//
//	"SystemAssignedIdentityRequiredForceNew": {"SystemAssigned"},
//	"SystemAssignedIdentityOptional":         {"SystemAssigned"},
//	"SystemAssignedIdentityOptionalForceNew": {"SystemAssigned"},
//	"SystemAssignedIdentityComputed":         {"SystemAssigned"},
//
//	"SystemOrUserAssignedIdentityRequired":         {"SystemAssigned", "UserAssigned"},
//	"SystemOrUserAssignedIdentityRequiredForceNew": {"SystemAssigned", "UserAssigned"},
//	"SystemOrUserAssignedIdentityOptional":         {"SystemAssigned", "UserAssigned"},
//	"SystemOrUserAssignedIdentityOptionalForceNew": {"SystemAssigned", "UserAssigned"},
//
//	"SystemAssignedUserAssignedIdentityRequired":         {"SystemAssigned", "UserAssigned", "SystemAssigned, UserAssigned"},
//	"SystemAssignedUserAssignedIdentityRequiredForceNew": {"SystemAssigned", "UserAssigned", "SystemAssigned, UserAssigned"},
//	"SystemAssignedUserAssignedIdentityOptional":         {"SystemAssigned", "UserAssigned", "SystemAssigned, UserAssigned"},
//	"SystemAssignedUserAssignedIdentityOptionalForceNew": {"SystemAssigned", "UserAssigned", "SystemAssigned, UserAssigned"},
//}
//
//var specialProps = map[string]map[string][]string{
//	"identity.type": identityTypeValues,
//}
//
//func IdentityTypePossibleValues(typ string) []string {
//	return identityTypeValues[typ]
//}
//
//func isSpecialProp(key string) bool {
//	_, ok := specialProps[key]
//	return ok
//}
