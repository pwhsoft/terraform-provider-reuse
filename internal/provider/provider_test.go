package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories wird für Akzeptanztests verwendet
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"reuse": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Hier können später Voraussetzungen für Tests geprüft werden,
	// z.B. ob bestimmte Umgebungsvariablen gesetzt sind
}
