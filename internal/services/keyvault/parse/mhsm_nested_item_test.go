package parse

import "testing"

func TestNewMHSMNestedItemID(t *testing.T) {
	cases := []struct {
		Scenario        string
		keyVaultBaseUrl string
		Expected        string
		Scope           string
		Name            string
		ExpectError     bool
	}{
		{
			Scenario:        "empty values",
			keyVaultBaseUrl: "",
			Expected:        "",
			ExpectError:     true,
		},
		{
			Scenario:        "valid, no port",
			keyVaultBaseUrl: "https://test.managedhsm.azure.net",
			Scope:           "/",
			Name:            "test",
			Expected:        "https://test.managedhsm.azure.net///test",
			ExpectError:     false,
		},
		{
			Scenario:        "valid, with port",
			keyVaultBaseUrl: "https://test.managedhsm.azure.net:443",
			Scope:           "/keys",
			Name:            "test",
			Expected:        "https://test.managedhsm.azure.net//keys/test",
			ExpectError:     false,
		},
	}
	for idx, tc := range cases {
		id, err := NewMHSMNestedItemID(tc.keyVaultBaseUrl, tc.Scope, RoleAssignmentType, tc.Name)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Resource ID '%s': %+v", tc.keyVaultBaseUrl, err)
				return
			}
			continue
		}
		if id.ID() != tc.Expected {
			t.Fatalf("Expected %d id for %q to be %q, got %q", idx, tc.keyVaultBaseUrl, tc.Expected, id)
		}
	}
}

func TestParseMHSMNestedItemID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    MHSMNestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net///test",
			ExpectError: true,
			Expected: MHSMNestedItemId{
				Name:         "test",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net///bird",
			ExpectError: true,
			Expected: MHSMNestedItemId{
				Name:         "bird",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net///bird",
			ExpectError: false,
			Expected: MHSMNestedItemId{
				Name:         "bird",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/",
			},
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net//keys/world",
			ExpectError: false,
			Expected: MHSMNestedItemId{
				Name:         "world",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/keys",
			},
		},
		{
			Input:       "https://my-keyvault.managedhsm.azure.net//keys/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
			Expected: MHSMNestedItemId{
				Name:         "fdf067c93bbb4b22bff4d8b7a9a56217",
				VaultBaseUrl: "https://my-keyvault.managedhsm.azure.net/",
				Scope:        "/keys",
			},
		},
	}

	for idx, tc := range cases {
		secretId, err := ParseMHSMNestedItemID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if secretId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.VaultBaseUrl != secretId.VaultBaseUrl {
			t.Fatalf("Expected %d 'KeyVaultBaseUrl' to be '%s', got '%s' for ID '%s'", idx, tc.Expected.VaultBaseUrl, secretId.VaultBaseUrl, tc.Input)
		}

		if tc.Expected.Name != secretId.Name {
			t.Fatalf("Expected 'Name' to be '%s', got '%s' for ID '%s'", tc.Expected.Name, secretId.Name, tc.Input)
		}

		if tc.Expected.Scope != secretId.Scope {
			t.Fatalf("Expected 'Scope' to be '%s', got '%s' for ID '%s'", tc.Expected.Scope, secretId.Scope, tc.Input)
		}

		if tc.Input != secretId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, secretId.ID())
		}
	}
}
