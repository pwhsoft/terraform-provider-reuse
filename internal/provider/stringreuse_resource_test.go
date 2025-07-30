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
			// Create und Read Testing
			{
				Config: `
resource "reuse_reuse" "test" {
    set_if_not_null_or_empty = "initial"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"reuse_reuse.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("initial"),
					),
					// Setter sollte nach Apply null sein
					statecheck.ExpectKnownValue(
						"reuse_reuse.test",
						tfjsonpath.New("set_if_not_null_or_empty"),
						knownvalue.Null(),
					),
				},
			},
			// Update und Read Testing
			{
				Config: `
resource "reuse_reuse" "test" {
    set_if_not_null_or_empty = "updated"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"reuse_reuse.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("updated"),
					),
					// Setter sollte nach Apply null sein
					statecheck.ExpectKnownValue(
						"reuse_reuse.test",
						tfjsonpath.New("set_if_not_null_or_empty"),
						knownvalue.Null(),
					),
				},
			},
			// Test: Wert bleibt erhalten wenn Setter leer ist
			{
				Config: `
resource "reuse_reuse" "test" {
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"reuse_reuse.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("updated"),
					),
				},
			},
			// Import Testing
			{
				ResourceName:      "reuse_reuse.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Setter ist write-only und sollte ignoriert werden
				ImportStateVerifyIgnore: []string{"set_if_not_null_or_empty"},
			},
		},
	})
}
