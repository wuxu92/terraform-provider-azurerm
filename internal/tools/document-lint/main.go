package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/mdparse/check"
)

func printHelp() {
	text := `USAGE: azdoc-lint [CMD] [OPTIONS]
CMD:
  check:	check documents and print the error information
  fix:	 	check and try to fix existing errors

OPTIONS:
`
	fmt.Printf("%s\n", text)
}

var config = struct {
	cmd          string
	dryRun       bool
	resource     string
	rp           string
	skipResource string
	skipRP       string
}{
	cmd:    "check",
	dryRun: true,
}

func parseArgs() {
	fs := flag.NewFlagSet("azdoc-check", flag.ExitOnError)
	fs.StringVar(&config.resource, "resource", os.Getenv("ONLY_RESOURCE"), "a list of resource names to check")
	fs.StringVar(&config.rp, "services", os.Getenv("ONLY_SERVICE"), "a list of services names to check")
	fs.StringVar(&config.skipResource, "skip-resource", os.Getenv("SKIP_RESOURCE"), "a list of resource names to skip the check")
	fs.StringVar(&config.skipRP, "skip-services", os.Getenv("SKIP_SERVICE"), "a list of rp names to skip the check")

	fs.Usage = func() {
		printHelp()
		fs.PrintDefaults()
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "check":
			config.cmd = "check"
			_ = fs.Parse(os.Args[2:])
		case "fix":
			config.cmd = "fix"
			config.dryRun = false
			_ = fs.Parse(os.Args[2:])
		default:
			fs.Usage()
		}
	}
	// update config in check package
	check.SetConfig(config.resource, config.rp, config.skipResource, config.skipRP, config.dryRun)
}

// pass a list of changed file to check the document
func main() {
	parseArgs()
	result := check.DiffAll(check.AzurermRegistersAll())
	if !result.HasDiff() {
		log.Printf("document linter runs success, time costs: %v", result.CostTime())
		return
	}

	if config.cmd == "fix" {
		if err := result.FixDocuments(); err != nil {
			log.Fatalf("error occurs when trying to fix documents: %v", err)
		}
	}
	log.Printf("%s\n", result.ToString())
	// fmt.Printf("document linter runs failed with: 1\n")
	os.Exit(1)
}
