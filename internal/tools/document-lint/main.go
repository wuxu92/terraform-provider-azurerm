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
	if len(result.GetResult()) == 0 {
		log.Printf("linter runs success")
		return
	}

	result.FixDocuments()
	log.Printf("%s\n", result.ToString())
	fmt.Printf("linter exists status: 1\n")
	os.Exit(1)
}
