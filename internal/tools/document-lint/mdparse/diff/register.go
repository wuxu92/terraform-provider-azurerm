package diff

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn"
)

type Registers struct {
	Registers []interface{}
}

func (r *Registers) add(i interface{}) {
	r.Registers = append(r.Registers, i)
}

func AzurermRegisters() Registers {
	res := Registers{
		Registers: []interface{}{
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
