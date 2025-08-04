// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccStringReuseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create und Read Testing - Setter setzt den Wert
			{
				Config: `
resource "reuse_string" "test" {
    set_if_not_null_or_empty = "initial"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"reusevalue_string.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("initial"),
					),
				},
			},
			// Update und Read Testing - Setter ändert den Wert
			{
				Config: `
resource "reuse_string" "test" {
    set_if_not_null_or_empty = "updated"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"reusevalue_string.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("updated"),
					),
				},
			},
			// Test: Wert bleibt erhalten wenn Setter weggelassen wird
			{
				Config: `
resource "reuse_string" "test" {
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"reusevalue_string.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("updated"),
					),
				},
			},
		},
	})
}
