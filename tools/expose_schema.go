package tools

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type TypedSDK = sdk.Resource

type TypedRegister = sdk.TypedServiceRegistration

type UntypedRegister = sdk.UntypedServiceRegistration

var GetFuncPossibleValues = validation.GetFunctionValues

type Resources struct {
	Typed   []sdk.Resource
	Untyped map[string]*schema.Resource // map from resource name to UntypedResource
}

func ResourceForSDKType(res TypedSDK) *schema.Resource {
	r := sdk.NewResourceWrapper(res)
	ins, _ := r.Resource()
	return ins
}

func AzurermProvider() *Resources {
	ng := automation.SoftwareUpdateConfigurationResource{}
	r := &Resources{
		Typed: []sdk.Resource{
			ng,
		},
	}
	return r
}

type Registers struct {
	//Typed       []sdk.TypedServiceRegistration
	//Untyped     []sdk.UntypedServiceRegistration
	Registers []interface{}
}

func (r *Registers) add(i interface{}) {
	r.Registers = append(r.Registers, i)
}

func AzurermRegisters() Registers {
	res := Registers{
		Registers: []interface{}{
			//automation.Registration{},
			//consumption.Registration{},
			//monitor.Registration{},
			//authorization.Registration{},
			//compute.Registration{},
			//costmanagement.Registration{},
			//loganalytics.Registration{},
			//network.Registration{},
			cdn.Registration{},
		},
	}
	return res
}

func AzurermRegistersAll() Registers {
	res := Registers{
		Registers: []interface{}{},
	}
	for _, r := range provider.SupportedTypedServices() {
		res.add(r)
	}

	for _, r := range provider.SupportedUntypedServices() {
		res.add(r)
	}
	return res
}
