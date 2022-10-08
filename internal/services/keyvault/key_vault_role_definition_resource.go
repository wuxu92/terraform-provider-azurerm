package keyvault

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/sdk/v7.3/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Permission struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}

func (p *Permission) toSDKDataAction() (pda, pnda *[]keyvault.DataAction) {
	var da, nda []keyvault.DataAction
	for _, d := range p.DataActions {
		da = append(da, keyvault.DataAction(d))
	}
	for _, nd := range p.NotDataActions {
		nda = append(nda, keyvault.DataAction(nd))
	}
	return &da, &nda
}

func (p *Permission) loadSDKDataAction(perm keyvault.Permission) Permission {
	p.Actions = pointer.ToSliceOfStrings(perm.Actions)
	p.NotActions = pointer.ToSliceOfStrings(perm.NotActions)
	if perm.DataActions != nil {
		for _, a := range *perm.DataActions {
			p.DataActions = append(p.DataActions, string(a))
		}
	}
	if perm.NotDataActions != nil {
		for _, a := range *perm.NotDataActions {
			p.NotDataActions = append(p.NotDataActions, string(a))
		}
	}
	return *p
}

type KeyVaultRoleDefinitionModel struct {
	Name             string       `tfschema:"name"`
	RoleDefinitionId string       `tfschema:"role_definition_id"`
	Scope            string       `tfschema:"scope"`
	VaultBaseUrl     string       `tfschema:"vault_base_url"`
	Description      string       `tfschema:"description"`
	AssignableScopes []string     `tfschema:"assignable_scopes"`
	Permission       []Permission `tfschema:"permission"`
	RoleType         string       `tfschema:"role_type"`
	ResourceId       string       `tfschema:"resource_id"`
}

func (k KeyVaultRoleDefinitionModel) id() string {
	return fmt.Sprintf("%s/%s/%s", k.VaultBaseUrl, k.Scope, k.RoleDefinitionId)
}

func (k *KeyVaultRoleDefinitionModel) ToSDKPermissions() *[]keyvault.Permission {
	var res []keyvault.Permission
	for _, p := range k.Permission {
		ins := keyvault.Permission{
			Actions:    pointer.FromSliceOfStrings(p.Actions),
			NotActions: pointer.FromSliceOfStrings(p.NotActions),
		}
		ins.DataActions, ins.NotDataActions = p.toSDKDataAction()
		res = append(res, ins)
	}
	return &res
}

func (k *KeyVaultRoleDefinitionModel) LoadSDKPermission(perms *[]keyvault.Permission) {
	if perms != nil {
		k.Permission = []Permission{}
		for _, p := range *perms {
			k.Permission = append(k.Permission, (&Permission{}).loadSDKDataAction(p))
		}
	}
}

type KeyVaultRoleDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = (*KeyVaultRoleDefinitionResource)(nil)

func (k KeyVaultRoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"vault_base_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Default:  "/",
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"/",
			}, false),
		},

		//lintignore:XS003
		"permission": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"not_actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range keyvault.PossibleDataActionValues() {
									res = append(res, string(v))
								}
								return
							}(), false),
						},
						Set: pluginsdk.HashString,
					},

					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range keyvault.PossibleDataActionValues() {
									res = append(res, string(v))
								}
								return
							}(), false),
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},

		"assignable_scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.Global),
					string(keyvault.Keys),
				}, false),
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (k KeyVaultRoleDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (k KeyVaultRoleDefinitionResource) ModelObject() interface{} {
	return &KeyVaultRoleDefinitionModel{}
}

func (k KeyVaultRoleDefinitionResource) ResourceType() string {
	return "azurerm_key_vault_role_definition"
}

func (k KeyVaultRoleDefinitionResource) createOrUpdateFunc(isUpdate bool) sdk.ResourceRunFunc {
	return func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
		client := meta.Client.KeyVault.MHSMRoleClient

		var model KeyVaultRoleDefinitionModel
		if err = meta.Decode(&model); err != nil {
			return err
		}

		if model.RoleDefinitionId == "" {
			model.RoleDefinitionId, err = uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
			}
		}

		if model.Name == "" {
			model.Name = model.RoleDefinitionId
		}

		id, err := parse.NewMHSMNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleDefinitionType, model.RoleDefinitionId)
		if err != nil {
			return err
		}

		//subscriptionID := meta.Client.Account.SubscriptionId
		existing, err := client.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
		if !utils.ResponseWasNotFound(existing.Response) {
			if err != nil {
				return fmt.Errorf("retreiving role definition by name %s: %v", model.RoleDefinitionId, err)
			}
			if !isUpdate {
				return meta.ResourceRequiresImport(k.ResourceType(), id)
			}
		} else if isUpdate {
			return fmt.Errorf("not found resource to update: %s", id)
		}

		var param keyvault.RoleDefinitionCreateParameters
		param.Properties = &keyvault.RoleDefinitionProperties{}
		prop := param.Properties
		prop.RoleName = utils.String(model.Name)
		prop.Description = utils.String(model.Description)
		prop.RoleType = keyvault.BuiltInRole
		prop.Permissions = model.ToSDKPermissions()

		var scopes []keyvault.RoleScope
		for _, role := range model.AssignableScopes {
			scopes = append(scopes, keyvault.RoleScope(role))
		}
		if len(scopes) > 0 {
			prop.AssignableScopes = &scopes
		}

		_, err = client.CreateOrUpdate(ctx, model.VaultBaseUrl, model.Scope, model.RoleDefinitionId, param)
		if err != nil {
			return fmt.Errorf("creating %s: %v", model.id(), err)
		}

		// has to set role_definition id before Read
		meta.ResourceData.Set("role_definition_id", model.RoleDefinitionId)

		meta.SetID(id)
		return nil
	}
}

func (k KeyVaultRoleDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    k.createOrUpdateFunc(false),
	}
}

func (k KeyVaultRoleDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			var model KeyVaultRoleDefinitionModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			if model.RoleDefinitionId == "" {
				return fmt.Errorf("role definition id is required")
			}

			id, err := parse.NewMHSMNestedItemID(model.VaultBaseUrl, model.Scope, parse.RoleDefinitionType, model.RoleDefinitionId)
			if err != nil {
				return err
			}

			client := meta.Client.KeyVault.MHSMRoleClient
			result, err := client.Get(ctx, model.VaultBaseUrl, model.Scope, model.RoleDefinitionId)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			prop := result
			model.Name = pointer.ToString(prop.RoleName)
			model.RoleDefinitionId = pointer.ToString(prop.Name) // prop.Name is role definition name
			model.Description = pointer.ToString(prop.Description)
			model.RoleType = string(prop.RoleType)
			model.ResourceId = pointer.ToString(prop.ID)

			if prop.AssignableScopes != nil {
				for _, r := range *prop.AssignableScopes {
					model.AssignableScopes = append(model.AssignableScopes, string(r))
				}
			}

			model.LoadSDKPermission(prop.Permissions)

			meta.SetID(id)
			return meta.Encode(&model)
		},
	}
}

func (k KeyVaultRoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 10,
		Func:    k.createOrUpdateFunc(true),
	}
}

func (k KeyVaultRoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			var model KeyVaultRoleDefinitionModel
			if err := meta.Decode(&model); err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", model.id())
			client := meta.Client.KeyVault.MHSMRoleClient
			if _, err := client.Delete(ctx, model.VaultBaseUrl, model.Scope, model.RoleDefinitionId); err != nil {
				return fmt.Errorf("deleting %s: %v", model.id(), err)
			}
			return nil
		},
	}
}

func (k KeyVaultRoleDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.MHSMNestedItemId
}
