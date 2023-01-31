package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/diff"
)

// pass a list of changed file to check the document
func main() {
	result := diff.DiffAll(diff.AzurermRegistersAll())
	if !result.HasDiff() {
		log.Printf("document linter runs success, time costs: %v", result.CostTime())
		return
	}

	result.FixDocuments()
	log.Printf("%s\n", result.ToString())
	fmt.Printf("document linter runs failed with: 1\n")
	os.Exit(1)
}
