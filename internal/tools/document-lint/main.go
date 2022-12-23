package main

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/diff"
)

// pass a list of changed file to check the document
func main() {
	result := diff.DiffAll(diff.AzurermRegistersAll())
	if len(result.GetResult()) == 0 {
		return
	}

	result.FixDocuments()
	log.Printf("%s", result.ToString())
	log.Printf("cost %s", result.CostTime())
	os.Exit(1)
}
