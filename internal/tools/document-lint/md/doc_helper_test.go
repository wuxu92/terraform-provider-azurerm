// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package md_test

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
)

func TestMDPathFor(t *testing.T) {
	cases := [][2]string{
		{
			"azurerm_api_management_api_policy",
			"api_management_api_policy.html.markdown",
		},
		{
			"not_exists",
			"",
		},
	}
	for _, c := range cases {
		got := md.MDPathFor(c[0])
		if !strings.Contains(got, c[1]) {
			t.Fatalf("%s: \nwant: %s,\ngot:  %s", c[0], c[1], got)
		}
	}
}

func TestResourceNameReg(t *testing.T) {
	var titleReg = regexp.MustCompile(`\npage_title:[^\n]*(azurerm_[a-zA-Z0-9_]+)"`)

	subs := titleReg.FindStringSubmatch(`---
subcategory: "AAD B2C"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_aadb2c_directory"
description: |-
  Manages an AAD B2C Directory.
---

# azurerm_aadb2c_directory

Manages an AAD B2C Directory.

## Example Usage`)
	t.Logf("%v", subs)
}

func TestAddDocumentPathTo(t *testing.T) {
	t.Skip()
	schemaJSONFile := "/home/wuxu/azure/terraform-provider-azurerm/internal/tools/schema-api/azurerm-schema.json"
	f, _ := os.Open(schemaJSONFile)
	bs, _ := io.ReadAll(f)
	var schema map[string]interface{}
	json.Unmarshal(bs, &schema)
	resources := schema["providerSchema"].(map[string]interface{})["resources"].(map[string]interface{})

	for name, value := range resources {
		mdPath := md.MDPathFor(name)
		value.(map[string]interface{})["document_path"] = mdPath
	}

	bs, _ = json.Marshal(schema)
	err := os.WriteFile(schemaJSONFile, bs, 066)
	if err != nil {
		t.Fatal(err)
	}
}
